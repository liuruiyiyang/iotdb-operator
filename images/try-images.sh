#!/bin/bash

# Licensed to the Apache Software Foundation (ASF) under one or more
# contributor license agreements.  See the NOTICE file distributed with
# this work for additional information regarding copyright ownership.
# The ASF licenses this file to You under the Apache License, Version 2.0
# (the "License"); you may not use this file except in compliance with
# the License.  You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

NAMESRV_DOCKERHUB_REPO=2019liurui/iotdb-namesrv
BROKER_DOCKERHUB_REPO=2019liurui/iotdb-iotdb
IOTDB_VERSION=4.5.0

start_namesrv_broker()
{
    TAG_SUFFIX=$1
    # Start nameserver
    docker run -d -v `pwd`/data/namesrv/logs:/home/iotdb/logs -v `pwd`/data/namesrv/store:/home/iotdb/store --name rmqnamesrv ${NAMESRV_DOCKERHUB_REPO}:${IOTDB_VERSION}${TAG_SUFFIX}
    # Start IoTDB
    docker run -d -v `pwd`/data/iotdb/logs:/home/iotdb/logs -v `pwd`/data/iotdb/store:/home/iotdb/store --name rmqbroker --link rmqnamesrv:namesrv -e "NAMESRV_ADDR=namesrv:9876" ${BROKER_DOCKERHUB_REPO}:${IOTDB_VERSION}${TAG_SUFFIX}
}

#if [ $# -lt 1 ]; then
#    echo -e "Usage: sh $0 BaseImage"
#    exit 2
#fi

export BASE_IMAGE=alpine

RMQ_CONTAINER=$(docker ps -a|awk '/rmq/ {print $1}')
if [[ -n "$RMQ_CONTAINER" ]]; then
   echo "Removing IoTDB Container..."
   docker rm -fv $RMQ_CONTAINER
   # Wait till the existing containers are removed
   sleep 5
fi

if [ ! -d "`pwd`/data" ]; then
  mkdir -p "data"
fi

echo "Play IoTDB namesrv and broker image"
echo "Starting IoTDB nodes..."

case "${BASE_IMAGE}" in
    alpine)
        start_namesrv_broker -alpine
    ;;
    centos)
        start_namesrv_broker
    ;;
    *)
        echo "${BASE_IMAGE} is not supported, supported base images: centos, alpine"
        exit 2
    ;;
esac

echo "Wait 10 seconds for service ready"

sleep 10

echo "Start producer..."
# Produce messages
docker exec -ti rmqbroker sh ./tools.sh org.apache.iotdb.example.quickstart.Producer
sleep 2
echo "Start consumer..."
# Consume messages
docker exec -ti rmqbroker sh ./tools.sh org.apache.iotdb.example.quickstart.Consumer
