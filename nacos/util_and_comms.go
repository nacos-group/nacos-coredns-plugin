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
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"net"
	"time"
)

var (
	DefaultCacheMillis = int64(5000)
	Version            = "Nacos-DNS:v1.0.1"
	CachePath          string
	LogPath            string
	SEPERATOR          = "@@"
	GZIP_MAGIC         = []byte("\x1F\x8B")
	EnableReceivePush  = true
	UDP_Port           = -1
	SERVER_PORT        = "8848"
)

func CurrentMillis() int64 {
	return time.Now().UnixNano() / 1e6
}

func TryDecompressData(data []byte) string {

	if !IsGzipFile(data) {
		return string(data)
	}

	reader, err := gzip.NewReader(bytes.NewReader(data))

	if err != nil {
		NacosClientLogger.Warn("failed to decompress gzip data", err)
		return ""
	}

	defer reader.Close()
	bs, err1 := ioutil.ReadAll(reader)

	if err1 != nil {
		NacosClientLogger.Warn("failed to decompress gzip data", err1)
		return ""
	}

	return string(bs)
}

func IsGzipFile(data []byte) bool {
	if len(data) < 2 {
		return false
	}

	return bytes.HasPrefix(data, GZIP_MAGIC)
}

var localIP = ""

func LocalIP() string {
	if localIP == "" {
		addrs, err := net.InterfaceAddrs()
		if err != nil {
			return ""
		}
		for _, address := range addrs {
			if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					localIP = ipnet.IP.String()
					break
				}
			}
		}
	}

	return localIP
}
