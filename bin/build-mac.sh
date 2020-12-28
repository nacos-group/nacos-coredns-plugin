#!/usr/bin/env bash
# cd GOPATH
cd $GOPATH/src/

# remove codes
rm -rf coredns
rm -rf nacos-coredns-plugin

# clone current codes
git clone https://github.com/coredns/coredns.git
git clone https://github.com/nacos-group/nacos-coredns-plugin.git

# cd coredns directory
cd $GOPATH/src/coredns
git checkout -b v1.6.7 v1.6.7
go get github.com/cihub/seelog

# copy nacos plugin to coredns
cp -r ../nacos-coredns-plugin/nacos plugin/
cp -r ../nacos-coredns-plugin/forward/setup.go plugin/forward

# insert nacos into plugin
sed -i '' '/hosts/a\
"nacos",\
' core/dnsserver/zdirectives.go
sed -i '' '/coredns\/plugin\/hosts/a\
_ "github.com/coredns/coredns/plugin/nacos"\
' core/plugin/zplugin.go
sed -i '' '/hosts:hosts/a\
nacos:nacos\
' plugin.cfg
# modify import
# build
make
