package prometheus

import (
	"context"

	hyperledgerv1alpha1 "github.com/Sumaid/besu-kubernetes/besu-operator/pkg/apis/hyperledger/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (r *ReconcilePrometheus) ensureConfigMap(request reconcile.Request,
	instance *hyperledgerv1alpha1.Prometheus,
	s *corev1.ConfigMap,
) (*reconcile.Result, error) {
	found := &corev1.ConfigMap{}
	err := r.client.Get(context.TODO(), types.NamespacedName{
		Name:      s.ObjectMeta.Name,
		Namespace: s.ObjectMeta.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {
		// Create the ConfigMap
		log.Info("Creating a new ConfigMap", "ConfigMap.Namespace", s.Namespace, "ConfigMap.Name", s.Name)
		err = r.client.Create(context.TODO(), s)

		if err != nil {
			// Creation failed
			log.Error(err, "Failed to create new ConfigMap", "ConfigMap.Namespace", s.Namespace, "ConfigMap.Name", s.Name)
			return &reconcile.Result{}, err
		} else {
			// Creation was successful
			return nil, nil
		}
	} else if err != nil {
		// Error that isn't due to the ConfigMap not existing
		log.Error(err, "Failed to get ConfigMap")
		return &reconcile.Result{}, err
	}

	return nil, nil
}

func (r *ReconcilePrometheus) ensureClusterRole(request reconcile.Request,
	instance *hyperledgerv1alpha1.Prometheus,
	s *rbacv1.ClusterRole,
) (*reconcile.Result, error) {
	found := &rbacv1.ClusterRole{}
	err := r.client.Get(context.TODO(), types.NamespacedName{
		Name: s.Name,
	}, found)
	if err != nil && errors.IsNotFound(err) {

		// Create the service
		log.Info("Creating a new ClusterRole", "ClusterRole.Namespace", s.Namespace, "ClusterRole.Name", s.Name)
		err = r.client.Create(context.TODO(), s)

		if err != nil {
			// Creation failed
			log.Error(err, "Failed to create new ClusterRole", "ClusterRole.Namespace", s.Namespace, "ClusterRole.Name", s.Name)
			return &reconcile.Result{}, err
		} else {
			// Creation was successful
			return nil, nil
		}
	} else if err != nil {
		// Error that isn't due to the service not existing
		log.Error(err, "Failed to get ClusterRole")
		return &reconcile.Result{}, err
	}

	return nil, nil
}

func (r *ReconcilePrometheus) ensureServiceAccount(request reconcile.Request,
	instance *hyperledgerv1alpha1.Prometheus,
	s *corev1.ServiceAccount,
) (*reconcile.Result, error) {
	found := &corev1.ServiceAccount{}
	err := r.client.Get(context.TODO(), types.NamespacedName{
		Name:      s.ObjectMeta.Name,
		Namespace: s.ObjectMeta.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {

		// Create the service
		log.Info("Creating a new ServiceAccount", "ServiceAccount.Namespace", s.Namespace, "ServiceAccount.Name", s.Name)
		err = r.client.Create(context.TODO(), s)

		if err != nil {
			// Creation failed
			log.Error(err, "Failed to create new ServiceAccount", "ServiceAccount.Namespace", s.Namespace, "ServiceAccount.Name", s.Name)
			return &reconcile.Result{}, err
		} else {
			// Creation was successful
			return nil, nil
		}
	} else if err != nil {
		// Error that isn't due to the service not existing
		log.Error(err, "Failed to get ServiceAccount")
		return &reconcile.Result{}, err
	}

	return nil, nil
}

func (r *ReconcilePrometheus) ensureClusterRoleBinding(request reconcile.Request,
	instance *hyperledgerv1alpha1.Prometheus,
	s *rbacv1.ClusterRoleBinding,
) (*reconcile.Result, error) {
	found := &rbacv1.ClusterRoleBinding{}
	err := r.client.Get(context.TODO(), types.NamespacedName{
		Name: s.Name,
	}, found)
	if err != nil && errors.IsNotFound(err) {

		// Create the service
		log.Info("Creating a new ClusterRoleBinding", "ClusterRoleBinding.Namespace", s.Namespace, "ClusterRoleBinding.Name", s.Name)
		err = r.client.Create(context.TODO(), s)

		if err != nil {
			// Creation failed
			log.Error(err, "Failed to create new ClusterRoleBinding", "ClusterRoleBinding.Namespace", s.Namespace, "ClusterRoleBinding.Name", s.Name)
			return &reconcile.Result{}, err
		} else {
			// Creation was successful
			return nil, nil
		}
	} else if err != nil {
		// Error that isn't due to the service not existing
		log.Error(err, "Failed to get ClusterRoleBinding")
		return &reconcile.Result{}, err
	}

	return nil, nil
}

func (r *ReconcilePrometheus) ensureService(request reconcile.Request,
	instance *hyperledgerv1alpha1.Prometheus,
	s *corev1.Service,
) (*reconcile.Result, error) {
	found := &corev1.Service{}
	err := r.client.Get(context.TODO(), types.NamespacedName{
		Name:      s.ObjectMeta.Name,
		Namespace: s.ObjectMeta.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {

		// Create the service
		log.Info("Creating a new Service", "Service.Namespace", s.Namespace, "Service.Name", s.Name)
		err = r.client.Create(context.TODO(), s)

		if err != nil {
			// Creation failed
			log.Error(err, "Failed to create new Service", "Service.Namespace", s.Namespace, "Service.Name", s.Name)
			return &reconcile.Result{}, err
		} else {
			// Creation was successful
			return nil, nil
		}
	} else if err != nil {
		// Error that isn't due to the service not existing
		log.Error(err, "Failed to get Service")
		return &reconcile.Result{}, err
	}

	return nil, nil
}

func (r *ReconcilePrometheus) ensureDeployment(request reconcile.Request,
	instance *hyperledgerv1alpha1.Prometheus,
	s *appsv1.Deployment,
) (*reconcile.Result, error) {
	found := &appsv1.Deployment{}
	err := r.client.Get(context.TODO(), types.NamespacedName{
		Name:      s.ObjectMeta.Name,
		Namespace: s.ObjectMeta.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {

		// Create the service
		log.Info("Creating a new Deployment", "Deployment.Namespace", s.Namespace, "Deployment.Name", s.Name)
		err = r.client.Create(context.TODO(), s)

		if err != nil {
			// Creation failed
			log.Error(err, "Failed to create new Deployment", "Deployment.Namespace", s.Namespace, "Deployment.Name", s.Name)
			return &reconcile.Result{}, err
		} else {
			// Creation was successful
			return nil, nil
		}
	} else if err != nil {
		// Error that isn't due to the service not existing
		log.Error(err, "Failed to get Deployment")
		return &reconcile.Result{}, err
	}

	return nil, nil
}
