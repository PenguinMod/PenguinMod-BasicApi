package main

import (
	"encoding/json"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2/log"
)

var recentCommits []GHApiCommit

type GHAuthor struct {
	Login     string `json:"login"`
	AvatarUrl string `json:"avatar_url"`
	HtmlUrl   string `json:"html_url"`
}

type GHCommit struct {
	Author  GHCommitAuthor `json:"author"`
	Message string         `json:"message"`
}

type GHCommitAuthor struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Date  string `json:"date"`
}

type GHApiCommit struct {
	Author  GHAuthor `json:"author"`
	Commit  GHCommit `json:"commit"`
	HtmlUrl string   `json:"html_url"`
}

type GHCommitsByDate []GHApiCommit

func (d GHCommitsByDate) Len() int      { return len(d) }
func (d GHCommitsByDate) Swap(i, j int) { d[i], d[j] = d[j], d[i] }
func (d GHCommitsByDate) Less(i, j int) bool {
	time1, _ := time.Parse(time.RFC3339, d[i].Commit.Author.Date)
	time2, _ := time.Parse(time.RFC3339, d[j].Commit.Author.Date)
	return time1.Before(time2)
}

func getRecentsCommits() {
	githubCommitApis := []string{
		"https://api.github.com/repos/PenguinMod/penguinmod.github.io/commits?per_page=50",
		"https://api.github.com/repos/PenguinMod/PenguinMod-Vm/commits?per_page=50",
		"https://api.github.com/repos/PenguinMod/PenguinMod-Home/commits?per_page=50",
		"https://api.github.com/repos/PenguinMod/PenguinMod-Blocks/commits?per_page=50",
		"https://api.github.com/repos/PenguinMod/PenguinMod-Paint/commits?per_page=50",
		"https://api.github.com/repos/PenguinMod/PenguinMod-Packager/commits?per_page=50",
		"https://api.github.com/repos/PenguinMod/PenguinMod-Render/commits?per_page=50",
		"https://api.github.com/repos/PenguinMod/PenguinMod-ExtensionsGallery/commits?per_page=50",
	}

	var newRecentCommits []GHApiCommit
	for i := 0; i < len(githubCommitApis); i++ {
		resp, err := http.Get(githubCommitApis[i])
		if err != nil {
			log.Errorf("Failed fetching %s: %s", githubCommitApis[i], err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			log.Errorf("Failed fetching %s: Non-OK status code; %s", githubCommitApis[i], strconv.Itoa(resp.StatusCode))
			continue
		}

		var apiResp []GHApiCommit
		err = json.NewDecoder(resp.Body).Decode(&apiResp)
		if err != nil {
			log.Errorf("Failed decoding response from %s: %s", githubCommitApis[i], err)
			continue
		}

		newRecentCommits = append(newRecentCommits, apiResp...)
	}

	sort.Sort(GHCommitsByDate(newRecentCommits))
	if len(newRecentCommits) >= 200 {
		recentCommits = newRecentCommits[:200]
	} else {
		recentCommits = newRecentCommits
	}
}
