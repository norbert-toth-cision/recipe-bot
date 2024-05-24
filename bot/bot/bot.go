package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"recipebot/queue"
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
		response += reportResult(result)
		queue.PushIntoQueue(result)
	}
	sendReport(discord, message.ChannelID, response)
}

func reportResult(result urlextract.WordResult) string {
	if result.UrlType != urlextract.NONE {
		return fmt.Sprintf("- %s: %s\n", result.UrlType, result.MatchedUrl.Hostname())
	}
	return ""
}

func sendReport(discord *discordgo.Session, channelId string, report string) {
	if report == "" {
		return
	}
	report = "Found URL(s):\n" + report
	_, err := discord.ChannelMessageSend(channelId, report)
	if err != nil {
		log.Println("Error when sending response: ", err)
	}
}
