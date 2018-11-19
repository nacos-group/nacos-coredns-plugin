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
	"net/http"
	"time"
	"strings"
	"io/ioutil"
	"strconv"
	"net/url"
)

var httpClient = http.Client{
	Timeout: time.Duration(10000 * time.Millisecond),
}

func encodeUrl(urlString string, params map[string]string) string {
	params["udpPort"] = strconv.Itoa(UDP_Port)

	if !strings.HasSuffix(urlString, "?") {
		urlString += "?"
	}

	if params == nil || len(params) == 0 {
		return urlString
	}

	u := url.Values{}
	for key, value := range params {
		u.Set(key, value)
	}

	return urlString + u.Encode()
}

func Get(url string, params map[string]string) string {
	if params == nil {
		params = make(map[string]string)
	}

	url = encodeUrl(url, params)

	req, err := http.NewRequest("GET",url, nil)
	if err != nil {
		NacosClientLogger.Error("failed to build request", err)
		return ""
	}

	req.Header.Add("Client-Version", Version)
	response, err := httpClient.Do(req)

	if err != nil || response.StatusCode != 200 {
		NacosClientLogger.Error("error while request from " + url, err)
		if err != nil {
			NacosClientLogger.Error("error while request from " + url, err)
		} else {
			NacosClientLogger.Warn("error while request from " + url + ", code: " + strconv.Itoa(response.StatusCode))
		}
		return ""
	}


	b, err := ioutil.ReadAll(response.Body)
	response.Body.Close()

	if err != nil {
		NacosClientLogger.Error("failed to get response body: " + url, err)
		return ""
	}

	bs := string(b)
	return bs
}