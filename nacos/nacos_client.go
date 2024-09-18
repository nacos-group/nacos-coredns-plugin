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
	"encoding/json"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cihub/seelog"
)

var (
	NacosClientLogger seelog.LoggerInterface
	LogConfig         string
	GrpcClient        *NacosGrpcClient
)

func init() {
	initLog()
}

type NacosClient struct {
	serviceMap ConcurrentMap
	udpServer  UDPServer
}

type NacosClientError struct {
	Msg string
}

func (err NacosClientError) Error() string {
	return err.Msg
}

var Inited = false

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func initDir() {
	dir, err := filepath.Abs(Home())
	if err != nil {
		os.Exit(1)
	}

	if LogPath == "" {
		LogPath = dir + string(os.PathSeparator) + "logs"
	}

	if CachePath == "" {
		CachePath = dir + string(os.PathSeparator) + "nacos-go-client-cache"
	}

	mkdirIfNecessary(CachePath)
}

func mkdirIfNecessary(path string) {
	if ok, _ := exists(path); !ok {
		err := os.Mkdir(path, 0755)
		if err != nil {
			NacosClientLogger.Warn("can not cread dir: "+path, err)
		}
	}
}

func initLog() {
	if NacosClientLogger != nil {
		return
	}

	initDir()
	var err error
	var nacosLogger seelog.LoggerInterface
	if LogConfig == "" || !Exist(LogConfig) {
		home := Home()
		LogConfig = home
		nacosLogger, err = seelog.LoggerFromConfigAsString(`
    		<seelog minlevel="info" maxlevel="error">
        		<outputs formatid="main">
            		<buffered size="1000" flushperiod="1000">
                		<rollingfile type="size" filename="` + home + `/logs/nacos-go-client/nacos-go-client.log" maxsize="100000000" maxrolls="10"/>
        			</buffered>
    			</outputs>
    			<formats>
        			<format id="main" format="%Date(2006-01-02 15:04:05.000) %LEVEL %Msg %n"/>
    			</formats>
			</seelog>
			`)
	} else {
		nacosLogger, err = seelog.LoggerFromConfigAsFile(LogConfig)
	}

	fmt.Println("log directory: " + LogConfig + "/logs/nacos-go-client/")

	if err != nil {
		fmt.Println("Failed to init log, ", err)
	} else {
		NacosClientLogger = nacosLogger
	}
}

func (nacosClient *NacosClient) asyncGetAllServiceNames() {
	for {
		time.Sleep(time.Duration(AllDoms.CacheSeconds) * time.Second)
		nacosClient.getAllServiceNames()
	}
}

//func (nacosClient *NacosClient) GetServerManager() (serverManager *ServerManager) {
//	return &nacosClient.serverManager
//}

func (nacosClient *NacosClient) GetUdpServer() (us UDPServer) {
	return nacosClient.udpServer
}

func (nacosClient *NacosClient) getAllServiceNames() {

	services := GrpcClient.GetAllServicesInfo()
	if services == nil {
		NacosClientLogger.Warn("No Service return from servers.")
		return
	}

	AllDoms.DLock.Lock()
	if AllDoms.Data == nil {
		allDoms := make(map[string]bool)
		// record all serviceNames return from server
		for _, service := range services {
			allDoms[service] = true
		}
		AllDoms.Data = allDoms
		AllDoms.CacheSeconds = 20 //刷新间隔
	} else {
		for _, service := range services {
			if !AllDoms.Data[service] {
				AllDoms.Data[service] = true
			}
		}
	}
	AllDoms.DLock.Unlock()
}

//func (nacosClient *NacosClient) SetServers(servers []string) {
//	nacosClient.serverManager.SetServers(servers)
//}

func (vc *NacosClient) Registered(service string) bool {
	defer AllDoms.DLock.RUnlock()
	AllDoms.DLock.RLock()
	_, ok1 := AllDoms.Data[service]

	return ok1
}

func (vc *NacosClient) loadCache() {
	NacosSdkCachePath := CachePath + "/naming/public/"
	files, err := ioutil.ReadDir(NacosSdkCachePath)
	if err != nil {
		NacosClientLogger.Critical(err)
	}

	for _, f := range files {
		fileName := NacosSdkCachePath + string(os.PathSeparator) + f.Name()
		b, err := ioutil.ReadFile(fileName)
		if err != nil {
			NacosClientLogger.Error("failed to read cache file: "+fileName, err)
		}

		s := string(b)
		var service model.Service
		err1 := json.Unmarshal([]byte(s), &service)

		if err1 != nil {
			continue
		}

		vc.serviceMap.Set(f.Name(), service)
	}

	NacosClientLogger.Info("finish loading cache, total: " + strconv.Itoa(len(files)))
}

func ProcessDomainString(s string) (model.Service, error) {
	var service model.Service
	err1 := json.Unmarshal([]byte(s), &service)

	if err1 != nil {
		NacosClientLogger.Error("failed to unmarshal json string: "+s, err1)
		return model.Service{}, err1
	}

	if len(service.Hosts) == 0 {
		NacosClientLogger.Warn("get empty ip list, ignore it, service: " + service.Name)
		return service, NacosClientError{"empty ip list"}
	}

	NacosClientLogger.Info("domain "+service.Name+" is updated, current ips: ", service.Hosts)
	return service, nil
}

