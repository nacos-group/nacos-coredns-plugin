/*
 * Copyright 1999-2018 Alibaba Group Holding Ltd.
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *      http://www.apache.org/licenses/LICENSE-2.0
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package nacos

import (
	"github.com/miekg/dns"
	"net"
	"github.com/coredns/coredns/plugin/proxy"
	"github.com/coredns/coredns/plugin"
	"time"
	"strconv"
	"encoding/json"
	"github.com/coredns/coredns/request"
	"context"
)

type Nacos struct {
	Next        plugin.Handler
	Zones       []string
	Proxy       proxy.Proxy
	NacosClientImpl  *NacosClient
	DNSCache    ConcurrentMap
}

func (vs *Nacos) String() string {
	b, err := json.Marshal(vs)

	if err != nil {
		return ""
	}

	return string(b)
}

// Lookup implements the ServiceBackend interface.
func (e *Nacos) Lookup(state request.Request, name string, typ uint16) (*dns.Msg, error) {
	key := name + strconv.Itoa(state.Family())
	msg, ok := e.DNSCache.Get(key)

	NacosClientLogger.Info("lookup " + name + " from upstream ")
	if ok {
		dnsCache := msg.(DnsCache)
		if !dnsCache.Updated() {
			msg1, err := e.Proxy.Lookup(state, name, typ)
			if err == nil {
				if len(msg1.Answer) > 0 {
					dnsCache.Msg = msg1
					dnsCache.LastUpdateMills = time.Now().UnixNano() / 1000000
					e.DNSCache.Set(key, dnsCache)
				}
			} else {
				NacosClientLogger.Warn("error while lookup dom: ", err)
			}
		}
		bs, err := json.Marshal(dnsCache.Msg)

		if err == nil {
			NacosClientLogger.Info("Forward " + name + " -> " + string(bs))
		}

		return dnsCache.Msg, nil
	} else {
		msg1, err := e.Proxy.Lookup(state, name, typ)
		if err == nil {
			dnsCache := DnsCache{Msg: msg1, LastUpdateMills: time.Now().UnixNano() / 1000000}
			e.DNSCache.Set(name, dnsCache)
		} else {
			NacosClientLogger.Warn("error while lookup dom: ", err)
		}

		bs, err := json.Marshal(msg1)

		if err == nil {
			NacosClientLogger.Info("Forward " + name + " -> " + string(bs))
		}

		return msg1, err
	}
}

func (vs *Nacos) managed(dom, clientIP string) bool {
	if _, ok := DNSDomains[dom]; ok {
		return false
	}

	defer AllDoms.DLock.RUnlock()

	AllDoms.DLock.RLock()
	_, ok1 := AllDoms.Data[dom]

	cacheKey := GetCacheKey(dom, clientIP)

	_, inCache := vs.NacosClientImpl.GetDomainCache().Get(cacheKey)

	return ok1 || inCache
}

func (vs *Nacos) getRecordBySession(dom, clientIP string) Instance {
	host := *vs.NacosClientImpl.SrvInstance(dom, clientIP)
	return host

}

func (vs *Nacos) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}

	name := state.QName()

	m := new(dns.Msg)

	clientIP := state.IP()
	if clientIP == "127.0.0.1" {
		clientIP = LocalIP()
	}

	if !vs.managed(name[:len(name)-1], clientIP) {
		dnsMsg, _ := vs.Lookup(state, name, state.QType())
		m.Answer = dnsMsg.Answer
		m.Extra = dnsMsg.Extra

	} else {
		hosts := make([]Instance, 0)
		host := vs.NacosClientImpl.SrvInstance(name[:len(name)-1], clientIP)
		hosts = append(hosts, *host)

		answer := make([]dns.RR, 0)
		extra := make([]dns.RR, 0)
		for _, host := range hosts {
			var rr dns.RR

			switch state.Family() {
			case 1:
				rr = new(dns.A)
				rr.(*dns.A).Hdr = dns.RR_Header{Name: state.QName(), Rrtype: dns.TypeA, Class: state.QClass(), Ttl: DNSTTL}
				rr.(*dns.A).A = net.ParseIP(host.IP).To4()
			case 2:
				rr = new(dns.AAAA)
				rr.(*dns.AAAA).Hdr = dns.RR_Header{Name: state.QName(), Rrtype: dns.TypeAAAA, Class: state.QClass(), Ttl: DNSTTL}
				rr.(*dns.AAAA).AAAA = net.ParseIP(host.IP)
			}

			srv := new(dns.SRV)
			srv.Hdr = dns.RR_Header{Name: "_" + state.Proto() + "." + state.QName(), Rrtype: dns.TypeSRV, Class: state.QClass(), Ttl: DNSTTL}
			port := host.Port
			srv.Port = uint16(port)
			srv.Target = "."

			extra = append(extra, srv)
			answer = append(answer, rr)
		}

		m.Answer = answer
		m.Extra = extra
		result, _ := json.Marshal(m.Answer)
		NacosClientLogger.Info("[RESOLVE]",  " [" + name[:len(name)-1] + "]  result: " + string(result) + ", clientIP: " + clientIP)
	}

	m.SetReply(r)
	m.Authoritative, m.RecursionAvailable, m.Compress = true, true, true

	state.SizeAndDo(m)
	m = state.Scrub(m)
	w.WriteMsg(m)
	return dns.RcodeSuccess, nil
}

func (vs *Nacos) Name() string { return "nacos" }
