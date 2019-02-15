This project provides a **DNS-F** client based on [CoreDNS](https://coredns.io/), which can help export those registed services on Nacos as DNS domain. **DNS-F** client is a dedicated agent process(side car) beside the application's process to foward the service discovery DNS domain query request to Nacos. 

## Quic Start
To build and run nacos coredns plugin, the OS must be Linux or Mac. And also, golang environments(GOPATH,GOROOT,GOHOME) must be configured correctly.

### Build
```
git clone https://github.com/nacos-group/nacos-coredns-plugin.git 
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
![image](https://cdn.nlark.com/lark/0/2018/png/7601/1542623914418-f529409b-c229-4ef9-aec3-b9c5df23c906.png)

### Test
dig $nacos_service_name @127.0.0.1 -p $dns_port

![image](https://cdn.nlark.com/lark/0/2018/png/7601/1542624023214-29cd9f71-0183-4231-b092-57535e8cfcfe.png)
