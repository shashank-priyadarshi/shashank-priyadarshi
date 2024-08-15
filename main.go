package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
	"go.uber.org/zap"
)

const (
	RSS_URL          = "https://blog.ssnk.in/rss.xml"
	GITHUB_STATS_URL = "https://github-readme-stats.vercel.app/api?username=shashank-priyadarshi&show_icons=true&hide_border=true&include_all_commits=true&card_width=600&custom_title=GitHub%20Open%20Source%20Stats&title_color=3B7EBF&text_color=474A4E&icon_color=3B7EBF&hide=contribs&show=prs_merged&theme=transparent#gh-light-mode-only"
	GITHUB_URL       = "https://api.github.com/graphql"
	LAYOUT           = "2006-01-02T15:04:05Z"
	UPDATED_AT       = "Last Updated "
	PUBLISHED_AT     = "Published "
)

var (
	logger, _    = zap.NewProduction()
	REPOCOUNT    = 6
	GITHUB_QUERY = `{
		"query": "query {viewer {repositories(first: %d ownerAffiliations: [OWNER] orderBy: {field: PUSHED_AT, direction: DESC}isArchived: false privacy: PUBLIC) {nodes {name url pushedAt refs(refPrefix: \"refs/heads/\", first: 5) {nodes {name target {... on Commit {history(first: 100, since: \"%s\", author: {emails: [\"shashank9163882019@gmail.com\"]}) {totalCount}}}}} pullRequests(states: MERGED first: 100 orderBy: {field: UPDATED_AT, direction: DESC}) {totalCount} issues(states: CLOSED first: 100 orderBy: {field: UPDATED_AT, direction: DESC}) {totalCount}}}}}"
	}`
	GITHUB_TOKEN = os.Getenv("GITHUB_TOKEN")
	LIST_ITEM    = `<li><a href="%s" target="_blank" rel="noopener noreferrer">%s</a> %s: %s</li>`
	HEADER       = `<div align="center"><p><a href="https://ssnk.in"><img src="https://img.shields.io/badge/-Website-3B7EBF?style=for-the-badge&amp;logo=amp&amp;logoColor=white" alt="Website Badge"></a> <a href="https://blog.ssnk.in"><img src="https://img.shields.io/badge/-Blog-3B7EBF?style=for-the-badge&amp;logo=Hashnode&amp;logoColor=white" alt="Blog Badge"></a> <a href="https://linkedin.com/in/shashank-priyadarshi"><img src="https://img.shields.io/badge/-LinkedIn-3B7EBF?style=for-the-badge&amp;logo=Linkedin&amp;logoColor=white" alt="Linkedin Badge"></a> <a href="https://peerlist.io/shasha"><img src="https://img.shields.io/badge/-PeerList-3B7EBF?style=for-the-badge&amp;logo=Peerlist&amp;logoColor=white" alt="Peerlist Badge"/></a></p><hr><p>Hi there 👋! I'm <a href="https://ssnk.in">Shashank</a>. I am a Backend Developer, currently building distributed payment solutions at <a href="https://npci.org.in">NPCI</a>. I like tinkering, and writing code, some of which I have pinned below. Sometimes I play <a href="https://www.chess.com/member/ttefabob">chess</a>, and then I procrastinate.</p><hr><p><img src="./assets/images/stats.svg"/></p><hr><h2>Highlights</h2><details><summary>Projects</summary><br /><ul>%s</ul></details><details><summary>Recent Blogposts</summary><br /><ul>%s</ul></details><hr></div>`
)

func main() {
	logger.Info("Fetching GitHub stats image")
	if err := fetchGitHubStatsSVG(); err != nil {
		logger.Sugar().Errorf("Error fetching GitHub stats image")
	}

	logger.Info("Successfully fetched GitHub stats image")

	logger.Info("Starting script to auto update README")
	markdown := markdown{}
	markdown.generateMarkdown()

	logger.Info("Successfully fetched and updated README")
	logger.Info("Writing README udpates to file")
	if err := os.WriteFile("README.md", []byte(markdown.body), 0644); err != nil {
		logger.Sugar().Errorf("error writing markdown buffer file: %s\n", err.Error())
		return
	}

	logger.Info("README update successful")
}

