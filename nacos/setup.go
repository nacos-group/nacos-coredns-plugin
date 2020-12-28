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
	"fmt"
	"github.com/caddyserver/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"strconv"
	"strings"
)

func init() {
	caddy.RegisterPlugin("nacos", caddy.Plugin{
		ServerType: "dns",
		Action:     setup,
	})
	fmt.Println("register nacos plugin")
}

func setup(c *caddy.Controller) error {
	fmt.Println("setup nacos plugin")
	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		vs, _ := NacosParse(c)
		vs.Next = next
		Inited = true
		return vs
	})
	return nil
}

func NacosParse(c *caddy.Controller) (*Nacos, error) {
	fmt.Println("init nacos plugin...")
	nacosImpl := Nacos{}
	var servers = make([]string, 0)
	serverPort := 8848
	for c.Next() {
		nacosImpl.Zones = c.RemainingArgs()

		if c.NextBlock() {
			for {
				switch v := c.Val(); v {
				case "nacos_server":
					servers = strings.Split(c.RemainingArgs()[0], ",")
					/* it is nacos_servera noop now */
				case "nacos_server_port":
					port, err := strconv.Atoi(c.RemainingArgs()[0])
					if err != nil {
						serverPort = port
					}
				case "cache_ttl":
					ttl, err := strconv.Atoi(c.RemainingArgs()[0])
					if err != nil {
						DNSTTL = uint32(ttl)
					}
				case "cache_dir":
					CachePath = c.RemainingArgs()[0]
				case "log_path":
					LogPath = c.RemainingArgs()[0]
				default:
					if c.Val() != "}" {
						return &Nacos{}, c.Errf("unknown property '%s'", c.Val())
					}
				}

				if !c.Next() {
					break
				}
			}

		}

		client := NewNacosClient(servers, serverPort)
		nacosImpl.NacosClientImpl = client
		nacosImpl.DNSCache = NewConcurrentMap()

		return &nacosImpl, nil
	}
	return &Nacos{}, nil
}
