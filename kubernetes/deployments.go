package kubernetes

import (
	"fmt"

	"github.com/aasaanjobs/alfred/build"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetDeployment retrieves the kubernetes deployment if found
func GetDeployment(name string) (*appsv1.Deployment, error) {
	clientSet := GetK8SClient()
	client := clientSet.AppsV1().Deployments(metav1.NamespaceDefault)
	deployment, err := client.Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve deployment %s; reason=%s", name, err.Error())
	}
	return deployment, nil
}

func ModifyDeployment(build *build.Response) *appsv1.Deployment {
	var deployment = build.Deployment
	deployment.ObjectMeta.Name = build.NewDeploymentName()
	for i, image := range build.Images {
		deployment.Spec.Template.Spec.Containers[i].Image = image
	}
	return deployment
}
