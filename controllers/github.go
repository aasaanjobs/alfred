package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/aasaanjobs/aj-alfred-ci/build"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"

	"github.com/aasaanjobs/aj-alfred-ci/utils"
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
		"committer": response.Commits[0].Committer.Name,
	})
	log.Info("Received request")
	go BuildAndDeploy(response, log)
	return nil
}

// Ping is the controller for healthcheck
func Ping(w http.ResponseWriter, r *http.Request) error {
	return nil
}
