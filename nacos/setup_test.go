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
	"github.com/mholt/caddy"
	"strings"
	"fmt"
	os "os"
)

func TestNacosParse(t *testing.T) {
	tests := []struct {
		input              string
		expectedPath       string
		expectedEndpoint   string
		expectedErrContent string // substring from the expected error. Empty for positive cases.
	}{
		{
			`nacos nacos.local {
			upstream 8.8.8.8:53
			nacos_server 192.168.0.1
        	nacos_server_port 8848
        	cache_ttl 1
			}
`, "skydns", "localhost:300", "",
		},
	}

	os.Unsetenv("nacos_server_list")

	for _, test := range tests {
		c := caddy.NewTestController("dns", test.input)
		nacosimpl, err :=NacosParse(c)
		if err != nil {
			t.Error("Failed to get instance.");
		} else {
			var passed bool
			fmt.Println(nacosimpl.NacosClientImpl.GetServerManager().GetServerList())
			for _, item := range nacosimpl.NacosClientImpl.GetServerManager().GetServerList() {
				if strings.Compare(item, "192.168.0.1") == 0 {
					t.Log("Passed")
					passed = true
				}
			}

			if !passed {
				t.Fatal("Failed")
			}
		}
	}
}
