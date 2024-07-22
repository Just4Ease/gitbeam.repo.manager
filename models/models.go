package models

type Repo struct {
	TimeCreated   string `json:"timeCreated"`
	TimeUpdated   string `json:"timeUpdated"`
	Name          string `json:"name"`
	Owner         string `json:"owner"`
	Description   string `json:"description"`
	URL           string `json:"url"`
	Languages     string `json:"language"`
	Meta          string `json:"meta"`
	ForkCount     int    `json:"forkCounts"`
	StarCount     int    `json:"starCounts"`
	OpenIssues    int    `json:"openIssues"`
	WatchersCount int    `json:"watchersCount"`
}
