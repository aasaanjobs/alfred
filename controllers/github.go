package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aasaanjobs/alfred/jira"

	k "github.com/aasaanjobs/alfred/kubernetes"

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

func buildAndDeploy(response utils.WebHookPayload, reDeploy bool, logger *logrus.Entry) {
	buildResponse, err := build.Build(response.Repository, response.Ref, logger)
	if err != nil {
		logger.Errorf("Build failed, reason: %s", err.Error())
		return
	}
	logger.Infof("Build Successful")

	k.DeployWorkload(buildResponse, reDeploy, logger)
}

func checkDeployCommand(commits []utils.Commit) (int, bool) {
	for index, commit := range commits {
		if utils.CommitHasDeployCmd(commit.Message) {
			return index, true
		}
	}
	return -1, false
}

// WebHook is the controller to handle github push webhooks
func WebHook(w http.ResponseWriter, r *http.Request) error {
	githubEvent := r.Header.Get("X-GitHub-Event")
	if githubEvent == "ping" {
		printer.Info("Received ping event from github, exiting.")
		return nil
	}
	var response utils.WebHookPayload
	if err := json.NewDecoder(r.Body).Decode(&response); err != nil {
		return err
	}
	printer.Info("Received Push from Github.", logrus.Fields{
		"project": response.Repository.Name, "ref": response.Ref,
	})

	// Check if the branch is already deployed or not, else check for explicit deploy command
	// in git commit message
	deploymentName := utils.GetDeploymentName(response.Repository.Name, response.Ref)
	var reDeploy bool
	if _, err := k.GetDeployment(deploymentName); err == nil {
		reDeploy = true
	}
	commitIndex, shouldDeploy := checkDeployCommand(response.Commits)
	if !shouldDeploy {
		printer.Info("No explicit deploy command received, exiting.")
		return nil
	}
	printer.Infof("Found deploy command in commit: %s", response.Commits[commitIndex].Message)

	// Initialise the logger
	log := printer.WithFields(logrus.Fields{
		"project":   response.Repository.Name,
		"ref":       response.Ref,
		"committer": response.Commits[commitIndex].Committer.Email,
	})
	log.Info("Received build and deploy request...")
	go buildAndDeploy(response, reDeploy, log)
	return nil
}

// Ping is the controller for healthcheck
func Ping(w http.ResponseWriter, r *http.Request) error {
	settings := jira.GetJiraConfig()
	fmt.Println(settings)
	return nil
}
