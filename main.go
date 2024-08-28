package main

import (
	"os"

	"github.com/shashank-priyadarshi/shashank-priyadarshi/core"
	"go.uber.org/zap"
)

var (
	logger, _ = zap.NewProduction()
)

const (
	GithubOpenSourceStatsURL = "https://github-readme-stats.vercel.app/api?username=shashank-priyadarshi&show_icons=true&hide_border=true&include_all_commits=true&card_width=600&custom_title=GitHub%20Open%20Source%20Stats&title_color=3B7EBF&text_color=474A4E&icon_color=3B7EBF&hide=contribs&show=prs_merged&theme=transparent#gh-light-mode-only"
	GithubStreaksURL         = "https://streak-stats.demolab.com/?user=shashank-priyadarshi"
)

func main() {
	logger.Info("Fetching GitHub stats image")
	statsMap := map[string]string{
		"open_source_stats": GithubOpenSourceStatsURL,
		"streak_stats":      GithubStreaksURL,
	}

	if err := core.FetchStatsSVG(statsMap); err != nil {
		logger.Sugar().Errorf("Error fetching GitHub stats images")
	}

	logger.Info("Successfully fetched GitHub stats images")

	logger.Info("Starting script to auto update README")
	markdown := core.Markdown{}
	markdown.GenerateMarkdown()

	logger.Info("Successfully fetched and updated README")
	logger.Info("Writing README udpates to file")
	if err := os.WriteFile("README.md", []byte(markdown.Body), 0644); err != nil {
		logger.Sugar().Errorf("error writing markdown buffer file: %s\n", err.Error())
		return
	}

	logger.Info("README update successful")
}
