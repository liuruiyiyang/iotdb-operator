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
ARG uid=2000
ARG gid=2000

# IoTDB is run with user `iotdb`, uid = 2000
# If you bind mount a volume from the host or a data container,
# ensure you use the same uid
RUN addgroup --gid ${gid} ${group} \
    && adduser --uid ${uid} -G ${group} ${user} -s /bin/bash -D

ARG version

# IoTDB version
ENV IOTDB_VERSION ${version}

# IoTDB home
ENV IOTDB_HOME  /home/iotdb/iotdb-${IOTDB_VERSION}

WORKDIR  ${IOTDB_HOME}

COPY iotdb-0.9.0.zip ${IOTDB_HOME}/

# Install
RUN unzip iotdb-0.9.0.zip; \
	mv iotdb-0.9.0*/* . ; \
	rmdir iotdb-0.9.0 ; \
	rm iotdb-0.9.0.zip ; \
	apk del .build-deps ; \
    rm -rf /var/cache/apk/* ; \
    rm -rf /tmp/*


RUN chown -R ${uid}:${gid} ${IOTDB_HOME}

# Expose iotdb ports
EXPOSE 6667 31999

## Export Java options
#RUN export JAVA_OPT=" -Duser.home=/opt"
#
## Add ${JAVA_HOME}/lib/ext as java.ext.dirs
#RUN sed -i 's/${JAVA_HOME}\/jre\/lib\/ext/${JAVA_HOME}\/jre\/lib\/ext:${JAVA_HOME}\/lib\/ext/' ${IOTDB_HOME}/bin/tools.sh

COPY genConfig.sh start.sh ${IOTDB_HOME}/server/sbin/

RUN chmod a+x ${IOTDB_HOME}/server/sbin/genConfig.sh \
 && chmod a+x ${IOTDB_HOME}/server/sbin/start.sh

USER ${user}

WORKDIR ${IOTDB_HOME}/server/sbin

ENV IOTDB_HOME  /home/iotdb/iotdb-${IOTDB_VERSION}/server

CMD ["/bin/bash", "./start.sh"]