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
	"math"
	"encoding/json"
)

type Domain struct {
	Name string `json:"dom"`
	Clusters string
	CacheMillis int64
	LastRefMillis int64
	Instances []Instance `json:"hosts"`
	Env string
	TTL int

}

func (domain Domain) getInstances() ([]Instance) {
	return domain.Instances
}

func (domain Domain) String() string {
	b, _ := json.Marshal(domain)
	return string(b)
}

func (domain Domain) SrvInstances() []Instance {
	var result = make([]Instance, 0)
	hosts := domain.getInstances()
	for _, host := range hosts {
		if host.Valid && host.Weight > 0 {
			for i := 0; i < int(math.Ceil(host.Weight)); i++ {
				result = append(result, host)
			}
		}
	}

	if len(result) <= 0{
		panic("no host to srv: " + domain.Name)
	}

	return result
}
