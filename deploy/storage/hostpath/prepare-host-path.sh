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

### Instruction
### This script helps you create the hostPath dir on your localhost. You may need sudo privilege to run this script!
### You may need to run this script on every node of your Kubernetes cluster to prepare hostPath.

### You can modify the following parameters according to your real deployment requirements ###
# SERVER_DATA_PATH should be the same with the hostPath field in iotdb_v1alpha1_iotdb_cr.yaml
SERVER_DATA_PATH=/data/iotdb/server
# IOTDB_UID and IOTDB_GID should be the same with docker image settings
IOTDB_UID=2000
IOTDB_GID=2000


### prepare hostPath dir for IoTDB-Operator
prepare_dir()
{
  mkdir -p $1
  chown -R $2  $1
  chgrp -R $3  $1
  chmod -R a+rw $1
}

prepare_dir $SERVER_DATA_PATH $IOTDB_UID $IOTDB_GID

echo "Changed hostPath $SERVER_DATA_PATH uid to $IOTDB_UID, gid to $IOTDB_GID"

