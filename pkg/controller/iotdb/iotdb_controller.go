package iotdb

import (
	"context"
	"strconv"
	cons "github.com/apache/iotdb-operator/pkg/constants"
	iotdbv1alpha1 "github.com/apache/iotdb-operator/pkg/apis/iotdb/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"time"
)

var log = logf.Log.WithName("controller_iotdb")

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
// belonging to the given broker CR name.
func labelsForIoTDB(name string) map[string]string {
	return map[string]string{"app": "iotdb", "iotdb_cr": name}
}

var cmd = []string{"/bin/bash", "-c", "echo Initial iotdb server"}

func (r *ReconcileIoTDB) getIoTDBStatefulSet(iotdbCluster *iotdbv1alpha1.IoTDB, iotdbGroupIndex int, replicaIndex int) *appsv1.StatefulSet {
	ls := labelsForIoTDB(iotdbCluster.Name)
	var a int32 = 1
	var c = &a
	var statefulSetName string
	if replicaIndex == 0 {
		statefulSetName = iotdbCluster.Name + "-" + strconv.Itoa(iotdbGroupIndex) + "-master"
	} else {
		statefulSetName = iotdbCluster.Name + "-" + strconv.Itoa(iotdbGroupIndex) + "-replica-" + strconv.Itoa(replicaIndex)
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
						Image: iotdbCluster.Spec.Image,
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
							Name:  cons.EnvIoTDBId,
							Value: strconv.Itoa(replicaIndex),
						}, {
							Name:  cons.EnvIoTDBClusterName,
							Value: iotdbCluster.Name + "-" + strconv.Itoa(iotdbGroupIndex),
						}, {
							Name:  cons.EnvIoTDBName,
							Value: iotdbCluster.Name + "-" + strconv.Itoa(iotdbGroupIndex),
						}},
						Ports: []corev1.ContainerPort{{
							ContainerPort: cons.JdbcPort,
							Name:          cons.JdbcPortName,
						}, {
							ContainerPort: cons.JMXPort,
							Name:          cons.JMXPortName,
						}},
						VolumeMounts: []corev1.VolumeMount{{
							MountPath: cons.LogMountPath,
							Name:      iotdbCluster.Spec.VolumeClaimTemplates[0].Name,
							SubPath:   cons.LogSubPathName + getPathSuffix(iotdbCluster, iotdbGroupIndex, replicaIndex),
						}, {
							MountPath: cons.StoreMountPath,
							Name:      iotdbCluster.Spec.VolumeClaimTemplates[0].Name,
							SubPath:   cons.StoreSubPathName + getPathSuffix(iotdbCluster, iotdbGroupIndex, replicaIndex),
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

func getPathSuffix(iotdbCluster *iotdbv1alpha1.IoTDB, iotdbGroupIndex int, replicaIndex int) string {
	if replicaIndex == 0 {
		return "/" + iotdbCluster.Name + "-" + strconv.Itoa(iotdbGroupIndex) + "-master"
	}
	return "/" + iotdbCluster.Name + "-" + strconv.Itoa(iotdbGroupIndex) + "-replica-" + strconv.Itoa(replicaIndex)
}

func getVolumeClaimTemplates(broker *iotdbv1alpha1.IoTDB) []corev1.PersistentVolumeClaim {
	switch broker.Spec.StorageMode {
	case cons.StorageModeNFS:
		return broker.Spec.VolumeClaimTemplates
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

	for iotdbGroupIndex := 0; iotdbGroupIndex < iotdbCluster.Spec.Size; iotdbGroupIndex++ {
		reqLogger.Info("Check IoTDB cluster " + strconv.Itoa(iotdbGroupIndex+1) + "/" + strconv.Itoa(iotdbCluster.Spec.Size))
		dep := r.getIoTDBStatefulSet(iotdbCluster, iotdbGroupIndex, 0)
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
			reqLogger.Error(err, "Failed to get broker master StatefulSet.")
		}
	}

	return reconcile.Result{Requeue: true, RequeueAfter: time.Duration(cons.RequeueIntervalInSecond) * time.Second}, nil
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
