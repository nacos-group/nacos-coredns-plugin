#!/usr/bin/env bash

#
# Copyright 1999-2024 Alibaba Group Holding Ltd.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

#
echo "inject 'nameserver 127.0.0.1' into /etc/resolv.conf"

cp /etc/resolv.conf /tmp/resolv.conf && sed -i '1i nameserver 127.0.0.1' /tmp/resolv.conf && cat /tmp/resolv.conf > /etc/resolv.conf
if [ $? -ne 0 ]; then
    echo "failed to inject 'nameserver 127.0.0.1' into /etc/resolv.conf"
    exit 1
fi
echo "namespace: $NACOS_NAMESPACE_ID"
if [ ! -z "$NACOS_NAMESPACE_ID" ]; then
    sed -i "s/public/${NACOS_NAMESPACE_ID}/g" /usr/local/nacos-agent/nacos-coredns.cfg
    if [ $? -ne 0 ]; then
        echo "failed to inject nacos namespace id into /usr/local/nacos-agent/nacos-coredns.cfg"
        exit 1
    fi
fi

echo "nacos server host: $NACOS_SERVER_HOST"
if [ ! -z "$NACOS_SERVER_HOST" ]; then
    sed -i "s/console.nacos.io:8848/${NACOS_SERVER_HOST}:8848/g" /usr/local/nacos-agent/nacos-coredns.cfg
    if [ $? -ne 0 ]; then
        echo "failed to inject nacos server host into /usr/local/nacos-agent/nacos-coredns.cfg"
        exit 1
    fi
fi

echo "accessKey: $NACOS_ACCESS_KEY"
if [ ! -z "$NACOS_ACCESS_KEY" ]; then
    sed -i "s/your_access_key/${NACOS_ACCESS_KEY}/g" /usr/local/nacos-agent/nacos-coredns.cfg
    if [ $? -ne 0 ]; then
        echo "failed to inject nacos accessKey into /usr/local/nacos-agent/nacos-coredns.cfg"
        exit 1
    fi
fi

echo "secretKey: $NACOS_SECRET_KEY"
if [ ! -z "$NACOS_SECRET_KEY" ]; then
    sed -i "s/your_secret_key/${NACOS_SECRET_KEY}/g" /usr/local/nacos-agent/nacos-coredns.cfg
    if [ $? -ne 0 ]; then
        echo "failed to inject nacos secretKey into /usr/local/nacos-agent/nacos-coredns.cfg"
        exit 1
    fi
fi

echo "startup nacos agent, listen port 53"
/usr/local/nacos-agent/coredns -conf /usr/local/nacos-agent/nacos-coredns.cfg -dns.port 53