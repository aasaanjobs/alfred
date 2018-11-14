package utils

import (
	"strings"

	"github.com/aasaanjobs/aj-alfred-ci/constants"
)

// CommitHasDeployCmd checks whether the user has specified the
// deploy command in commit or not
func CommitHasDeployCmd(commitMessage string) bool {
	if strings.Index(commitMessage, constants.CommitDeployCmd) >= 0 {
		return true
	}
	return false
}

// CommitHasNoDeployCmd checks whether the user has specified the
// don't deploy command in commit or not
func CommitHasNoDeployCmd(commitMessage string) bool {
	if strings.Index(commitMessage, constants.CommitNoDeployCmd) >= 0 {
		return true
	}
	return false
}

// CommitHasDeleteDeployCmd checks whether the user has specified the
// delete deployment command in commit or not
func CommitHasDeleteDeployCmd(commitMessage string) bool {
	if strings.Index(commitMessage, constants.CommitDeleteDeployCmd) >= 0 {
		return true
	}
	return false
}
