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
	"os"
	"strings"
)

func TestServerManager_NextServer(t *testing.T) {
	os.Setenv("nacos_server_list", "2.2.2.2,3.3.3.3")
	sm := ServerManager{}
	sm.RefreshServerListIfNeed()
	ip := sm.NextServer()
	if strings.Compare(ip, "2.2.2.2") == 0 ||
		strings.Compare(ip, "3.3.3.3") == 0 {
		t.Log("ServerManager.NextServer test is passed.")
	}
}

func TestServerManager_RefreshServerListIfNeed(t *testing.T) {
	os.Setenv("nacos_server_list", "2.2.2.2,3.3.3.3")
	sm := ServerManager{}
	sm.RefreshServerListIfNeed()

	if len(sm.serverList) == 2 {
		t.Log("ServerManager.RefreshServerListIfNeed test is passed.")
	}

}
