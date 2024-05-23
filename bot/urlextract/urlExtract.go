package urlextract

import (
	"log"
	"net/url"
	"strings"
)

const TIKTOK_HOST = "tiktok.com"

type UrlType int

const (
	VIDEO_TIKTOK UrlType = iota
	TEXT
	NONE
)

func (u UrlType) String() string {
	switch u {
	case VIDEO_TIKTOK:
		return "TikTok"
	case TEXT:
		return "text based"
	default:
		return "not a URL"
	}
}

type WordResult struct {
	UrlType    UrlType
	MatchedUrl *url.URL
}

func ExtractUrlsFromText(text string) (chan WordResult, int) {
	var fields = strings.Fields(text)
	results := make(chan WordResult, len(fields))

	first := fields[:len(fields)/2]
	second := fields[len(fields)/2:]

	if len(first) > 0 {
		go processSlice(first, results)
	}
	go processSlice(second, results)

	return results, len(fields)
}

func processSlice(words []string, results chan WordResult) {
	for _, field := range words {
		processUrl(field, results)
	}
}

func processUrl(field string, result chan WordResult) {
	parsedUrl, err := url.ParseRequestURI(field)

	switch {
	case err != nil:
		log.Println("Word not contains URL ", field)
		result <- WordResult{UrlType: NONE}
	case strings.Contains(parsedUrl.Hostname(), TIKTOK_HOST):
		log.Println("TikTok URL identified ", parsedUrl)
		result <- WordResult{UrlType: VIDEO_TIKTOK, MatchedUrl: parsedUrl}
	default:
		log.Println("Valid (probably text) URL identified ", parsedUrl)
		result <- WordResult{UrlType: TEXT, MatchedUrl: parsedUrl}
	}
}
