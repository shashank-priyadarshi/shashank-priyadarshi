package core

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"
)

type GitHubData struct {
	Data struct {
		Viewer struct {
			Repositories struct {
				Nodes []struct {
					Name     string    `json:"name"`
					URL      string    `json:"url"`
					PushedAt time.Time `json:"pushedAt"`
					Refs     struct {
						Nodes []struct {
							Name   string `json:"name"`
							Target struct {
								History struct {
									TotalCount int `json:"totalCount"`
								} `json:"history"`
							} `json:"target"`
						} `json:"nodes"`
					} `json:"refs"`
					PullRequests struct {
						TotalCount int `json:"totalCount"`
					} `json:"pullRequests"`
					Issues struct {
						TotalCount int `json:"totalCount"`
					} `json:"issues"`
				} `json:"nodes"`
			} `json:"repositories"`
		} `json:"viewer"`
	} `json:"data"`
}

type githubData struct {
	commits, prs, issues int
	projects             []Item
	list                 string
}

var (
	logger, _ = zap.NewProduction()
)

func (g *githubData) fetchGitHubData() {
	query := fmt.Sprintf(GithubQuery, RepoCount, getDateRange())
	body, err := httpClient("POST", GithubURL, fmt.Sprintf("Bearer %s", GithubToken), []byte(query), nil)
	if err != nil {
		logger.Sugar().Errorf("error making http request to %s: %s\n", GithubURL, err.Error())
		return
	}

	githubData := GitHubData{}
	err = json.Unmarshal(body, &githubData)
	if err != nil {
		logger.Sugar().Errorf("error unmarshaling github data: %s\n", err.Error())
		return
	}

	githubDataLength := len(githubData.Data.Viewer.Repositories.Nodes)

	for i := 0; i < githubDataLength; i++ {
		githubDataNode := githubData.Data.Viewer.Repositories.Nodes[i]

		if strings.EqualFold(githubDataNode.Name, "shashank-priyadarshi") {
			continue
		}

		g.prs, g.issues, g.projects = g.prs+githubDataNode.Issues.TotalCount, g.issues+githubDataNode.Issues.TotalCount, append(g.projects, Item{
			Title:     githubDataNode.Name,
			Permalink: githubDataNode.URL,
			Updated:   githubDataNode.PushedAt.String()[:10],
		})

		for _, val := range githubDataNode.Refs.Nodes {
			if strings.Contains(val.Name, "bot") {
				continue
			}
			g.commits += val.Target.History.TotalCount
		}

		g.list += fmt.Sprintf(ListItem, githubDataNode.URL, githubDataNode.Name, UpdatedAt, githubDataNode.PushedAt.String()[:10])
	}
}

func getDateRange() (date string) {
	date = time.Now().AddDate(0, -3, 0).Format(Layout)
	return
}