func NewNacosClient(namespaceId string, serverHosts []string, accessKey string, secretKey string) *NacosClient {
	fmt.Println("init nacos client.")
	initLog()
	vc := NacosClient{NewConcurrentMap(), UDPServer{}}
	vc.loadCache()
	vc.udpServer.vipClient = &vc
	//init grpcClient
	var err error
	GrpcClient, err = NewNacosGrpcClient(namespaceId, serverHosts, &vc, accessKey, secretKey)
	if err != nil {
		NacosClientLogger.Error("init nacos-grpc-client failed.", err)
	}

	if EnableReceivePush {
		go vc.udpServer.StartServer()
	}

	AllDoms = AllDomsMap{}
	AllDoms.Data = make(map[string]bool)
	AllDoms.DLock = sync.RWMutex{}
	AllDoms.CacheSeconds = 20

	vc.getAllServiceNames()

	go vc.asyncGetAllServiceNames()

	//go vc.asyncUpdateDomain()

	NacosClientLogger.Info("cache-path: " + CachePath)
	return &vc
}

func (vc *NacosClient) GetDomainCache() ConcurrentMap {
	return vc.serviceMap
}

func (vc *NacosClient) GetDomain(name string) (*Domain, error) {
	item, _ := vc.serviceMap.Get(name)

	if item == nil {
		domain := Domain{}
		ss := strings.Split(name, SEPERATOR)

		domain.Name = ss[0]
		domain.CacheMillis = DefaultCacheMillis
		domain.LastRefMillis = CurrentMillis()
		vc.serviceMap.Set(name, domain)
		item = domain
		return nil, NacosClientError{"domain not found: " + name}
	}

	domain := item.(Domain)

	return &domain, nil
}

func (vc *NacosClient) asyncUpdateDomain() {
	for {
		for k, v := range vc.serviceMap.Items() {
			service := v.(model.Service)
			ss := strings.Split(k, SEPERATOR)

			serviceName := ss[0]
			var clientIP string
			if len(ss) > 1 && ss[1] != "" {
				clientIP = ss[1]
			}

			if uint64(CurrentMillis())-service.LastRefTime > service.CacheMillis && vc.Registered(serviceName) {
				vc.getServiceNow(serviceName, &vc.serviceMap, clientIP)
			}
		}
		time.Sleep(1 * time.Second)
	}

}

func GetCacheKey(dom, clientIP string) string {
	return dom + SEPERATOR + clientIP
}

func (vc *NacosClient) getServiceNow(serviceName string, cache *ConcurrentMap, clientIP string) model.Service {
	service := GrpcClient.GetService(serviceName)

	cache.Set(serviceName, service)

	NacosClientLogger.Info("dom "+serviceName+" updated: ", service)

	return service
}

func (vc *NacosClient) SrvInstance(serviceName, clientIP string) *model.Instance {
	item, hasService := vc.serviceMap.Get(serviceName)
	var service model.Service
	if !hasService {
		service = vc.getServiceNow(serviceName, &vc.serviceMap, clientIP)
		vc.serviceMap.Set(serviceName, service)
	} else {
		service = item.(model.Service)
	}

	//select healthy instances
	var hosts []model.Instance
	for _, host := range service.Hosts {
		if host.Healthy && host.Enable && host.Weight > 0 {
			hosts = append(hosts, host)
		}
	}

	if len(hosts) == 0 {
		NacosClientLogger.Warn("no healthy instances for " + serviceName)
		return nil
	}

	i, indexOk := indexMap.Get(serviceName)
	var index int

	if !indexOk {
		index = rand.Intn(len(hosts))
	} else {
		index = i.(int)
		index += 1
		if index >= len(hosts) {
			index = index % len(hosts)
		}
	}

	indexMap.Set(serviceName, index)

	return &hosts[index]
}

func (vc *NacosClient) SrvInstances(domainName, clientIP string) []model.Instance {
	cacheKey := GetCacheKey(domainName, clientIP)
	item, hasDom := vc.serviceMap.Get(cacheKey)
	var dom model.Service

	if !hasDom {
		dom = vc.getServiceNow(domainName, &vc.serviceMap, clientIP)
		vc.serviceMap.Set(cacheKey, dom)
	} else {
		dom = item.(model.Service)
	}

	var hosts []model.Instance
	//select healthy instances
	for _, host := range dom.Hosts {
		if host.Healthy && host.Enable && host.Weight > 0 {
			for i := 0; i < int(math.Ceil(host.Weight)); i++ {
				hosts = append(hosts, host)
			}
		}
	}

	return hosts
}

func (vc *NacosClient) Contains(dom, clientIP string, host model.Instance) bool {
	hosts := vc.SrvInstances(dom, clientIP)

	for _, host1 := range hosts {
		if reflect.DeepEqual(host1, host) {
			return true
		}
	}

	return false
}
