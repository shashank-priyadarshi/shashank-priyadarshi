package core

import "os"

const (
	RssURL                   = "https://hashnode.ssnk.in/rss.xml"
	GithubOpenSourceStatsURL = "https://github-readme-stats.vercel.app/api?username=shashank-priyadarshi&show_icons=true&hide_border=true&include_all_commits=true&card_width=600&custom_title=GitHub%20Open%20Source%20Stats&title_color=3B7EBF&text_color=474A4E&icon_color=3B7EBF&hide=contribs&show=prs_merged&theme=transparent#gh-light-mode-only"
	GithubStreaksURL         = "https://streak-stats.demolab.com/?user=shashank-priyadarshi"
	GithubURL                = "https://api.github.com/graphql"
	Layout                   = "2006-01-02T15:04:05Z"
	UpdatedAt                = "Last Updated "
	PublishedAt              = "Published "
)

var (
	RepoCount   = 6
	GithubQuery = `{
		"query": "query {viewer {repositories(first: %d ownerAffiliations: [OWNER] orderBy: {field: PUSHED_AT, direction: DESC}isArchived: false privacy: PUBLIC) {nodes {name url pushedAt refs(refPrefix: \"refs/heads/\", first: 5) {nodes {name target {... on Commit {history(first: 100, since: \"%s\", author: {emails: [\"shashank9163882019@gmail.com\"]}) {totalCount}}}}} pullRequests(states: MERGED first: 100 orderBy: {field: UPDATED_AT, direction: DESC}) {totalCount} issues(states: CLOSED first: 100 orderBy: {field: UPDATED_AT, direction: DESC}) {totalCount}}}}}"
	}`
	GithubToken = os.Getenv("GITHUB_TOKEN")
	ListItem    = `<li><a href="%s" target="_blank" rel="noopener noreferrer">%s</a> %s: %s</li>`
	Header      = `<div align="center"><p><a href="https://ssnk.in"><img src="https://img.shields.io/badge/-Website-3B7EBF?style=for-the-badge&amp;logo=amp&amp;logoColor=white" alt="Website Badge"></a> <a href="https://hashnode.ssnk.in"><img src="https://img.shields.io/badge/-Blog-3B7EBF?style=for-the-badge&amp;logo=Hashnode&amp;logoColor=white" alt="Blog Badge"></a> <a href="https://linkedin.com/in/shashank-priyadarshi"><img src="https://img.shields.io/badge/-LinkedIn-3B7EBF?style=for-the-badge&amp;logo=Linkedin&amp;logoColor=white" alt="Linkedin Badge"></a> <a href="https://peerlist.io/shasha"><img src="https://img.shields.io/badge/-PeerList-3B7EBF?style=for-the-badge&amp;logo=Peerlist&amp;logoColor=white" alt="Peerlist Badge"/></a></p><hr><p>Hi there 👋! I'm <a href="https://ssnk.in">Shashank</a>. I am a Backend Developer, currently building distributed payment solutions at <a href="https://npci.org.in">NPCI</a>. I like tinkering, and writing code, some of which I have pinned below. Sometimes I play <a href="https://www.chess.com/member/ttefabob">chess</a>, and then I procrastinate.</p><hr><p><img src="./assets/images/streak_stats.svg"/></p><hr><p><img src="./assets/images/open_source_stats.svg"/></p><hr><h2>Highlights</h2><details><summary>Projects</summary><br /><ul>%s</ul></details><details><summary>Recent Blogposts</summary><br /><ul>%s</ul></details><hr></div></br>`
)

type Item struct {
	Title, Permalink, Updated string
}
