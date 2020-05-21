package aceapplication

import (
	"context"
	v1 "k8s.io/api/apps/v1"

	ibmacedraynerv1alpha1 "operators/ace-app-operator/pkg/apis/ibmacedrayner/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_aceapplication")
var CONTROLLER_VERSION = "0.1.0"

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new AceApplication Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileAceApplication{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("aceapplication-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource AceApplication
	err = c.Watch(&source.Kind{Type: &ibmacedraynerv1alpha1.AceApplication{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner AceApplication
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &ibmacedraynerv1alpha1.AceApplication{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileAceApplication implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileAceApplication{}

// ReconcileAceApplication reconciles a AceApplication object
type ReconcileAceApplication struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a AceApplication object and makes changes based on the state read
// and what is in the AceApplication.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileAceApplication) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling AceApplication")

	// Fetch the AceApplication instance
	instance := &ibmacedraynerv1alpha1.AceApplication{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
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

	// Define a new deploy object
	deploy := newDeployForCR(instance)

	// Set AceApplication instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, deploy, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if this deploy already exists
	found := &v1.Deployment{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: deploy.Name, Namespace: deploy.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Deploy", "Deploy.Namespace", deploy.Namespace, "Deploy.Name", deploy.Name)
		err = r.client.Create(context.TODO(), deploy)
		if err != nil {
			return reconcile.Result{}, err
		}

		// deploy created successfully - don't requeue
		return reconcile.Result{}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}

	// deploy already exists - don't requeue
	reqLogger.Info("Skip reconcile: deploy already exists", "Deploy.Namespace", found.Namespace, "Deploy.Name", found.Name)
	return reconcile.Result{}, nil
}

func newDeployForCR(aceApp *ibmacedraynerv1alpha1.AceApplication) *v1.Deployment {
	labels := map[string]string{
		"app": aceApp.Name,
	}
	annotations := map[string]string{
		"ace-operator/controller-version": CONTROLLER_VERSION,
	}

	labelSelector := metav1.LabelSelector{
		MatchLabels:      labels,
	}

	template := corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels:                     labels,
			Annotations:                annotations,
		},
		Spec:       corev1.PodSpec{
			Volumes:                       nil, // Precisamos criar um volume para cada item de configuração do ACE
			InitContainers:                nil,
			Containers:                    []corev1.Container{ // Terminar de Configurar este Objeto
				{
					Name:                     aceApp.Spec.AceServerName,
					Image:                    aceApp.Spec.AceBaseImage+":"+aceApp.Spec.AceImageTag,
					Command:                  nil,
					Args:                     nil,
					Ports:                    nil,
					EnvFrom:                  nil,
					Env:                      nil,
					Resources:                corev1.ResourceRequirements{},
					VolumeMounts:             nil,
					VolumeDevices:            nil,
					LivenessProbe:            nil,
					ReadinessProbe:           nil,
					StartupProbe:             nil,
					Lifecycle:                nil,
					TerminationMessagePath:   "",
					TerminationMessagePolicy: "",
					ImagePullPolicy:          "",
					SecurityContext:          nil,
					Stdin:                    false,
					StdinOnce:                false,
					TTY:                      false,
				},
			},
			RestartPolicy:                 "",
			TerminationGracePeriodSeconds: nil,
			ActiveDeadlineSeconds:         nil,
			DNSPolicy:                     "",
			NodeSelector:                  aceApp.Spec.NodeSelectorLabels,
			ServiceAccountName:            aceApp.Spec.ServiceAccountName,
		},
	}

	deployment := v1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:                       aceApp.Name,
			Namespace:                  aceApp.Namespace,
			Labels:                     labels,
			Annotations:                annotations,
		},
		Spec:       v1.DeploymentSpec{
			Replicas:                aceApp.Spec.Replicas,
			Selector:                &labelSelector,
			Template:                template,
			Strategy:                v1.DeploymentStrategy{},
			MinReadySeconds:         0,
			RevisionHistoryLimit:    nil,
			Paused:                  false,
			ProgressDeadlineSeconds: nil,
		},
	}

	return &deployment
	
}


