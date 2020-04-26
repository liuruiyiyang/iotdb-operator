package iotdb

import (
	"context"
	"encoding/json"
	iotdbv1alpha1 "github.com/apache/iotdb-operator/pkg/apis/iotdb/v1alpha1"
	cons "github.com/apache/iotdb-operator/pkg/constants"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"strconv"
	"time"
)

var log = logf.Log.WithName("controller_iotdb")

type TimeSegmentDescriptor struct {
	// Time string of time segment start time
	ScaleTime string `json:"scaleTime"`
	// The schema segment instance number in this time segment
	SchemaSegments []SchemaSegmentDescriptor `json:"schemaSplit"`
	//
}

type SchemaSegmentDescriptor struct {

	WriteRead string `json:"writeRead"`

	ReadOnly []string `json:"readOnly"`

}

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new IoTDB Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileIoTDB{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("iotdb-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource IoTDB
	err = c.Watch(&source.Kind{Type: &iotdbv1alpha1.IoTDB{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner IoTDB
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &iotdbv1alpha1.IoTDB{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileIoTDB implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileIoTDB{}

// ReconcileIoTDB reconciles a IoTDB object
type ReconcileIoTDB struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// labelsForIoTDB returns the labels for selecting the resources
func labelsForIoTDB(name string) map[string]string {
	return map[string]string{"app": "iotdb", "iotdb_cr": name}
}

func labelsForIkr(name string) map[string]string {
	return map[string]string{"app": "ikr", "ikr_cr": name}
}

var cmd = []string{"/bin/bash", "-c", "echo Initial iotdb server"}

func (r *ReconcileIoTDB) getIKRStatefulSet(iotdbCluster *iotdbv1alpha1.IoTDB, deployDescriptor string) *appsv1.StatefulSet {
	var statefulSetName string
	statefulSetName = iotdbCluster.Name + "-ikr"
	ls := labelsForIkr(statefulSetName)

	dep := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      statefulSetName,
			Namespace: iotdbCluster.Namespace,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: &iotdbCluster.Spec.IkrSize,
			Selector: &metav1.LabelSelector{
				MatchLabels: ls,
			},
			UpdateStrategy: appsv1.StatefulSetUpdateStrategy{
				Type: appsv1.RollingUpdateStatefulSetStrategyType,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: ls,
				},
				Spec: corev1.PodSpec{
					//HostNetwork: true,
					//DNSPolicy:   "ClusterFirstWithHostNet",
					Containers: []corev1.Container{{
						Image: iotdbCluster.Spec.IkrImage,
						Name:  cons.IKRContainerName,
						Lifecycle: &corev1.Lifecycle{
							PostStart: &corev1.Handler{
								Exec: &corev1.ExecAction{
									Command: cmd,
								},
							},
						},
						ImagePullPolicy: iotdbCluster.Spec.ImagePullPolicy,
						Env: []corev1.EnvVar{
							{
								Name:  cons.DeploymentDescriptor,
								Value: deployDescriptor,
							}},
						Ports: []corev1.ContainerPort{{
							ContainerPort: cons.IkrPort,
							Name:          cons.IkrPortName,
						}, {
							ContainerPort: cons.IkrJmxPort,
							Name:          cons.IkrJmxPortName,
						}},
						VolumeMounts: []corev1.VolumeMount{{
							MountPath: "/home/ikr/ikr-0.9.0/logs",
							Name:      iotdbCluster.Spec.VolumeClaimTemplates[0].Name,
							SubPath:   cons.LogSubPathName + "-ikr",
						}},
					}},
					Volumes: getVolumes(iotdbCluster),
				},
			},
			VolumeClaimTemplates: getVolumeClaimTemplates(iotdbCluster),
		},
	}
	// Set IoTDB instance as the owner and controller
	controllerutil.SetControllerReference(iotdbCluster, dep, r.scheme)

	return dep

}

func (r *ReconcileIoTDB) getIoTDBStatefulSet(iotdbCluster *iotdbv1alpha1.IoTDB, timeSegmentIndex int, schemaSegmentIndex int, replicaIndex int, syncReceiverAddress string) *appsv1.StatefulSet {
	var ls map[string]string
	var a int32 = 1
	var c = &a
	var statefulSetName string
	var isMaster string
	if replicaIndex == 0 {
		isMaster = "true"
		ls = labelsForIoTDB(iotdbCluster.Name+"-master")
		statefulSetName = iotdbCluster.Name + "-" + strconv.Itoa(timeSegmentIndex) + "-" + strconv.Itoa(schemaSegmentIndex) + "-master"
	} else {
		isMaster = "false"
		ls = labelsForIoTDB(iotdbCluster.Name+"-replica")
		statefulSetName = iotdbCluster.Name + "-" + strconv.Itoa(timeSegmentIndex) + "-" + strconv.Itoa(schemaSegmentIndex) + "-replica-" + strconv.Itoa(replicaIndex)
	}

	dep := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      statefulSetName,
			Namespace: iotdbCluster.Namespace,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: c,
			Selector: &metav1.LabelSelector{
				MatchLabels: ls,
			},
			UpdateStrategy: appsv1.StatefulSetUpdateStrategy{
				Type: appsv1.RollingUpdateStatefulSetStrategyType,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: ls,
				},
				Spec: corev1.PodSpec{
					//HostNetwork: true,
					//DNSPolicy:   "ClusterFirstWithHostNet",
					Containers: []corev1.Container{{
						Image: iotdbCluster.Spec.IotdbImage,
						Name:  cons.IoTDBContainerName,
						Lifecycle: &corev1.Lifecycle{
							PostStart: &corev1.Handler{
								Exec: &corev1.ExecAction{
									Command: cmd,
								},
							},
						},
						ImagePullPolicy: iotdbCluster.Spec.ImagePullPolicy,
						Env: []corev1.EnvVar{
							{
							Name:  cons.DeploymentDescriptor,
							Value: strconv.Itoa(replicaIndex),
						}, {
							Name:  cons.EnvIoTDBSyncReceiverAddress,
							Value: syncReceiverAddress,
						}, {
							Name:  cons.EnvIoTDBIsMaster,
							Value: isMaster,
						}},
						Ports: []corev1.ContainerPort{{
							ContainerPort: cons.JdbcPort,
							Name:          cons.JdbcPortName,
						}, {
							ContainerPort: cons.JMXPort,
							Name:          cons.JMXPortName,
						}, {
							ContainerPort: cons.SyncPort,
							Name:          cons.SyncPortName,
						}},
						VolumeMounts: []corev1.VolumeMount{{
							MountPath: cons.LogMountPath,
							Name:      iotdbCluster.Spec.VolumeClaimTemplates[0].Name,
							SubPath:   cons.LogSubPathName + getPathSuffix(iotdbCluster, timeSegmentIndex, schemaSegmentIndex, replicaIndex),
						}, {
							MountPath: cons.StoreMountPath,
							Name:      iotdbCluster.Spec.VolumeClaimTemplates[0].Name,
							SubPath:   cons.StoreSubPathName + getPathSuffix(iotdbCluster, timeSegmentIndex, schemaSegmentIndex, replicaIndex),
						}},
					}},
					Volumes: getVolumes(iotdbCluster),
				},
			},
			VolumeClaimTemplates: getVolumeClaimTemplates(iotdbCluster),
		},
	}
	// Set IoTDB instance as the owner and controller
	controllerutil.SetControllerReference(iotdbCluster, dep, r.scheme)

	return dep

}

