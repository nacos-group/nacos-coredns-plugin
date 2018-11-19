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
	"os"
	"net"
	"strconv"
	"math/rand"
	json "encoding/json"
	"time"
)

type UDPServer struct {
	port int
	host string
	vipClient *NacosClient
}

type PushData struct {
	PushType string `json:"type"`
	Data string `json:"data"`
	LastRefTime int64 `json:"lastRefTime"`
	
}

func getUdpPort() int {
	return 0
}

func (us *UDPServer) tryListen() (*net.UDPConn, bool) {
	addr, err := net.ResolveUDPAddr("udp", us.host+":"+ strconv.Itoa(us.port))
	if err != nil {
		NacosClientLogger.Error("Can't resolve address: ", err)
		return nil , false
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		NacosClientLogger.Error("Error listening:", err)
		return nil, false
	}

	return conn, true
}

func (us *UDPServer) SetNacosClient(nc *NacosClient) {
	us.vipClient = nc
}

func (us *UDPServer) StartServer(){
	var conn *net.UDPConn

	for i := 0; i < 3; i++ {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		port := r.Intn(1000) + 54951
		us.port = port
		conn1, ok := us.tryListen()

		if ok {
			conn = conn1
			NacosClientLogger.Info("udp server start, port: " + strconv.Itoa(port))
			break
		}

		if !ok && i == 2 {
			NacosClientLogger.Critical("failed to start udp server after trying 3 times.")
			os.Exit(1)
		}
	}

	UDP_Port = us.port

	defer conn.Close()
	for {
		us.handleClient(conn)
	}
}

func (us *UDPServer) handleClient(conn *net.UDPConn) {
	data := make([]byte, 4024)
	n, remoteAddr, err := conn.ReadFromUDP(data)
	if err != nil {
		NacosClientLogger.Error("failed to read UDP msg because of ", err)
		return
	}

	s := TryDecompressData(data[:n])

	NacosClientLogger.Info("receive push: " + s + " from: ", remoteAddr)

	var pushData PushData
	err1 := json.Unmarshal([]byte(s), &pushData)
	if err1 != nil {
		NacosClientLogger.Warn("failed to process push data, ", err1)
		return
	}

	domain, err1 := ProcessDomainString(pushData.Data)
	NacosClientLogger.Info("receive domain: " , domain)

	if err1 != nil {
		NacosClientLogger.Warn("failed to process push data: " + s, err1)
	}

	key := GetCacheKey(domain.Name, LocalIP())

	us.vipClient.domainMap.Set(key, domain)

	ack := make(map[string]string)
	ack["type"] = "push-ack"
	ack["lastRefTime"] = strconv.FormatInt(pushData.LastRefTime, 10)
	ack["data"] = ""

	bs,_ := json.Marshal(ack)

	conn.WriteToUDP(bs, remoteAddr)
}

