package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	"v2.staffjoy.com/errorpages"
)

type breaktimeEpisode struct {
	// Generic page stuff
	Title        string // Used in <title>
	Description  string // SEO matters
	CSSId        string // e.g. 'careers'
	Version      string // e.g. master-1, for cachebusting
	TemplateName string
	CsrfField    template.HTML
	// Message stuff
	Name              string
	SoundcloudTrackID string
	Body              template.HTML
	CoverPhoto        string
	Date              string
}

// For new episodes add markdown to breaktime-content folder
// and cover photo to breaktime-cover
var episodes = map[string]breaktimeEpisode{
	"kaldis-coffee": {
		Title:             "Breaktime by Staffjoy Episode 1 - Kaldi’s Coffee",
		Name:              "Episode 1 - Tyler Zimmer, Owner of Kaldi's Coffee",
		Description:       "In the first episode of Breaktime by Staffjoy we interview Tyler Zimmer, owner of Kaldi’s Coffee, to learn the business strategies that helped them grow.",
		SoundcloudTrackID: "293713246",
		Date:              "November 29, 2016",
	},
	"sump-coffee": {
		Title:             "Breaktime by Staffjoy Episode 2 - Specialty Coffee at Sump",
		Name:              "How Sump Brought Specialty Coffee to St. Louis",
		Description:       "In the second episode of Breaktime by Staffjoy, we interview Scott Carey, owner of Sump Coffee, to learn how he brought specialty coffee to St. Louis.",
		SoundcloudTrackID: "297151821",
		Date:              "December 14, 2016",
	},
	"workshop-cafe": {
		Title:             "Breaktime by Staffjoy Episode 3 - Workshop Cafe’s Technology and Brick and Mortar Success",
		Name:              "How Workshop Cafe Combined Technology with Brick and Mortar",
		Description:       "In the third episode of Breaktime, we interview Rich Menendez, CEO of Workshop Cafe, to learn how he blends technology and brick and mortar to perfection.",
		SoundcloudTrackID: "297605514",
		Date:              "January 3, 2017",
	},
}

func breaktimeEpisodeHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/html; charset=UTF-8")
	vars := mux.Vars(req)
	slug := vars["slug"]

	episode, ok := episodes[slug]
	if !ok {
		res.WriteHeader(http.StatusNotFound)
		errorpages.NotFound(res)
		return
	}
	res.WriteHeader(http.StatusOK)
	body, ok := breaktimeSource[slug]
	if !ok {
		logger.Panicf("cannot find episode body for slug %s", slug)
	}
	episode.Body = template.HTML(body)
	episode.Version = config.GetDeployVersion()
	episode.CSSId = "breaktimeEpisode"
	episode.CoverPhoto = fmt.Sprintf("/assets/breaktime-cover/%s.jpg", slug)

	err := tmpl.ExecuteTemplate(res, "breaktime_episode.tmpl", episode)
	if err != nil {
		logger.Panicf("Unable to render page %s - %s", episode.Title, err)
	}
}

type breaktimeList struct {
	// Generic page stuff
	Title        string // Used in <title>
	Description  string // SEO matters
	CSSId        string // e.g. 'careers'
	Version      string // e.g. master-1, for cachebusting
	TemplateName string
	CsrfField    template.HTML
	// Message stuff
	Episodes map[string]breaktimeEpisode
}

func breaktimeListHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/html; charset=UTF-8")
	res.WriteHeader(http.StatusOK)
	p := &breaktimeList{
		Title: "Breaktime by Staffjoy",
	}
	p.Version = config.GetDeployVersion()
	p.CSSId = "breaktimeList"
	p.Episodes = episodes
	p.Description = "Breaktime is a blog and podcast series by Staffjoy that interviews top business leaders to learn their top tips for improving and growing their businesses."

	err := tmpl.ExecuteTemplate(res, "breaktime_list.tmpl", p)
	if err != nil {
		logger.Panicf("Unable to render page %s - %s", p.Title, err)
	}
}
