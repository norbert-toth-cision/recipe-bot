package bot

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"recipebot/config"
	"recipebot/queue"
	"recipebot/urlextract"
)

const (
	Name = "RecipeBot"
)

type RecipeBot struct {
	config   config.Config
	session  *discordgo.Session
	msgQueue queue.Queue
}

func (rb *RecipeBot) Configure(config config.Config) error {
	rb.config = config
	rb.msgQueue = new(queue.RMQueue)
	rb.msgQueue.Configure(config)
	return nil
}

func (rb *RecipeBot) Start() error {
	log.Println("Starting", Name)
	var err error
	rb.session, err = discordgo.New("Bot " + rb.config.GetString(config.BotToken))
	if err != nil {
		return err
	}

	rb.session.AddHandler(rb.OnNewMessage)
	err = rb.session.Open()
	if err != nil {
		return err
	}
	return nil
}

func (rb *RecipeBot) Stop() error {
	log.Println("Closing", Name, "bye!")
	qErr := rb.msgQueue.Close()
	sErr := rb.session.Close()
	return errors.Join(qErr, sErr)
}

func (rb *RecipeBot) OnNewMessage(discord *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == discord.State.User.ID {
		return
	}
	results, count := urlextract.ExtractUrlsFromText(message.Content)
	rb.respondStatus(discord, message, results, count)
}

func (rb *RecipeBot) respondStatus(discord *discordgo.Session, message *discordgo.MessageCreate, results chan urlextract.WordResult, count int) {
	var response string
	for range count {
		result := <-results
		if result.UrlType != urlextract.NONE {
			response += reportResult(result)
			rb.msgQueue.SendMessage(result)
		}
	}
	sendReport(discord, message.ChannelID, response)
}

func reportResult(result urlextract.WordResult) string {
	return fmt.Sprintf("- %s: %s\n", result.UrlType, result.MatchedUrl.Hostname())
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
