package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"recipebot/urlextract"
)

func OnNewMessage(discord *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == discord.State.User.ID {
		return
	}
	results, count := urlextract.ExtractUrlsFromText(message.Content)
	respondStatus(discord, message, results, count)
}

func respondStatus(discord *discordgo.Session, message *discordgo.MessageCreate, results chan urlextract.WordResult, count int) {
	var response string
	for range count {
		result := <-results
		if result.UrlType != urlextract.NONE {
			response += fmt.Sprintf("- %s: %s\n", result.UrlType, result.MatchedUrl.Hostname())
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
