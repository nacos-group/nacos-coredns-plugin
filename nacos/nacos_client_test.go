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
	"strings"
	"net/http/httptest"
	"net/http"
	"strconv"
)

func TestNacosClient_GetDomain(t *testing.T) {
	s := `{"dom":"hello123","cacheMillis":10000,"useSpecifiedURL":false,"hosts":[{"valid":true,"marked":false,"metadata":{},"instanceId":"","port":81,"ip":"2.2.2.2","weight":1.0,"enabled":true}],"checksum":"c7befb32f3bb5b169f76efbb0e1f79eb1542236821437","lastRefTime":1542236821437,"env":"","clusters":""}`
	server := httptest.NewServer(http.HandlerFunc(func (w http.ResponseWriter, req *http.Request){
		if req.URL.EscapedPath() == "/nacos/v1/ns/api/srvIPXT" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(s))
		} else if req.URL.EscapedPath() == "/nacos/v1/ns/api/allDomNames"  {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("{\"count\":1,\"doms\":[\"hello123\"]}"))
		}

	}))

	port, _ := strconv.Atoi(strings.Split(server.URL, ":")[2])

	defer server.Close()

	vc := NacosClient{NewConcurrentMap(), UDPServer{}, ServerManager{}, port}
	vc.udpServer.vipClient = &vc
	vc.SetServers([]string{strings.Split(strings.Split(server.URL, "http://")[1], ":")[0]})
	instance := vc.SrvInstance("hello123", "127.0.0.1")
	if strings.Compare(instance.IP, "2.2.2.2") == 0 {
		t.Log("Passed")
	}
}
