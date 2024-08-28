package core

import (
	"fmt"
)

type Markdown struct {
	// title string // website, blog, linkedin, profile views
	Body string // intro, active repo list(with link), OSS stats(commits, PRs, PRs merged, issues), latest articles, devcard
}

func (m *Markdown) GenerateMarkdown() {
	logger.Info("Starting README generation")
	githubData := githubData{}

	logger.Info("Fetching GitHub data")
	githubData.fetchGitHubData()

	logger.Info("Fetching RSS data")
	rss := rss{}
	rss.fetchRSSData()

	logger.Info("Setting README body")
	m.Body = fmt.Sprintf(Header, githubData.list, rss.list)
}
