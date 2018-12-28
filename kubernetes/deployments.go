package kubernetes

import (
	"fmt"

	"github.com/aasaanjobs/alfred/jira"
	"github.com/aasaanjobs/alfred/utils"

	"github.com/sirupsen/logrus"

	"github.com/aasaanjobs/alfred/build"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
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

func modifyDeployment(build *build.Response) (*appsv1.Deployment, string) {
	var deployment = build.Deployment
	deployment.ObjectMeta.Name = build.NewDeploymentName()
	deployment.Spec.Selector.MatchLabels = map[string]string{
		"app": build.NewDeploymentName(),
	}
	deployment.Spec.Template.ObjectMeta.Labels = map[string]string{
		"app": build.NewDeploymentName(),
	}
	for i, image := range build.Images {
		deployment.Spec.Template.Spec.Containers[i].Image = image
	}
	return deployment, deployment.ObjectMeta.Name
}

// DeployWorkload deploys the build image to kubernetes cluster
func DeployWorkload(build *build.Response, reDeploy bool, logger *logrus.Entry) error {
	client := GetK8SClient().AppsV1().Deployments(apiv1.NamespaceDefault)
	workload, workloadName := modifyDeployment(build)
	if !reDeploy {
		// New Deployment
		logger.Infof("Deploying new workload %s", workloadName)
		if _, err := client.Create(workload); err != nil {
			logger.Errorf("Failed to deploy workload %s, reason: %s", workloadName, err.Error())
			return err
		}
		logger.Infof("Successfully deployed workload %s", workloadName)
		exposedIP, err := ExposeService(workloadName, build.Service, logger)
		if err != nil {
			return err
		}
		if err := jira.UpdateIssue(exposedIP, utils.GetJiraID(build.Feature), logger); err != nil {
			return err
		}
		return nil
	}
	// Existing Deployment
	logger.Infof("Updating workload %s", workloadName)
	if _, err := client.Update(workload); err != nil {
		logger.Errorf("Failed to update workload %s, reason: %s", workloadName, err.Error())
		return err
	}
	logger.Infof("Successfully updated workload %s", workloadName)
	return nil
}
