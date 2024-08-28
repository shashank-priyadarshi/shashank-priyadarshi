package core

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

//c := &Chess{}
//
//c.getStats("").
//	getRatingLine("chess_blitz", "⚡", "Blitz", 52).
//	getRatingLine("chess_bullet", "🚅", "Bullet", 52).
//	getRatingLine("chess_rapid", "⏲️", "Rapid", 53).
//	getRatingLine("tactics", "🧩", "Tactics", 52).
//	getRatingLine("chess_daily", "☀️", "Daily", 53)
//
//data := Gist{
//	Files: Files{
//		Content: Content{
//			Content: strings.ReplaceAll(c.data.String(), "separator", "\n"),
//		},
//	},
//}
//
//body, err := json.Marshal(data)
//if err != nil {
//	panic(err)
//}
//
//fmt.Println(string(body))
//
//err = updateGist(Title, body)
//fmt.Println(err)

func updateGist(title string, content []byte) error {
	githubToken := os.Getenv("GITHUB_TOKEN")
	gistID := os.Getenv("GIST_ID")

	_, err := httpClient("PATCH", fmt.Sprintf("https://api.github.com/gists/%s", gistID), fmt.Sprintf("token %s", githubToken), nil, nil)
	if err != nil {
		return err
	}

	return nil
}

const (
	Title = "♟︎ Chess.com Ratings"
)

var (
	LiveURL    = "https://www.chess.com/stats/live/{format}/{user}"
	PuzzlesURL = "https://www.chess.com/stats/{format}/{user}"
	DailyURL   = "https://www.chess.com/stats/{format}/chess/{user}"
	StatsURL   = "https://api.chess.com/pub/player/%s/stats"
)

type Chess struct {
	data  strings.Builder
	stats *Stats
}

func (c *Chess) getRatingLine(key, emoji, format string, maxLineLen int) *Chess {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf(`%s \%s`, emoji, format))
	title := builder.String()
	builder.Reset()

	var rating int
	switch format {
	case "Blitz":
		rating = c.stats.ChessBlitz.Last.Rating
	case "Bullet":
		rating = c.stats.ChessBullet.Last.Rating
	case "Rapid":
		rating = c.stats.ChessRapid.Last.Rating
	case "Daily":
		rating = c.stats.ChessDaily.Last.Rating
	case "Tactics":
		rating = c.stats.Tactics.Highest.Rating
	default:
		builder.WriteString("N/A")
	}

	builder.WriteString(fmt.Sprintf("%d 📈", rating))
	value := builder.String()
	builder.Reset()

	separation := maxLineLen - (len(title) + len(value))

	builder.WriteString(" ")
	for range separation {
		builder.WriteString(".")
	}
	builder.WriteString(" ")

	separator := builder.String()

	c.data.WriteString(fmt.Sprintf("%s%s%sseparator", title, separator, value))
	return c
}

func (c *Chess) getStats(user string) *Chess {

	headers := map[string]string{
		"User-Agent": fmt.Sprintf("chess-com-box-py @%s", user),
	}

	body, err := httpClient("GET", fmt.Sprintf(StatsURL, user), "", nil, headers)
	if err != nil {
		panic(err)
	}

	var s Stats
	err = json.Unmarshal(body, &s)
	if err != nil {
		panic(err)
	}

	c.stats = &s
	return c
}

type Gist struct {
	Files Files `json:"files"`
}

type Files struct {
	Content Content `json:"content"`
}

type Content struct {
	Content string `json:"content"`
}

type Stats struct {
	ChessDaily struct {
		Last struct {
			Rating int `json:"rating"`
			Date   int `json:"date"`
			Rd     int `json:"rd"`
		} `json:"last"`
		Record struct {
			Win            int `json:"win"`
			Loss           int `json:"loss"`
			Draw           int `json:"draw"`
			TimePerMove    int `json:"time_per_move"`
			TimeoutPercent int `json:"timeout_percent"`
		} `json:"record"`
	} `json:"chess_daily"`
	ChessRapid struct {
		Last struct {
			Rating int `json:"rating"`
			Date   int `json:"date"`
			Rd     int `json:"rd"`
		} `json:"last"`
		Best struct {
			Rating int    `json:"rating"`
			Date   int    `json:"date"`
			Game   string `json:"game"`
		} `json:"best"`
		Record struct {
			Win  int `json:"win"`
			Loss int `json:"loss"`
			Draw int `json:"draw"`
		} `json:"record"`
	} `json:"chess_rapid"`
	ChessBullet struct {
		Last struct {
			Rating int `json:"rating"`
			Date   int `json:"date"`
			Rd     int `json:"rd"`
		} `json:"last"`
		Best struct {
			Rating int    `json:"rating"`
			Date   int    `json:"date"`
			Game   string `json:"game"`
		} `json:"best"`
		Record struct {
			Win  int `json:"win"`
			Loss int `json:"loss"`
			Draw int `json:"draw"`
		} `json:"record"`
	} `json:"chess_bullet"`
	ChessBlitz struct {
		Last struct {
			Rating int `json:"rating"`
			Date   int `json:"date"`
			Rd     int `json:"rd"`
		} `json:"last"`
		Best struct {
			Rating int    `json:"rating"`
			Date   int    `json:"date"`
			Game   string `json:"game"`
		} `json:"best"`
		Record struct {
			Win  int `json:"win"`
			Loss int `json:"loss"`
			Draw int `json:"draw"`
		} `json:"record"`
	} `json:"chess_blitz"`
	Fide    int `json:"fide"`
	Tactics struct {
		Highest struct {
			Rating int `json:"rating"`
			Date   int `json:"date"`
		} `json:"highest"`
		Lowest struct {
			Rating int `json:"rating"`
			Date   int `json:"date"`
		} `json:"lowest"`
	} `json:"tactics"`
	PuzzleRush struct {
		Best struct {
			TotalAttempts int `json:"total_attempts"`
			Score         int `json:"score"`
		} `json:"best"`
	} `json:"puzzle_rush"`
}
