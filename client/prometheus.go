package main

import (
	"fmt"
	rbacv1beta1 "k8s.io/api/rbac/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"os"
	"path/filepath"
)

func main() {
	// Create client
	var kubeconfig string
	kubeconfig, ok := os.LookupEnv("KUBECONFIG")
	if !ok {
		kubeconfig = filepath.Join(homedir.HomeDir(), ".kube", "config")
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	kubeclient := clientset.RbacV1beta1().ClusterRoles()

	// Create resource object
	object := &rbacv1beta1.ClusterRole{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ClusterRole",
			APIVersion: "rbac.authorization.k8s.io/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "prometheus",
		},
		Rules: []rbacv1beta1.PolicyRule{
			rbacv1beta1.PolicyRule{
				Verbs: []string{
					"get",
					"list",
					"watch",
				},
				APIGroups: []string{},
				Resources: []string{
					"nodes",
					"nodes/proxy",
					"services",
					"endpoints",
					"pods",
				},
			},
			rbacv1beta1.PolicyRule{
				Verbs: []string{
					"get",
					"list",
					"watch",
				},
				APIGroups: []string{
					"extensions",
				},
				Resources: []string{
					"ingresses",
				},
			},
			rbacv1beta1.PolicyRule{
				Verbs: []string{
					"get",
				},
				NonResourceURLs: []string{
					"/metrics",
				},
			},
		},
	}

	// Manage resource
	_, err = kubeclient.Create(object)
	if err != nil {
		panic(err)
	}
	fmt.Println("ClusterRole Created successfully!")
}