func getPathSuffix(iotdbCluster *iotdbv1alpha1.IoTDB, timeSegmentIndex int, schemaSegmentIndex int, replicaIndex int) string {
	if replicaIndex == 0 {
		return "/" + iotdbCluster.Name + "-" + strconv.Itoa(timeSegmentIndex) + "-" + strconv.Itoa(schemaSegmentIndex) + "-master"
	}
	return "/" + iotdbCluster.Name + "-" + strconv.Itoa(timeSegmentIndex) + "-" + strconv.Itoa(schemaSegmentIndex) + "-replica-" + strconv.Itoa(replicaIndex)
}

func getVolumeClaimTemplates(server *iotdbv1alpha1.IoTDB) []corev1.PersistentVolumeClaim {
	switch server.Spec.StorageMode {
	case cons.StorageModeNFS:
		return server.Spec.VolumeClaimTemplates
	case cons.StorageModeEmptyDir, cons.StorageModeHostPath:
		fallthrough
	default:
		return nil
	}
}

func getVolumes(iotdb *iotdbv1alpha1.IoTDB) []corev1.Volume {
	switch iotdb.Spec.StorageMode {
	case cons.StorageModeNFS:
		return nil
	case cons.StorageModeEmptyDir:
		volumes := []corev1.Volume{{
			Name: iotdb.Spec.VolumeClaimTemplates[0].Name,
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{}},
		}}
		return volumes
	case cons.StorageModeHostPath:
		fallthrough
	default:
		volumes := []corev1.Volume{{
			Name: iotdb.Spec.VolumeClaimTemplates[0].Name,
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: iotdb.Spec.HostPath,
				}},
		}}
		return volumes
	}
}

