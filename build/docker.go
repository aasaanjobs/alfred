package build

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"
)

// GetDockerFiles returns a list of all dockerfile configurations in the directory
func GetDockerFiles(targetDir string, logger *logrus.Entry) ([]string, error) {
	logger.Infof("Reading dockerfiles in %s", targetDir)
	files, err := ioutil.ReadDir(targetDir)
	if err != nil {
		return nil, err
	}
	var dockerFiles []string
	for _, f := range files {
		if strings.HasPrefix(f.Name(), "Dockerfile") {
			dockerFiles = append(dockerFiles, f.Name())
		}
	}
	if len(dockerFiles) == 0 {
		logger.Warnf("Failed to find Dockerfile/s in %s", targetDir)
		return nil, fmt.Errorf("No Dockerfiles in project")
	}
	return dockerFiles, nil
}

// DockerBuildImage Runs docker build for the dockerfile provided
func DockerBuildImage(workDir, dockerFile, project, tag string, logger *logrus.Entry) (string, error) {
	logger.Infof("Building docker image from %s", dockerFile)
	imageTag := fmt.Sprintf("gcr.io/aj-cloud-staging/%s:%s", project, tag)
	cmd := exec.Command("docker", "build", "-f", dockerFile, "-t", imageTag, ".")
	cmd.Dir = workDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		logger.Error("Docker build failed")
		return "", err
	}
	return imageTag, nil
}

// DockerPushImage pushes the image to google cloud registry
func DockerPushImage(image string, logger *logrus.Entry) error {
	logger.Infof("Pushing docker image %s", image)
	cmd := exec.Command("docker", "push", image)
	if err := cmd.Run(); err != nil {
		logger.Error("Failed to push docker image")
		return err
	}
	return nil
}

// RunDocker builds the images as per the steps provided in the
// Dockerfile/s and uploads the images to remote cloud registries
func RunDocker(targetDir, feature, project string, logger *logrus.Entry) ([]string, error) {
	dockerFiles, err := GetDockerFiles(targetDir, logger)
	if err != nil {
		return nil, err
	}
	// Change the working directory
	if err := os.Chdir(targetDir); err != nil {
		return nil, err
	}
	var dockerImages []string
	for _, file := range dockerFiles {
		image, err := DockerBuildImage(targetDir, file, project, feature, logger)
		if err != nil {
			return nil, err
		}
		if err := DockerPushImage(image, logger); err != nil {
			return nil, err
		}
		dockerImages = append(dockerImages, image)
	}
	return dockerImages, nil
}
