package build

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/aasaanjobs/alfred/utils"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
)

// Response represents the build context and result
type Response struct {
	Deployment  *appsv1.Deployment
	Service     *apiv1.Service
	Images      []string
	Feature     string
	ProjectName string
}

// NewDeploymentName returns the Kubernetes workload name to assigned for the build
func (b *Response) NewDeploymentName() string {
	return fmt.Sprintf("%s-%s", b.ProjectName, strings.ToLower(utils.GetJiraID(b.Feature)))
}

// Build pulls the repository locally and builds the source using Dockerfile/s provided
func Build(repo utils.Repository, branch string, logger *logrus.Entry) (*Response, error) {
	logger.Infof("Running build sequence")
	featureName := utils.GetFeatureName(branch)
	var workDir = "/tmp"
	targetDir := fmt.Sprintf("%s/%d_%s", workDir, int32(time.Now().Unix()), featureName)
	defer os.RemoveAll(targetDir)
	if _, err := clone(repo.CloneURL, targetDir, branch, logger); err != nil {
		return nil, err
	}
	deployment, err := retrieveDeployment(targetDir, logger)
	if err != nil {
		return nil, err
	}
	service, err := retrieveService(targetDir, logger)
	if err != nil {
		return nil, err
	}
	dockerImages, err := RunDocker(targetDir, featureName, repo.Name, logger)
	if err != nil {
		return nil, err
	}
	return &Response{
		Deployment:  deployment,
		Feature:     featureName,
		ProjectName: repo.Name,
		Images:      dockerImages,
		Service:     service,
	}, nil
}
