/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Package constants defines some global constants
package constants

const (
	// IoTDBContainerName is the name of broker container
	IoTDBContainerName = "broker"

	// BasicCommand is basic command of exec function
	BasicCommand = "sh"

	// AdminToolDir is the RocketMQ Admin directory in operator image
	AdminToolDir = "/home/rocketmq/rocketmq-4.5.0/bin/mqadmin"

	// StoreConfigDir is the directory of config file
	StoreConfigDir = "/home/rocketmq/store/config"

	// TopicJsonDir is the directory of topics.json
	TopicJsonDir = "/home/rocketmq/store/config/topics.json"

	// SubscriptionGroupJsonDir is the directory of subscriptionGroup.json
	SubscriptionGroupJsonDir = "/home/rocketmq/store/config/subscriptionGroup.json"

	// UpdateIoTDBConfig is update broker config command
	UpdateIoTDBConfig = "updateIoTDBConfig"

	// ParamNameServiceAddress is the name of name server list parameter
	ParamNameServiceAddress = "namesrvAddr"

	// EnvNameServiceAddress is the container environment variable name of name server list
	EnvNameServiceAddress = "NAMESRV_ADDR"

	// EnvReplicationMode is the container environment variable name of replication mode
	EnvReplicationMode = "REPLICATION_MODE"

	// EnvIoTDBId is the container environment variable name of broker id
	EnvIoTDBId = "BROKER_ID"

	// EnvIoTDBClusterName is the container environment variable name of broker cluster name
	EnvIoTDBClusterName = "BROKER_CLUSTER_NAME"

	// EnvIoTDBName is the container environment variable name of broker name
	EnvIoTDBName = "BROKER_NAME"

	// LogMountPath is the directory of RocketMQ log files
	LogMountPath = "/home/iotdb/iotdb-0.9.0/server/logs"

	// StoreMountPath is the directory of RocketMQ store files
	StoreMountPath = "/home/iotdb/iotdb-0.9.0/server/data"

	// LogSubPathName is the sub-path name of log dir under mounted host dir
	LogSubPathName = "logs"

	// StoreSubPathName is the sub-path name of store dir under mounted host dir
	StoreSubPathName = "data"

	// NameServiceMainContainerPort is the main port number of name server container
	NameServiceMainContainerPort = 9876

	// NameServiceMainContainerPortName is the main port name of name server container
	NameServiceMainContainerPortName = "main"

	// JdbcPort is the VIP port number of broker container
	JdbcPort = 6667

	// JdbcPortName is the VIP port name of broker container
	JdbcPortName = "jdbc"

	// JMXPort is the main port number of broker container
	JMXPort = 31999

	// JMXPortName is the main port name of broker container
	JMXPortName = "jmx"

	// IoTDBHighAvailabilityContainerPort is the high availability port number of broker container
	IoTDBHighAvailabilityContainerPort = 10912

	// IoTDBHighAvailabilityContainerPortName is the high availability port name of broker container
	IoTDBHighAvailabilityContainerPortName = "ha"

	// StorageModeNFS is the name of NFS storage mode
	StorageModeNFS = "NFS"

	// StorageModeEmptyDir is the name of EmptyDir storage mode
	StorageModeEmptyDir = "EmptyDir"

	// StorageModeHostPath is the name pf HostPath storage mode
	StorageModeHostPath = "HostPath"

	// RestartIoTDBPodIntervalInSecond is restart broker pod interval in second
	RestartIoTDBPodIntervalInSecond = 30

	// MinMetadataJsonFileSize is the threshold value if file length is lower than this will be considered as invalid
	MinMetadataJsonFileSize = 5

	// MinIpListLength is the threshold value if the name server list parameter length is shorter than this will be considered as invalid
	MinIpListLength = 8

	// CheckConsumeFinishIntervalInSecond is the interval of checking whether the consumption process is finished in second
	CheckConsumeFinishIntervalInSecond = 5

	// RequeueIntervalInSecond is an universal interval of the reconcile function
	RequeueIntervalInSecond = 6

	// Topic is the topic field index of the output when using command check consume progress
	Topic = 0

	// IoTDBName is the broker name field index of the output when using command check consume progress
	IoTDBName = 1

	// Diff is the diff field index of the output when using command check consume progress
	Diff = 6

	// TopicListTopic is the topic field index of the output when using command check topic list
	TopicListTopic = 1

	// TopicListConsumerGroup is the consumer group field index of the output when using command check topic list
	TopicListConsumerGroup = 2
)
