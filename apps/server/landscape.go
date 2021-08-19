package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	"github.com/tidwall/gjson"
	"gopkg.in/yaml.v2"
	"local.packages/common"
)

const (
	LandscapeUrl = "https://raw.githubusercontent.com/cncf/landscape/master/landscape.yml"
)

var client http.Client

// パフォーマンス改善、GitHub API のレートリミット対策としてキャッシュ
var stars map[string]int64

type LandScape struct {
	Landscape []struct {
		Category      string
		Name          string
		Subcategories []struct {
			Subcategory string
			Name        string
			Items       []SubItem
		}
	}
}

type SubItem struct {
	Extra struct {
		Accepted         string
		DevStatsUrl      string `yaml:"dev_stats_url"`
		ArtworkUrl       string `yaml:"artwork_url"`
		StackOverflowUrl string `yaml:"stack_overflow_url"`
		BlogUrl          string `yaml:"blog_url"`
		SlackUrl         string `yaml:"slack_url"`
		YoutubeUrl       string `yaml:"youtube_url"`
	}
	Name        string
	Description string
	HomepageUrl string `yaml:"homepage_url"`
	Project     string
	RepoUrl     string `yaml:"repo_url"`
	Logo        string
	Twitter     string
	Crunchbase  string
}

type Project common.Project

func getCicdProjects() ([]Project, error) {
	resp, err := client.Get(LandscapeUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return findCicdProjects(b)
	}
	return nil, nil
}

func findCicdProjects(data []byte) ([]Project, error) {
	l := LandScape{}

	var err = yaml.Unmarshal([]byte(data), &l)
	if err != nil {
		return nil, err
	}

	ml := getMemberList(l)
	for _, category := range l.Landscape {
		if category.Name == "App Definition and Development" {
			for _, sub := range category.Subcategories {
				if sub.Name == "Continuous Integration & Delivery" {
					var wg sync.WaitGroup
					list := make([]Project, len(sub.Items))
					for i, proj := range sub.Items {
						wg.Add(1)
						go func(i int, proj SubItem) {
							defer wg.Done()
							list[i] = Project{
								Name:        proj.Name,
								Description: proj.Description,
								HomepageUrl: proj.HomepageUrl,
								Project:     getProject(proj.Project, proj.Crunchbase, ml),
								RepoUrl:     proj.RepoUrl,
								Crunchbase:  proj.Twitter,
								StarCount:   getStarCount(proj.RepoUrl),
							}
						}(i, proj)
					}
					wg.Wait()
					return list, nil
				}
			}
		}
	}
	return nil, nil
}

func getMemberList(l LandScape) map[string]bool {
	ml := make(map[string]bool)
	for _, category := range l.Landscape {
		if category.Name == "CNCF Members" {
			for _, sub := range category.Subcategories {
				for _, proj := range sub.Items {
					ml[proj.Crunchbase] = true
				}
			}
		}
	}
	return ml
}

func getProject(project, crunchbase string, ml map[string]bool) string {
	if project != "" {
		return project
	}
	if _, ok := ml[crunchbase]; ok {
		return "member"
	}
	return ""
}

func getStarCount(repoUrl string) int64 {
	if repoUrl == "" {
		return 0
	}

	if v, ok := stars[repoUrl]; ok {
		return v
	}

	apiUrl := strings.Replace(repoUrl, "github.com", "api.github.com/repos", 1)
	fmt.Println("Access to ", apiUrl)
	resp, err := client.Get(apiUrl)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return 0
	}

	if resp.StatusCode == http.StatusOK {
		c := gjson.GetBytes(b, "stargazers_count").Int()
		stars[repoUrl] = c
		return c
	} else {
		fmt.Println(string(b))
		return 0
	}
}