var isInitializing = true

// Reconcile reads that state of the cluster for a IoTDB object and makes changes based on the state read
// and what is in the IoTDB.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileIoTDB) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling IoTDB")

	// Fetch the IoTDB iotdbCluster
	iotdbCluster := &iotdbv1alpha1.IoTDB{}
	err := r.client.Get(context.TODO(), request.NamespacedName, iotdbCluster)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	iotdbCluster.Status.TimeSegmentNum = len(iotdbCluster.Spec.Sharding)

	// create replica istances
	for timeSegmentIndex := 0; timeSegmentIndex < iotdbCluster.Status.TimeSegmentNum; timeSegmentIndex++ {
		thisSegmentNum := iotdbCluster.Spec.Sharding[timeSegmentIndex].SchemaSegmentNum
		for segmentIndex := 0; segmentIndex < thisSegmentNum; segmentIndex++ {
			for replicaIndex := 1; replicaIndex <= iotdbCluster.Spec.ReplicaNum; replicaIndex++ {
				reqLogger.Info("Check Replica IoTDB instance", "timeSegmentIndex", strconv.Itoa(timeSegmentIndex+1), "schemaSegmentIndex", strconv.Itoa(segmentIndex+1), "replicaIndex", replicaIndex)
				dep := r.getIoTDBStatefulSet(iotdbCluster, timeSegmentIndex, segmentIndex, replicaIndex, "")
				// Check if the statefulSet already exists, if not create a new one
				found := &appsv1.StatefulSet{}
				err = r.client.Get(context.TODO(), types.NamespacedName{Name: dep.Name, Namespace: dep.Namespace}, found)
				if err != nil && errors.IsNotFound(err) {
					reqLogger.Info("Creating a new Replica IoTDB StatefulSet.", "StatefulSet.Namespace", dep.Namespace, "StatefulSet.Name", dep.Name)
					err = r.client.Create(context.TODO(), dep)
					if err != nil {
						reqLogger.Error(err, "Failed to create new StatefulSet", "StatefulSet.Namespace", dep.Namespace, "StatefulSet.Name", dep.Name)
					}
				} else if err != nil {
					reqLogger.Error(err, "Failed to get StatefulSet.")
				}
			}
		}
	}

	// generate deploymentDescriptor for IKR
	replicaPodList := &corev1.PodList{}
	labelSelector := labels.SelectorFromSet(labelsForIoTDB(iotdbCluster.Name+"-replica"))
	listOps := &client.ListOptions{
		Namespace:     iotdbCluster.Namespace,
		LabelSelector: labelSelector,
	}

	err = r.client.List(context.TODO(), listOps, replicaPodList)
	if err != nil {
		reqLogger.Error(err, "Failed to list pods.", "IoTDB.Namespace", iotdbCluster.Namespace, "IoTDB.Name", iotdbCluster.Name)
		return reconcile.Result{Requeue: true}, err
	}


	notReady := false
	for _, pod := range replicaPodList.Items {
		if !reflect.DeepEqual(pod.Status.Phase, corev1.PodRunning) {
			log.Info("pod " + pod.Name + " phase is " + string(pod.Status.Phase) + ", wait for a moment...")
			notReady = true
		}
	}
	if !notReady {
		isInitializing = false
	}

	if !isInitializing {
		replicaHostIps := getIoTDBHostIps(replicaPodList.Items)
		// deploy master instances
		readOnlyIndex := 0
		for timeSegmentIndex := 0; timeSegmentIndex < iotdbCluster.Status.TimeSegmentNum; timeSegmentIndex++ {
			thisSegmentNum := iotdbCluster.Spec.Sharding[timeSegmentIndex].SchemaSegmentNum
			for segmentIndex := 0; segmentIndex < thisSegmentNum; segmentIndex++ {
				var currentSyncReceiverAddress string
				for i := 0;i<iotdbCluster.Spec.ReplicaNum; i++{
					currentSyncReceiverAddress = currentSyncReceiverAddress + replicaHostIps[readOnlyIndex+i] + " "
				}
				reqLogger.Info("Check IoTDB RW time-schema segment " + strconv.Itoa(timeSegmentIndex+1) + "-" + strconv.Itoa(segmentIndex+1) + "/" + strconv.Itoa(iotdbCluster.Status.TimeSegmentNum) + "-" + strconv.Itoa(thisSegmentNum))
				dep := r.getIoTDBStatefulSet(iotdbCluster, timeSegmentIndex, segmentIndex, 0, currentSyncReceiverAddress)
				// Check if the statefulSet already exists, if not create a new one
				found := &appsv1.StatefulSet{}
				err = r.client.Get(context.TODO(), types.NamespacedName{Name: dep.Name, Namespace: dep.Namespace}, found)
				if err != nil && errors.IsNotFound(err) {
					reqLogger.Info("Creating a new Master IoTDB StatefulSet.", "StatefulSet.Namespace", dep.Namespace, "StatefulSet.Name", dep.Name)
					err = r.client.Create(context.TODO(), dep)
					if err != nil {
						reqLogger.Error(err, "Failed to create new StatefulSet", "StatefulSet.Namespace", dep.Namespace, "StatefulSet.Name", dep.Name)
					}
				} else if err != nil {
					reqLogger.Error(err, "Failed to get StatefulSet.")
				}
				readOnlyIndex += iotdbCluster.Spec.ReplicaNum
			}
		}


		masterPodList := &corev1.PodList{}
		labelSelector := labels.SelectorFromSet(labelsForIoTDB(iotdbCluster.Name+"-master"))
		listOps := &client.ListOptions{
			Namespace:     iotdbCluster.Namespace,
			LabelSelector: labelSelector,
		}

		err = r.client.List(context.TODO(), listOps, masterPodList)
		if err != nil {
			reqLogger.Error(err, "Failed to list pods.", "IoTDB.Namespace", iotdbCluster.Namespace, "IoTDB.Name", iotdbCluster.Name)
			return reconcile.Result{Requeue: true}, err
		}


		notReady := false
		isMasterInitializing := true
		for _, pod := range masterPodList.Items {
			if !reflect.DeepEqual(pod.Status.Phase, corev1.PodRunning) {
				log.Info("pod " + pod.Name + " phase is " + string(pod.Status.Phase) + ", wait for a moment...")
				notReady = true
			}
		}
		if !notReady {
			isMasterInitializing = false
		}

		if !isMasterInitializing {
			masterHostIps := getIoTDBHostIps(masterPodList.Items)
			deploymentDescriptor := make([]TimeSegmentDescriptor, iotdbCluster.Status.TimeSegmentNum)
			writeReadInstanceIndex := 0
			readOnlyInstanceIndex := 0
			for timeSegmentIndex := 0; timeSegmentIndex < iotdbCluster.Status.TimeSegmentNum; timeSegmentIndex++ {
				deploymentDescriptor[timeSegmentIndex].ScaleTime = iotdbCluster.Spec.Sharding[timeSegmentIndex].ScaleTime
				thisSegmentNum := iotdbCluster.Spec.Sharding[timeSegmentIndex].SchemaSegmentNum
				schemaSegments := make([]SchemaSegmentDescriptor, thisSegmentNum)
				for segmentIndex := 0; segmentIndex < thisSegmentNum; segmentIndex++ {
					var readOnlyList []string
					if iotdbCluster.Spec.ReplicaNum == 0 {
						readOnlyList = make([]string, 0)
					} else {
						for i:=0;i<iotdbCluster.Spec.ReplicaNum;i++{
							readOnlyList = append(readOnlyList, replicaHostIps[readOnlyInstanceIndex + i] + ":6667")
						}
					}
					schemaSegments[segmentIndex].ReadOnly = readOnlyList
					schemaSegments[segmentIndex].WriteRead = masterHostIps[writeReadInstanceIndex] + ":6667"
					writeReadInstanceIndex++
					readOnlyInstanceIndex += iotdbCluster.Spec.ReplicaNum
				}
				deploymentDescriptor[timeSegmentIndex].SchemaSegments = schemaSegments
			}
			deploymentDescriptorJson, err := json.Marshal(deploymentDescriptor)
			if err != nil {
				reqLogger.Error(err, "Failed to generate deployment descriptor JSON", "IoTDB.Namespace", iotdbCluster.Namespace, "IoTDB.Name", iotdbCluster.Name)
				return reconcile.Result{Requeue: true}, err
			}
			reqLogger.Info("Generated deploymentDescriptorJson: " + string(deploymentDescriptorJson))

			found := &appsv1.StatefulSet{}

			dep := r.getIKRStatefulSet(iotdbCluster, string(deploymentDescriptorJson))
			err = r.client.Get(context.TODO(), types.NamespacedName{Name: dep.Name, Namespace: dep.Namespace}, found)
			if err != nil && errors.IsNotFound(err) {
				reqLogger.Info("Creating a new IKR StatefulSet.", "StatefulSet.Namespace", dep.Namespace, "StatefulSet.Name", dep.Name)
				err = r.client.Create(context.TODO(), dep)
				reqLogger.Info("Created a new IKR StatefulSet.", "StatefulSet.Namespace", dep.Namespace, "StatefulSet.Name", dep.Name)
				if err != nil {
					reqLogger.Error(err, "Failed to create StatefulSet of IKR", "StatefulSet.Namespace", dep.Namespace, "StatefulSet.Name", dep.Name)
				}
				// StatefulSet created successfully - return and requeue
				return reconcile.Result{Requeue: true}, nil
			} else if err != nil {
				reqLogger.Error(err, "Failed to get IKR StatefulSet.")
			}
		}
	}

	return reconcile.Result{Requeue: true, RequeueAfter: time.Duration(cons.RequeueIntervalInSecond) * time.Second}, nil
}

func getIoTDBHostIps(pods []corev1.Pod) []string {
	var ips []string
	for _, pod := range pods {
		ips = append(ips, pod.Status.PodIP)
	}
	return ips
}

// newPodForCR returns a busybox pod with the same name/namespace as the cr
func newPodForCR(cr *iotdbv1alpha1.IoTDB) *corev1.Pod {
	labels := map[string]string{
		"app": cr.Name,
	}
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-pod",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    "busybox",
					Image:   "busybox",
					Command: []string{"sleep", "3600"},
				},
			},
		},
	}
}
