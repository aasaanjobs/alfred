package constants

// CommitDeployCmd is the command specified in commit message
// by which the build & deployment process should start
const CommitDeployCmd string = "#deployBranch"

// CommitNoDeployCmd is the command specified in commit message
// by which the CI should recognize not to build the feature
const CommitNoDeployCmd string = "#dontDeployBranch"

// CommitDeleteDeployCmd is the command specified in commit message
// by which the CI should recognize that the feature deployment should be deleted
const CommitDeleteDeployCmd string = "#deleteDeployedBranch"