func fetchGitHubStatsSVG() (err error) {
	body, err := httpClient("GET", GITHUB_STATS_URL, "", "", "")
	if err != nil {
		logger.Sugar().Errorf("error making http request to %s: %s\n", GITHUB_STATS_URL, err.Error())
		return
	}
	if err = os.WriteFile("./assets/images/stats.svg", body, 0644); err != nil {
		logger.Sugar().Errorf("error writing stats svg file: %s\n", err.Error())
		return
	}
	return
}

type markdown struct {
	// title string // website, blog, linkedin, profile views
	body string // intro, active repo list(with link), OSS stats(commits, PRs, PRs merged, issues), latest articles, devcard
}

func (m *markdown) generateMarkdown() {
	logger.Info("Starting README generation")
	githubData := githubData{}

	logger.Info("Fetching GitHub data")
	githubData.fetchGitHubData()
	
	logger.Info("Fetching RSS data")
	rss := rss{}
	rss.fetchRSSData()

	logger.Info("Setting README body")
	m.body = fmt.Sprintf(HEADER, githubData.list, rss.list)
}

type githubData struct {
	commits, prs, issues int
	projects             []item
	list                 string
}

func (g *githubData) fetchGitHubData() {
	query := fmt.Sprintf(GITHUB_QUERY, REPOCOUNT, getDateRange())
	body, err := httpClient("POST", GITHUB_URL, query, "Authorization", fmt.Sprintf("Bearer %s", GITHUB_TOKEN))
	if err != nil {
		logger.Sugar().Errorf("error making http request to %s: %s\n", GITHUB_URL, err.Error())
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

		g.prs, g.issues, g.projects = g.prs+githubDataNode.Issues.TotalCount, g.issues+githubDataNode.Issues.TotalCount, append(g.projects, item{
			title:     githubDataNode.Name,
			permalink: githubDataNode.URL,
			updated:   githubDataNode.PushedAt.String()[:10],
		})

		for _, val := range githubDataNode.Refs.Nodes {
			if strings.Contains(val.Name, "bot") {
				continue
			}
			g.commits += val.Target.History.TotalCount
		}

		g.list += fmt.Sprintf(LIST_ITEM, githubDataNode.URL, githubDataNode.Name, UPDATED_AT, githubDataNode.PushedAt.String()[:10])
	}
}

func getDateRange() (date string) {
	date = time.Now().AddDate(0, -3, 0).Format(LAYOUT)
	return
}

type rss struct {
	articles []item
	list     string
}

func (r *rss) fetchRSSData() {
	xmlParser := gofeed.NewParser()
	feed, err := xmlParser.ParseURL(RSS_URL)
	if err != nil {
		logger.Sugar().Errorf("error parsing rss: %s\n", err.Error())
		return
	}
	feedLength := len(feed.Items)
	for i := 0; i < feedLength-1; i++ {
		feedItem := feed.Items[i]
		r.articles = append(r.articles, item{
			title:     feedItem.Title,
			permalink: feedItem.Link,
			updated:   feedItem.Updated,
		})
		r.list += fmt.Sprintf(LIST_ITEM, feedItem.Link, feedItem.Title, PUBLISHED_AT, feedItem.PublishedParsed.String()[:10])
	}
}

type item struct {
	title, permalink, updated string
}

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

func httpClient(method, url, body, authMechanism, auth string) (resBody []byte, err error) {
	client := http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewReader([]byte(body)))
	if err != nil {
		return
	}

	switch authMechanism {
	case "Authorization":
		req.Header.Add(authMechanism, auth)
	default:
	}

	response, err := client.Do(req)
	if err != nil {
		logger.Sugar().Info(err)
		return
	}

	if response.StatusCode != 200 {
		err = fmt.Errorf("status code: %v", response.StatusCode)
		return
	}

	defer response.Body.Close()

	resBody, err = io.ReadAll(response.Body)
	if err != nil {
		return
	}

	return
}
