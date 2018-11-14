package utils

// Committer is the one who submitted the commit message
type Committer struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Commit is the commit object
type Commit struct {
	Message   string    `json:"message"`
	Committer Committer `json:"committer"`
}

// Repository is the git repo object
type Repository struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	CloneURL string `json:"clone_url"`
}

// WebHookPayload is the json body that is sent by github
type WebHookPayload struct {
	Ref        string     `json:"ref"`
	Repository Repository `json:"repository"`
	Commits    []Commit   `json:"commits"`
}
