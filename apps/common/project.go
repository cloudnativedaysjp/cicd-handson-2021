package common

type Project struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	HomepageUrl string `json:"homepage_url"`
	Project     string `json:"project"`
	RepoUrl     string `json:"repo_url"`
	Crunchbase  string `json:"crunchbase"`
	StarCount int64 `json:"star_count"`
}

type Projects []Project
