package config

import (
	"os"
)

var (
	//github
	AccessToken string
	Owner       string
	Repo        string
)

func init() {
	LoadWithGithub()
}

func LoadWithGithub() {
	AccessToken = getEnv("GITHUB_ACCESS_TOKEN")
	Owner = getEnv("GITHUB_OWNER")
	Repo = getEnv("GITHUB_REPO")
}

func getEnv(key string) string {
	return os.Getenv(key)
}
