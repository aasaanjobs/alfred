package jira

import (
	"fmt"

	"github.com/andygrunwald/go-jira"
	"github.com/sirupsen/logrus"
	"github.com/trivago/tgo/tcontainer"
)

// UpdateIssue updates the JIRA card with the exposed IP
func UpdateIssue(exposedIP string, jiraID string, logger *logrus.Entry) error {
	issue := jira.Issue{
		Key: jiraID,
		Fields: &jira.IssueFields{
			Unknowns: tcontainer.MarshalMap{
				settings.ExposedIPField: fmt.Sprintf("http://%s", exposedIP),
			},
		},
	}
	if _, _, err := GetJiraClient().Issue.Update(&issue); err != nil {
		logger.Errorf("Failed to update JIRA issue %s with exposed IP, reason: %s",
			jiraID,
			err.Error(),
		)
		return err
	}
	if _, _, err := GetJiraClient().Issue.AddComment(jiraID, &jira.Comment{
		Body: fmt.Sprintf("Feature has been deployed and exposed to IP: %s", exposedIP),
	}); err != nil {
		return nil
	}
	return nil
}
