package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
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
	urlType    UrlType
	matchedUrl *url.URL
}

func OnNewMessage(discord *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == discord.State.User.ID {
		return
	}

	var fields = strings.Fields(message.Content)
	results := make(chan WordResult, len(fields))

	for _, field := range fields {
		go processUrl(field, results)
	}

	go respondStatus(discord, message, results, len(fields))
}

func processUrl(field string, result chan WordResult) {
	parsedUrl, err := url.ParseRequestURI(field)

	switch {
	case err != nil:
		log.Println("Word not contains URL ", field)
		result <- WordResult{urlType: NONE}
	case strings.Contains(parsedUrl.Hostname(), TIKTOK_HOST):
		log.Println("TikTok URL identified ", parsedUrl)
		result <- WordResult{urlType: VIDEO_TIKTOK, matchedUrl: parsedUrl}
	default:
		log.Println("Valid (probably text) URL identified ", parsedUrl)
		result <- WordResult{urlType: TEXT, matchedUrl: parsedUrl}
	}
}

func respondStatus(discord *discordgo.Session, message *discordgo.MessageCreate, results chan WordResult, count int) {
	var response string
	for range count {
		result := <-results
		if result.urlType != NONE {
			response += fmt.Sprintf("- %s: %s\n", result.urlType, result.matchedUrl.Hostname())
		}
	}
	if response == "" {
		return
	}
	response = "Found URL(s):\n" + response
	_, err := discord.ChannelMessageSend(message.ChannelID, response)
	if err != nil {
		log.Println("Error when sending response: ", err)
	}
}
