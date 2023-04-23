package midjourney

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const (
	baseURL = "https://discord.com/api/v9/interactions"
	version = "1077969938624553050"
	id      = "938956540159881230"
)

// ClientOptions contains the configuration options for the BotClient.
type ClientOptions struct {
	BotToken           string `json:"bot_token"`
	AuthorizationToken string `json:"authorization_token"`
	GuildId            string `json:"guild_id"`
	ChannelId          string `json:"channel_id"`
	ApplicationId      string `json:"application_id"`
}

// EventHandleFunc is a function type for handling message events.
type EventHandleFunc func(*discordgo.MessageCreate)

// BotClient is a simple Discord bot client using the DiscordGo library.
type BotClient struct {
	options         *ClientOptions
	dg              *discordgo.Session
	MessageId       string `json:"message_id,omitempty"`
	SessionId       string `json:"session_id,omitempty"`
	messageHandlers EventHandleFunc
}

var botClient *BotClient

// NewClient creates a new BotClient and connects it to Discord.
func NewClient(options *ClientOptions) (*BotClient, error) {
	botClient = &BotClient{
		options: options,
	}

	err := botClient.Connect()
	if err != nil {
		return nil, err
	}

	return botClient, nil
}

func GetBotClient() *BotClient {
	return botClient
}

// Connect establishes a connection to Discord and sets up the message event handler.
func (c *BotClient) Connect() error {
	if c.options.BotToken == "" {
		return errors.New("bot token is required")
	}

	dg, err := discordgo.New("Bot " + c.options.BotToken)
	if err != nil {
		return err
	}
	c.dg = dg
	c.dg.AddHandler(c.handleMessage)
	if err := c.dg.Open(); err != nil {
		return err
	}

	return nil
}

// Disconnect closes the connection to Discord.
func (c *BotClient) Disconnect() {
	c.dg.Close()
}

// handleMessage is the internal message handler for the Discord bot client.
func (c *BotClient) handleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	handler := c.messageHandlers
	handler(m)
}

func (c *BotClient) OnMessage(handleFunc EventHandleFunc) {
	c.messageHandlers = handleFunc
}

// SendMessage sends a message to the specified channel.
func (c *BotClient) SendMessage(channelID, message string) (*discordgo.Message, error) {
	if channelID == "" || message == "" {
		return nil, errors.New("channelID and message cannot be empty")
	}

	return c.dg.ChannelMessageSend(channelID, message)
}

// GenerateImage GenerateImage
func (c *BotClient) GenerateImage(prompt string) (int, []byte, error) {
	payload := Payload{
		Type:          2,
		ApplicationID: c.options.ApplicationId,
		GuildID:       c.options.GuildId,
		ChannelID:     c.options.ChannelId,
		SessionID:     generateSessionId(16),
		Data: Data{
			Version: version,
			ID:      id,
			Name:    "imagine",
			Type:    1,
			Options: []Options{
				{
					Type:  3,
					Name:  "prompt",
					Value: prompt,
				},
			},
			ApplicationCommand: ApplicationCommand{
				ID:                       id,
				ApplicationID:            c.options.ApplicationId,
				Version:                  version,
				DefaultMemberPermissions: nil,
				Type:                     1,
				Nsfw:                     false,
				Name:                     "imagine",
				Description:              "Create images with Midjourney",
				DmPermission:             true,
				Options: []ApplicationCommandOptions{
					{
						Type:        3,
						Name:        "prompt",
						Description: "The prompt to imagine",
						Required:    true,
					},
				},
			},
			Attachments: nil,
		},
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return 0, nil, err
	}
	return c.DoHTTPRequest(jsonPayload)

}

func (c *BotClient) Upscale(index uint8, messageId, customId string) (int, []byte, error) {
	payload := UpscalePayload{
		Type:          3,
		GuildId:       c.options.GuildId,
		ChannelId:     c.options.ChannelId,
		MessageFlags:  0,
		MessageId:     messageId,
		ApplicationId: c.options.ApplicationId,
		SessionId:     generateSessionId(16),
		Data: UpscaleData{
			ComponentType: 2,
			CustomId:      fmt.Sprintf("MJ::JOB::upsample::%d::%s", index, customId),
		},
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return 0, nil, err
	}
	return c.DoHTTPRequest(jsonPayload)
}

func (c *BotClient) Variation(index uint8, messageId, customId string) (int, []byte, error) {
	payload := UpscalePayload{
		Type:          3,
		GuildId:       c.options.GuildId,
		ChannelId:     c.options.ChannelId,
		MessageFlags:  0,
		MessageId:     messageId,
		ApplicationId: c.options.ApplicationId,
		SessionId:     generateSessionId(16),
		Data: UpscaleData{
			ComponentType: 2,
			CustomId:      fmt.Sprintf("MJ::JOB::variation::%d::%s", index, customId),
		},
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return 0, nil, err
	}
	return c.DoHTTPRequest(jsonPayload)
}

func (c *BotClient) Reset(messageId, customId string) (int, []byte, error) {
	payload := UpscalePayload{
		Type:          3,
		GuildId:       c.options.GuildId,
		ChannelId:     c.options.ChannelId,
		MessageFlags:  0,
		MessageId:     messageId,
		ApplicationId: c.options.ApplicationId,
		SessionId:     generateSessionId(16),
		Data: UpscaleData{
			ComponentType: 2,
			CustomId:      fmt.Sprintf("MJ::JOB::reroll::0::%s::SOLO", customId),
		},
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return 0, nil, err
	}
	return c.DoHTTPRequest(jsonPayload)
}

// DoHTTPRequest makes a generic HTTP request and returns the response body as a string.
func (c *BotClient) DoHTTPRequest(data []byte) (int, []byte, error) {
	client := &http.Client{}

	req, err := http.NewRequest("POST", baseURL, strings.NewReader(string(data)))
	if err != nil {
		return 0, nil, err
	}
	req.Header.Add("Authorization", c.options.AuthorizationToken)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36")
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, err
	}

	return resp.StatusCode, bodyBytes, nil
}

func generateSessionId(length int) string {
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(randomBytes)
}

func (c *BotClient) ExtractDataFromURL(url string) (string, error) {
	lastUnderscoreIndex := strings.LastIndex(url, "_")
	if lastUnderscoreIndex == -1 {
		return "", fmt.Errorf("no underscore found")
	}

	lastDotIndex := strings.LastIndex(url, ".")
	if lastDotIndex == -1 {
		return "", fmt.Errorf("no dot found")
	}

	return url[lastUnderscoreIndex+1 : lastDotIndex], nil
}
