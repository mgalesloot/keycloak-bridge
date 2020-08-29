package keycloak

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func getKubeConfig() string {
	var kubeconfig *string
	var home, _ = os.UserHomeDir()
	if home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	return *kubeconfig
}

func getSecretData(secretName string, namespace string) map[string][]byte {

	kubeconfig := getKubeConfig()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	var keycloakSecret *v1.Secret
	keycloakSecret, err = clientset.CoreV1().Secrets(namespace).Get(context.TODO(), secretName, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		fmt.Printf("Secret %s in namespace %s not found\n", secretName, namespace)
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		fmt.Printf("Error getting secret %s in namespace %s: %v\n",
			secretName, namespace, statusError.ErrStatus.Message)
	} else if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Found secret %s in namespace %s\n", secretName, namespace)

	return keycloakSecret.Data
}
