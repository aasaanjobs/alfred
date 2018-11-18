package build

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"

	"github.com/sirupsen/logrus"

	"gopkg.in/src-d/go-git.v4/plumbing"

	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"

	"gopkg.in/src-d/go-git.v4"
	appsv1 "k8s.io/api/apps/v1"
)

func retrieveDeployment(targetDir string, log *logrus.Entry) (*appsv1.Deployment, error) {
	log.Info("Reading deployment.json")
	jsonFile, err := os.Open(path.Join(targetDir, "deployment.json"))
	if err != nil {
		log.Errorf("Failed to read deployment.json, reason: %s", err.Error())
		return nil, err
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var deployment *appsv1.Deployment
	json.Unmarshal([]byte(byteValue), &deployment)
	return deployment, nil
}

func clone(url, targetDir, branch string, logger *logrus.Entry) (*git.Repository, error) {
	logger.Infof("Cloning repository %s to %s", url, targetDir)
	repo, err := git.PlainClone(targetDir, false, &git.CloneOptions{
		URL:               url,
		Progress:          os.Stdout,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		Auth: &http.BasicAuth{
			Username: "soheltarir",
			Password: "Blackhole@1719",
		},
		ReferenceName: plumbing.ReferenceName(branch),
		SingleBranch:  true,
	})
	if err != nil {
		logger.Errorf("Repostory clone failed")
		return nil, err
	}
	return repo, nil
}
