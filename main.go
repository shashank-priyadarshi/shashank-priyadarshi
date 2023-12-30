package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/mmcdole/gofeed"
)

const (
	RSS_URL      = "https://blog.ssnk.in/rss.xml"
	GITHUB_URL   = "https://api.github.com/graphql"
	GITHUB_QUERY = `{
		"query": "query { viewer { repositories(first: 100, isFork: false, ownerAffiliations: [OWNER]) { nodes { name url isFork defaultBranchRef { target { repository { updatedAt } ... on Commit { history(first: 100, since: \"2023-01-01T00:00:00Z\") { totalCount } } } } pullRequests(states: MERGED, first: 100, orderBy: {field: UPDATED_AT, direction: DESC}) { totalCount } issues(states: CLOSED, first: 100, orderBy: {field: UPDATED_AT, direction: DESC}) { totalCount } } } } }"
	}`
	LAYOUT       = "2006-01-02T15:04:05Z07:00"
	UPDATED_AT   = "Last Updated "
	PUBLISHED_AT = "Published "
)

var (
	GITHUB_TOKEN = os.Getenv("GITHUB_TOKEN")
	LIST_ITEM    = `<li><a href="%s" target="_blank" rel="noopener noreferrer">%s</a> %s: %s</li>`
	HEADER       = `<div align="center"><p><a href="https://ssnk.in"><img src="https://img.shields.io/badge/-Website-3B7EBF?style=for-the-badge&amp;logo=amp&amp;logoColor=white" alt="Website Badge"></a> <a href="https://blog.ssnk.in"><img src="https://img.shields.io/badge/-Blog-3B7EBF?style=for-the-badge&amp;logo=Hashnode&amp;logoColor=white" alt="Blog Badge"></a> <a href="https://linkedin.com/in/shashank-priyadarshi"><img src="https://img.shields.io/badge/-LinkedIn-3B7EBF?style=for-the-badge&amp;logo=Linkedin&amp;logoColor=white" alt="Linkedin Badge"> <a href="https://peerlist.io/shasha"></a><img src="https://img.shields.io/badge/-PeerList-3B7EBF?style=for-the-badge&amp;logo=Peerlist&amp;logoColor=white" alt="Peerlist Badge"/></p><hr><p>Hi there 👋! I'm <a href="https://ssnk.in">Shashank</a>. I am a Backend Developer, currently building distributed payment solutions at <a href="https://npci.org.in">NPCI</a>. I like tinkering, and writing code, some of which I have pinned below. Sometimes I play <a href="https://www.chess.com/member/ttefabob">chess</a>, and then I procrastinate.</p><hr><p><a href="https://github.com/shashank-priyadarshi/shashank-priyadarshi#gh-light-mode-only"><img src="https://github-readme-stats.vercel.app/api?username=shashank-priyadarshi&amp;show_icons=true&amp;hide_border=true&amp;include_all_commits=true&amp;card_width=600&amp;title_color=3B7EBF&amp;text_color=474A4E&amp;icon_color=3B7EBF&amp;hide=contribs&amp;show=prs_merged&amp;theme=transparent#gh-light-mode-only" alt="GitHub-Stats-Card-Light"></a></p><hr><h2>Highlights</h2><details><summary>Projects</summary><br /><ul>%s</ul></details><details><summary>Recent Blogposts</summary><br /><ul>%s</ul></details><hr></div>`
)

func main() {
	markdown := markdown{}
	markdown.generateMarkdown()

	if err := os.WriteFile("README.md", []byte(markdown.body), 0644); err != nil {
		fmt.Printf("error writing markdown buffer file: %s\n", err.Error())
		return
	}
}

type markdown struct {
	title string // website, blog, linkedin, profile views
	body  string // intro, active repo list(with link), OSS stats(commits, PRs, PRs merged, issues), latest articles, devcard
}

func (m *markdown) generateMarkdown() {
	githubData := githubData{}
	githubData.fetchGitHubData()
	rss := rss{}
	rss.fetchRSSData()

	m.body = fmt.Sprintf(HEADER, githubData.list, rss.list)
}

type githubData struct {
	commits, prs, issues int
	projects             []item
	list                 string
}

func (g *githubData) fetchGitHubData() {
	body, err := httpClient("POST", GITHUB_URL, GITHUB_QUERY, GITHUB_TOKEN)
	if err != nil {
		fmt.Printf("error making http request to %s: %s\n", GITHUB_URL, err.Error())
		return
	}

	githubData := GitHubData{}
	err = json.Unmarshal(body, &githubData)
	if err != nil {
		fmt.Printf("error unmarshaling github data: %s\n", err.Error())
		return
	}

	githubDataLength := len(githubData.Data.Viewer.Repositories.Nodes)

	for i := 0; i < githubDataLength-1; i++ {
		githubDataNode := githubData.Data.Viewer.Repositories.Nodes[i]
		if githubDataNode.IsFork {
			fmt.Printf("%s is a forked repository", githubDataNode.Name)
			continue
		}

		// updatedAt, err := time.Parse(LAYOUT, githubDataNode.DefaultBranchRef.Target.Repository.UpdatedAt)
		// if err != nil {
		// 	fmt.Printf("error unmarshaling updatedAt %s for repository %s: %s\n", githubDataNode.DefaultBranchRef.Target.Repository.UpdatedAt, githubDataNode.Name, err.Error())
		// }

		g.commits, g.prs, g.issues, g.projects = g.commits+githubDataNode.DefaultBranchRef.Target.History.TotalCount, g.prs+githubDataNode.Issues.TotalCount, g.issues+githubDataNode.Issues.TotalCount, append(g.projects, item{
			title:     githubDataNode.Name,
			permalink: githubDataNode.URL,
			updated:   githubDataNode.DefaultBranchRef.Target.Repository.UpdatedAt[:10],
		})

		g.list += fmt.Sprintf(LIST_ITEM, githubDataNode.URL, githubDataNode.Name, UPDATED_AT, githubDataNode.DefaultBranchRef.Target.Repository.UpdatedAt[:10])
	}
}

type rss struct {
	articles []item
	list     string
}

func (r *rss) fetchRSSData() {
	xmlParser := gofeed.NewParser()
	feed, err := xmlParser.ParseURL(RSS_URL)
	if err != nil {
		fmt.Printf("error parsing rss: %s\n", err.Error())
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
					Name, URL        string
					IsFork           bool
					DefaultBranchRef struct {
						Target struct {
							Repository struct {
								UpdatedAt string
							}
							History struct {
								TotalCount int
							}
						}
					}
					PullRequests struct {
						TotalCount int
					}
					Issues struct {
						TotalCount int
					}
				}
			}
		}
	}
}

func httpClient(method, url, body, auth string) (resBody []byte, err error) {

	client := http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewReader([]byte(body)))
	if err != nil {
		return
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", auth))

	response, err := client.Do(req)
	if err != nil || response.StatusCode != 200 {
		return
	}
	defer response.Body.Close()

	resBody, err = io.ReadAll(response.Body)
	if err != nil {
		return
	}

	return
}
