package kubernetes

import (
	"github.com/sirupsen/logrus"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func modifyService(name string, service *apiv1.Service) *apiv1.Service {
	service.ObjectMeta.Name = name
	service.Spec.Selector = map[string]string{
		"app": name,
	}
	return service
}

// ExposeService deploys a Kubernetes service and waits for the Load Balancer external IP
func ExposeService(name string, service *apiv1.Service, logger *logrus.Entry) (string, error) {
	client := GetK8SClient().CoreV1().Services(metav1.NamespaceDefault)
	if _, err := client.Create(modifyService(name, service)); err != nil {
		logger.Errorf("Failed to expose service %s, reason: %s", name, err.Error())
		return "", err
	}
	logger.Infof("Waiting for external IP...")
	for {
		service, _ := client.Get(name, metav1.GetOptions{})
		var ingress = service.Status.LoadBalancer.Ingress
		if len(ingress) > 0 && ingress[0].IP != "" {
			logger.Infof("Received External IP %s", ingress[0].IP)
			return ingress[0].IP, nil
		}
	}
}
