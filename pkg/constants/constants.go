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
	// IoTDBContainerName is the name of iotdb container
	IoTDBContainerName = "server"

	IKRContainerName = "ikr"

	// DeploymentDescriptor is the container environment variable name of broker id
	DeploymentDescriptor = "DEPLOY_DESCRIPTOR_JSON"

	// EnvIoTDBTimeSegmentIndex is the container environment variable name of broker cluster name
	EnvIoTDBTimeSegmentIndex = "BROKER_CLUSTER_NAME"

	// EnvIoTDBSchemaSegmentIndex is the container environment variable name of broker name
	EnvIoTDBSchemaSegmentIndex = "BROKER_NAME"

	// LogMountPath is the directory of IoTDB log files
	LogMountPath = "/home/iotdb/iotdb-0.9.0/server/logs"

	// StoreMountPath is the directory of IoTDB store files
	StoreMountPath = "/home/iotdb/iotdb-0.9.0/server/data"

	// LogSubPathName is the sub-path name of log dir under mounted host dir
	LogSubPathName = "logs"

	// StoreSubPathName is the sub-path name of store dir under mounted host dir
	StoreSubPathName = "data"

	JdbcPort = 6667

	JdbcPortName = "jdbc"

	JMXPort = 31999

	JMXPortName = "jmx"

	IkrPort = 6666

	IkrPortName = "ikr"

	IkrJmxPort = 31998

	// must contain only alpha-numeric characters (a-z, 0-9), and hyphens (-)
	IkrJmxPortName = "ikr-jmx"

	// StorageModeNFS is the name of NFS storage mode
	StorageModeNFS = "NFS"

	// StorageModeEmptyDir is the name of EmptyDir storage mode
	StorageModeEmptyDir = "EmptyDir"

	// StorageModeHostPath is the name pf HostPath storage mode
	StorageModeHostPath = "HostPath"

	// RequeueIntervalInSecond is an universal interval of the reconcile function
	RequeueIntervalInSecond = 6
)
