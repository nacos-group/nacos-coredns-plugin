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
	"net/http"
	"strings"
	"net/http/httptest"
)

func TestGet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func (w http.ResponseWriter, req *http.Request){

		if req.URL.EscapedPath() != "/nacos/v1/ns/api/srvIPXT" {
			t.Errorf("Except to path '/person',got '%s'", req.URL.EscapedPath())
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("hello"))
		}

		}))
	defer server.Close()

	s := Get(server.URL + "/nacos/v1/ns/api/srvIPXT", nil)

	if strings.Compare(s, "hello") == 0 {
		t.Log("Success to test http client get")
	} else {
		t.Fatal("Failed to test http client get")
	}
}
