package utils

import (
	"strconv"
	"strings"
)

// GetFeatureName retrieves the feature name from the branch ref
func GetFeatureName(branch string) string {
	return strings.Replace(branch, "refs/heads/feature/", "", 1)
}

// GetJiraID retrieves the JIRA ID from the branch ref
func GetJiraID(branch string) string {
	featureName := GetFeatureName(branch)
	tokens := strings.Split(featureName, "-")
	// Check whether the second element is number or not
	if _, err := strconv.Atoi(tokens[1]); err != nil {
		panic("Invalid feature name")
	}
	return strings.Join(tokens[0:2], "-")
}

// GetDeploymentName returns the workload name by which the feature
// will be deployed to the Kubernetes cluster
func GetDeploymentName(project, branch string) string {
	return strings.Join([]string{
		strings.ToLower(project),
		strings.ToLower(GetJiraID(branch)),
	}, "-")
}
