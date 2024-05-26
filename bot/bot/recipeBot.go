package bot

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"recipebot/environment"
	"recipebot/urlProcessor"
	"recipebot/urlextract"
)

const (
	Name = "RecipeBot"
)

type RecipeBot struct {
	config     *environment.RecipeBotConfig
	session    *discordgo.Session
	processors []urlProcessor.UrlProcessor
}

func (rb *RecipeBot) Configure(config *environment.RecipeBotConfig) *RecipeBot {
	rb.config = config
	return rb
}

func (rb *RecipeBot) AddProcessor(p urlProcessor.UrlProcessor) *RecipeBot {
	if rb.processors == nil {
		rb.processors = make([]urlProcessor.UrlProcessor, 0, 5)
	}
	rb.processors = append(rb.processors, p)
	return rb
}

func (rb *RecipeBot) Start() error {
	log.Println("Starting", Name)
	var err error
	rb.session, err = discordgo.New("Bot " + rb.config.BotToken)
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
	sErr := rb.session.Close()
	return errors.Join(sErr)
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
			processed, err := rb.processUrl(result)
			response += reportInitialResult(result, processed, err)
		}
	}
	sendReport(discord, message.ChannelID, response)
}

func reportInitialResult(result urlextract.WordResult, processed *urlProcessor.Result, err error) string {
	if err != nil {
		return fmt.Sprintf("- %s: %s\n", result.MatchedUrl.Hostname(), err)
	}
	return fmt.Sprintf("- [%s](%s)\n", result.MatchedUrl.Hostname(), processed.StoredUrl)
}

func (rb *RecipeBot) processUrl(exUrl urlextract.WordResult) (*urlProcessor.Result, error) {
	for _, processor := range rb.processors {
		if processor.CanHandle(exUrl.UrlType) {
			request := new(urlProcessor.Request)
			request.Details = exUrl
			return processor.Process(request)
		}
	}
	return nil, fmt.Errorf("cannot process URL with type %s", exUrl.UrlType)
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
