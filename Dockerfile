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

FROM alpine:latest AS builder
RUN apk --no-cache add git make
RUN wget https://go.dev/dl/go1.23.0.linux-amd64.tar.gz -O /root/go1.23.0.linux-amd64.tar.gz && tar -xzf /root/go1.23.0.linux-amd64.tar.gz -C /root/
#RUN git clone https://github.com/nacos-group/nacos-coredns-plugin.git && cp nacos-coredns-plugin/bin/build.sh ./
COPY bin/build.sh /root/build.sh
ENV GOROOT=/root/go GOPATH=/root/go-code
COPY . /root/go-code/src/nacos-coredns-plugin
RUN mkdir -p /root/go-code/src/nacos-coredns-plugin /root/go-code/pkg /root/go-code/bin  && sh /root/build.sh
RUN mkdir -p /usr/local/nacos-agent/

FROM registry.cn-hangzhou.aliyuncs.com/nacos-agent/alpine:latest

COPY --from=builder /root/go-code/src/coredns/coredns /usr/local/nacos-agent/
COPY --from=builder /root/go-code/src/nacos-coredns-plugin/conf/nacos-coredns.cfg /usr/local/nacos-agent/
COPY start.sh /usr/local/nacos-agent/
CMD ["sh","/usr/local/nacos-agent/start.sh"]