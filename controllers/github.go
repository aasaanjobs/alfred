package controllers

import (
	"encoding/json"
	"net/http"

	k "github.com/aasaanjobs/alfred/kubernetes"
	apiv1 "k8s.io/api/core/v1"

	"github.com/aasaanjobs/alfred/build"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"

	"github.com/aasaanjobs/alfred/utils"
)

var printer *logrus.Logger

func init() {
	printer = logrus.New()
	// Set the formatter
	printer.Formatter = &prefixed.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	}
}

func BuildAndDeploy(response utils.WebHookPayload, logger *logrus.Entry) {
	// Initialise the logger
	buildResponse, err := build.Build(response.Repository, response.Ref, logger)
	if err != nil {
		logger.Errorf("Build failed, reason: %s", err.Error())
	} else {
		logger.Infof("Build Successful, %q", buildResponse.Deployment)
	}
	deploymentsClient := k.GetK8SClient().AppsV1().Deployments(apiv1.NamespaceDefault)
	if _, err := deploymentsClient.Create(k.ModifyDeployment(buildResponse)); err != nil {
		logger.Errorf("Deployment failed, reason: %s", err.Error())
	} else {
		logger.Infof("Deployment Successful")
	}
}

// WebHook is the controller to handle github push webhooks
func WebHook(w http.ResponseWriter, r *http.Request) error {
	var response utils.WebHookPayload
	if err := json.NewDecoder(r.Body).Decode(&response); err != nil {
		return err
	}
	// Initialise the logger
	log := printer.WithFields(logrus.Fields{
		"project":   response.Repository.Name,
		"ref":       response.Ref,
		"committer": response.Commits[0].Committer.Email,
	})
	log.Info("Received build and deploy request...")
	go BuildAndDeploy(response, log)
	return nil
}

// Ping is the controller for healthcheck
func Ping(w http.ResponseWriter, r *http.Request) error {
	return nil
}
