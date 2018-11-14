package build

import (
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/aasaanjobs/aj-alfred-ci/utils"
	appsv1 "k8s.io/api/apps/v1"
)

type BuildResponse struct {
	Deployment *appsv1.Deployment
	Images     []string
	feature    string
}

// Build pulls the repository locally and builds the source using Dockerfile/s provided
func Build(repo utils.Repository, branch string, logger *logrus.Entry) (*BuildResponse, error) {
	logger.Infof("Running build sequence")
	featureName := utils.GetFeatureName(branch)
	var workDir = "/tmp"
	targetDir := fmt.Sprintf("%s/%d_%s", workDir, int32(time.Now().Unix()), featureName)
	defer os.RemoveAll(targetDir)
	if _, err := clone(repo.CloneURL, targetDir, branch, logger); err != nil {
		return nil, err
	}
	// deployment, err := retrieveDeployment(targetDir)
	// if err != nil {
	// 	return nil, err
	// }
	dockerImages, err := RunDocker(targetDir, featureName, repo.Name, logger)
	if err != nil {
		return nil, err
	}
	return &BuildResponse{
		// Deployment: deployment,
		feature: featureName,
		Images:  dockerImages,
	}, nil
}
