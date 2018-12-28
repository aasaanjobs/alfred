package jira

import (
	"log"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"

	"github.com/andygrunwald/go-jira"
)

// Settings represents the Jira configuration
type Settings struct {
	URL            string `env:"JIRA_SITE_URL,required"`
	Username       string `env:"JIRA_ADMIN_USER,required"`
	Password       string `env:"JIRA_ADMIN_TOKEN,required"`
	ExposedIPField string `env:"JIRA_EXPOSED_IP_FIELD,required"`
}

var jiraClient *jira.Client
var settings Settings

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	config := Settings{}
	if err := env.Parse(&config); err != nil {
		panic(err)
	}
	log.Printf("Initialised Jira Config: %+v\n", config)
	settings = config
	tp := jira.BasicAuthTransport{
		Username: config.Username,
		Password: config.Password,
	}
	client, err := jira.NewClient(tp.Client(), config.URL)
	if err != nil {
		panic(err)
	}
	jiraClient = client
}

// GetJiraClient returns a handler to JIRA client
func GetJiraClient() *jira.Client {
	return jiraClient
}

// GetJiraConfig returns a handler to the JIRA settings
func GetJiraConfig() Settings {
	return settings
}
