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
	"strings"
	"math/rand"
	"reflect"
	"os"
)

type ServerManager struct {
	serverList      []string
	lastRefreshTime int64
	cursor          int
}

// get nacos ip list from address by env
func (manager *ServerManager) RefreshServerListIfNeed() []string {
	if CurrentMillis()-manager.lastRefreshTime < 60*1000 && len(manager.serverList) > 0 {
		return manager.serverList
	}

	nacosServerList := os.Getenv("nacos_server_list")

	var list []string
	list = strings.Split(nacosServerList, ",")

	var servers []string
	for _, line := range list {
		if line != "" {
			servers = append(servers, strings.TrimSpace(line))
		}
	}

	if len(servers) > 0 {
		if !reflect.DeepEqual(manager.serverList, servers) {
			NacosClientLogger.Info("server list is updated, old: ", manager.serverList, ", new: ", servers)
		}
		manager.serverList = servers

		manager.lastRefreshTime = CurrentMillis()
	}

	return manager.serverList
}

func (manager *ServerManager) NextServer() string {
	manager.RefreshServerListIfNeed()

	if len(manager.serverList) == 0 {
		panic("no nacos server avialible.")
	}

	return manager.serverList[rand.Intn(len(manager.serverList))]
}

func (manager *ServerManager) SetServers(servers []string) {
	manager.serverList = servers
}

func (manager *ServerManager) GetServerList() []string {
	return manager.serverList
}
