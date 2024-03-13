package mittelweg

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bernardoforcillo/diskursus"
)

const (
	ApplicationID string = "936929561302675456"
	SessionID     string = "2fb980f65e5c9a77c96ca01f2c242cf6"
)

type MidjourneyClient struct {
	client    diskursus.DiscordClient
	channelID string
}

func NewClient(authToken string, channelID string) *MidjourneyClient {
	client := diskursus.NewClient(authToken)
	return &MidjourneyClient{
		client:    *client,
		channelID: channelID,
	}
}

func (c *MidjourneyClient) Imagine(prompt string, waitUntilGenerated bool) (*GeneratedImage, error) {
	channel, err := c.client.Channel(c.channelID)
	if err != nil {
		return nil, err
	}
	interactionsReq := &diskursus.InteractionRequest{
		Type:          2,
		ApplicationID: ApplicationID,
		GuildID:       channel.GuildID,
		ChannelID:     c.channelID,
		SessionID:     SessionID,
		Data: map[string]any{
			"version": "1166847114203123795",
			"id":      "938956540159881230",
			"name":    "imagine",
			"type":    "1",
			"options": []map[string]any{
				{
					"type":  3,
					"name":  "prompt",
					"value": prompt,
				},
			},
			"application_command": map[string]any{
				"id":                         "938956540159881230",
				"application_id":             ApplicationID,
				"version":                    "1166847114203123795",
				"default_permission":         true,
				"default_member_permissions": nil,
				"type":                       1,
				"nsfw":                       false,
				"name":                       "imagine",
				"description":                "Create images with Midjourney",
				"dm_permission":              true,
				"options": []map[string]any{
					{
						"type":        3,
						"name":        "prompt",
						"description": "The prompt to imagine",
						"required":    true,
					},
				},
				"attachments": []any{},
			},
		},
	}
	err = c.client.SendInteraction(interactionsReq)
	if err != nil {
		return nil, err
	}
	if waitUntilGenerated {
		time.Sleep(8 * time.Second)
		var txtmessage string
		var result *diskursus.Message
		for txtmessage != "fast" {
			result, err = c.SearchGeneratedMessage(prompt)
			if err != nil {
				return nil, err
			}
			regexPattern := `^\*\*(.*) --v\s+\d+(\.\d+)?\*\*\s*-\s*<@\d+> \((.*?)\)$`
			regex, err := regexp.Compile(regexPattern)
			if err != nil {
				return nil, err
			}
			matches := regex.FindStringSubmatch(result.Content)
			if regex.MatchString(result.Content) && len(matches) == 4 && matches[1] == prompt {
				txtmessage = matches[3]
			}
			time.Sleep(4 * time.Second)
		}
		return &GeneratedImage{mj: c, message: *result}, nil
	}
	return nil, nil
}

func (c MidjourneyClient) SearchGeneratedMessage(prompt string) (*diskursus.Message, error) {
	messages, err := c.client.Messages(c.channelID)
	if err != nil {
		return nil, err
	}
	var result diskursus.Message
	regex := regexp.MustCompile(`^\*\*(.*) --v\s+\d+(\.\d+)?\*\*\s*-\s*<@\d+> \((.*?)\)$`)
	currentTime := time.Now()
	for _, message := range *messages {
		creationTime := message.Timestamp
		if currentTime.Sub(creationTime) < 2*time.Minute && regex.MatchString(message.Content) {
			matches := regex.FindStringSubmatch(message.Content)
			if strings.Compare(strings.Trim(matches[1], " "), strings.Trim(prompt, " ")) == 0 {
				result = message
				break
			}
		}
	}
	return &result, nil
}

func (c MidjourneyClient) SearchUpscaledMessage(prompt string, image int) (*diskursus.Message, error) {
	messages, err := c.client.Messages(c.channelID)
	if err != nil {
		return nil, err
	}
	var matches []string
	var result diskursus.Message
	regex := regexp.MustCompile(`^\*\*(.*) --v\s+\d+(\.\d+)?\*\*\s*- Image #(\d+) \s*<@\d+>$`)
	for _, message := range *messages {
		if regex.MatchString(message.Content) {
			matches = regex.FindStringSubmatch(message.Content)
			if strings.Compare(strings.Trim(matches[1], " "), strings.Trim(prompt, " ")) == 0 && matches[3] == strconv.Itoa(image) {
				result = message
				break
			}
		}
	}
	return &result, nil
}

type GeneratedImage struct {
	mj      *MidjourneyClient
	message diskursus.Message
}

func (g GeneratedImage) URL() string {
	return g.message.Attachments[0].Url
}

func (g *GeneratedImage) Upscale(index int, waitUntilGenerated bool) (*UpscaledImage, error) {
	if index < 0 || index > 3 {
		return nil, fmt.Errorf("upscale index must be between 0 and 3")
	}
	channel, err := g.mj.client.Channel(g.mj.channelID)
	if err != nil {
		return nil, err
	}
	flags := 0
	interactionsReq := &diskursus.InteractionRequest{
		Type:          3,
		ApplicationID: ApplicationID,
		GuildID:       channel.GuildID,
		ChannelID:     channel.ID,
		MessageFlags:  &flags,
		MessageID:     &g.message.ID,
		SessionID:     SessionID,
		Data: map[string]any{
			"component_type": 2,
			"custom_id":      g.message.Components[0].Components[index].CustomID,
		},
	}
	err = g.mj.client.SendInteraction(interactionsReq)
	if err != nil {
		return nil, err
	}
	if waitUntilGenerated {
		time.Sleep(4 * time.Second)
	}
	prompt := strings.Split(g.message.Content, "**")[1]
	regex := regexp.MustCompile(`\s*--v\s+\d+(\.\d+)?\s*`)
	prompt = regex.ReplaceAllString(prompt, "")
	result, err := g.mj.SearchUpscaledMessage(prompt, index + 1)
	if err != nil {
		return nil, err
	}
	upscaled := UpscaledImage{
		mj:      g.mj,
		message: *result,
	}
	return &upscaled, nil
}

type UpscaledImage struct {
	mj      *MidjourneyClient
	message diskursus.Message
}

func (u UpscaledImage) Message() diskursus.Message {
	return u.message
}

func (u UpscaledImage) URL() string {
	return u.message.Attachments[0].Url
}
