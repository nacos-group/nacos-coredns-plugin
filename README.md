This project provides a coredns plugin, which can export nacos services as dns domain.

## Quic Start
To build an run nacos coredns plugin, the OS must be Linux or Mac. And also, golang environments(GOPATH,GOROOT,GOHOME) must be configured correctly.

### Build
```
git clone git@gitlab.alibaba-inc.com:middleware/nacos-coredns-plugin.git 
cp nacos-coredns-plugin/bin/build.sh ~/
cd ~/
sh build.sh
```
### Configuration
To run nacos coredns plugin, you need a configuration file. A possible file may be as bellow:
```
. {
    nacos {
        upstream /etc/resolv.conf
        nacos_server 127.0.0.1
        nacos_server_port 8848
   }
 }
```
* upstream: domain names those not registered in nacos will be forwarded to upstream.
* nacos_server: Ips of nacos server, seperated by comma if there are two or more nacos servers
* nacos_server_port: Nacos server port

### Run
* Firstly, you need to deploy nacos server. [Here](https://github.com/alibaba/nacos)
* Secondly, register service on nacos.
* Then run ```$GOPATH/src/coredns/coredns -conf $path_to_corefile -dns.port $dns_port```
![image](http://git.cn-hangzhou.oss-cdn.aliyun-inc.com/uploads/middleware/nacos-coredns-plugin/5f231b0ad438d9fc62299a7d94e42cca/image.png)

### Test
dig $nacos_service_name @127.0.0.1 -p $dns_port

![image](http://git.cn-hangzhou.oss-cdn.aliyun-inc.com/uploads/middleware/nacos-coredns-plugin/464d59644bbb76c227a3900d8459bc9c/image.png)
