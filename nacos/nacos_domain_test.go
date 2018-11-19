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
)

func TestDomain_SrvInstances(t *testing.T) {
	domain := Domain{}
	domain.CacheMillis = 10000
	domain.Clusters = "DEFAULT"

	//test weight
	domain.Instances = []Instance{Instance{IP: "2.2.2.2", Port: 80, Weight: 2, AppUseType: "publish", Valid: true, Site: "et2"}}
	instances := domain.SrvInstances()
	if len(instances) == 2 {
		t.Log("Domain.srvInstances weight passed.")
	}

	//test valid
	defer func() {
		if err := recover(); err != nil {
			if strings.HasPrefix(err.(string), "no host to srv: ") {
				t.Log("Domain.srvInstances valid passed.")
			}
		}
	}()
	domain.Instances = []Instance{Instance{IP: "2.2.2.2", Port: 80, Weight: 2, AppUseType: "publish", Valid: false, Site: "et2"}}
	domain.SrvInstances()

}
