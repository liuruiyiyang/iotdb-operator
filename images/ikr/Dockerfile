#
# Licensed to the Apache Software Foundation (ASF) under one or more
# contributor license agreements.  See the NOTICE file distributed with
# this work for additional information regarding copyright ownership.
# The ASF licenses this file to You under the Apache License, Version 2.0
# (the "License"); you may not use this file except in compliance with
# the License.  You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

FROM openjdk:8-alpine

RUN apk add --no-cache bash gettext nmap-ncat openssl busybox-extras

ARG user=iotdb
ARG group=iotdb
ARG uid=3000
ARG gid=3000

# IoTDB is run with user `iotdb`, uid = 3000
# If you bind mount a volume from the host or a data container,
# ensure you use the same uid
RUN addgroup --gid ${gid} ${group} \
    && adduser --uid ${uid} -G ${group} ${user} -s /bin/bash -D

ARG version

# Rocketmq version
ENV IOTDB_VERSION ${version}

# Rocketmq home
ENV IOTDB_HOME  /home/iotdb/iotdb-${IOTDB_VERSION}

WORKDIR  ${IOTDB_HOME}

# Install
RUN set -eux; \
    apk add --virtual .build-deps curl gnupg unzip; \
    curl https://dist.apache.org/repos/dist/release/iotdb/${IOTDB_VERSION}/iotdb-all-${IOTDB_VERSION}-bin-release.zip -o iotdb.zip; \
    curl https://dist.apache.org/repos/dist/release/iotdb/${IOTDB_VERSION}/iotdb-all-${IOTDB_VERSION}-bin-release.zip.asc -o iotdb.zip.asc; \
    #https://www.apache.org/dist/iotdb/KEYS
	curl https://www.apache.org/dist/iotdb/KEYS -o KEYS; \
	\
	gpg --import KEYS; \
    gpg --batch --verify iotdb.zip.asc iotdb.zip; \
    unzip iotdb.zip; \
	mv iotdb-all*/* . ; \
	rmdir iotdb-all* ; \
	rm iotdb.zip iotdb.zip.asc KEYS; \
	apk del .build-deps ; \
    rm -rf /var/cache/apk/* ; \
    rm -rf /tmp/*

# Copy customized scripts
COPY runserver-customize.sh ${IOTDB_HOME}/bin/

RUN chown -R ${uid}:${gid} ${IOTDB_HOME}

# Expose ikr port
EXPOSE 9876

# Override customized scripts for ikr
# Export Java options
# Add ${JAVA_HOME}/lib/ext as java.ext.dirs
RUN mv ${IOTDB_HOME}/bin/runserver-customize.sh ${IOTDB_HOME}/bin/runserver.sh \
 && chmod a+x ${IOTDB_HOME}/bin/runserver.sh \
 && chmod a+x ${IOTDB_HOME}/bin/mqnamesrv \
 && export JAVA_OPT=" -Duser.home=/opt" \
 && sed -i 's/${JAVA_HOME}\/jre\/lib\/ext/${JAVA_HOME}\/jre\/lib\/ext:${JAVA_HOME}\/lib\/ext/' ${IOTDB_HOME}/bin/tools.sh

USER ${user}

WORKDIR ${IOTDB_HOME}/bin

CMD ["/bin/bash", "mqnamesrv"]
