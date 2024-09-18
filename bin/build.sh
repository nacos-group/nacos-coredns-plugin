#!/usr/bin/env bash
# cd GOPATH
cd $GOPATH/src/
export PATH=$PATH:$GOROOT/bin
# remove codes
rm -rf coredns

# clone current codes
git clone https://github.com/coredns/coredns.git

# cd coredns directory
cd $GOPATH/src/coredns
git checkout -b v1.9.3 v1.9.3

# copy nacos plugin to coredns
cp -r ../nacos-coredns-plugin/nacos plugin/
cp -r ../nacos-coredns-plugin/forward/setup.go plugin/forward
cp -r ../nacos-coredns-plugin/conf conf

# insert nacos into plugin
sed -i '/hosts/a\\t"nacos",' core/dnsserver/zdirectives.go
sed -i '/coredns\/plugin\/hosts/a\\t_ "github.com/coredns/coredns/plugin/nacos"' core/plugin/zplugin.go
sed -i '/hosts:hosts/a\nacos:nacos' plugin.cfg

go mod tidy

# build
make