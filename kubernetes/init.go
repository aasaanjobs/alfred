package kubernetes

import (
	"flag"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var k8sClient *kubernetes.Clientset

func init() {
	var kubeconfig *string
	home := homedir.HomeDir()
	kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	flag.Parse()
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	k8sClient = clientset
}

// GetK8SClient returns the handle to the configured kubernetes client
func GetK8SClient() *kubernetes.Clientset {
	return k8sClient
}
