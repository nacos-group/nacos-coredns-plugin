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
	"testing"
	"net"
	"time"
	"strings"
)

func TestUDPServer_StartServer(t *testing.T) {
	s := `{"dom":"hello123","cacheMillis":10000,"useSpecifiedURL":false,"hosts":[{"valid":true,"marked":false,"metadata":{},"instanceId":"","port":80,"ip":"2.2.2.2","weight":1.0,"enabled":true}],"checksum":"c7befb32f3bb5b169f76efbb0e1f79eb1542236821437","lastRefTime":1542236821437,"env":"","clusters":""}`
	us := UDPServer{}
	us.vipClient = &NacosClient{NewConcurrentMap(), UDPServer{}, ServerManager{}, 8848 }
	go us.StartServer()

	time.Sleep(100000)
	sip := net.ParseIP("127.0.0.1")
	srcAddr := &net.UDPAddr{IP: net.IPv4zero, Port: 0}
	dstAddr := &net.UDPAddr{IP: sip, Port: us.port}
	conn, err := net.DialUDP("udp", srcAddr, dstAddr)
	if err != nil {
		t.Error("Udp server test failed")
	}
	defer conn.Close()
	conn.Write([]byte(s))
	data := make([]byte, 4024)
	n, _, _ := conn.ReadFromUDP(data)
	if strings.Contains(string(data[:n]), "push-ack") {
		t.Log("Udp server test passed.")
	}
}
