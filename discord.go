package main

import (
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Some Discord environment variables.
var discordToken string
var discordStatusChannel string
var discordUpdateChannel string
var currentStatus Status
var currentUpdate Update

// Current status of the website.
type Status struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// Latest update from the Discord 'updates' channel.
type Update struct {
	ID           string `json:"id"`
	GuildID      string `json:"guildId"`
	ChannelID    string `json:"channelId"`
	CreatedTs    int64  `json:"createdTimestamp"`
	EditedTs     int64  `json:"editedTimestamp"`
	AuthorID     string `json:"authorId"`
	AuthorName   string `json:"authorName"`
	AuthorImage  string `json:"authorImage"`
	Content      string `json:"content"`
	CleanContent string `json:"cleanContent"`
	Image        string `json:"image"`
}

// Start the Discord bot responsible for the site status and 'What's new?' card.
func startDiscordBot() {
	// create Discord session
	dg, err := discordgo.New("Bot " + discordToken)
	if err != nil {
		log.Fatalf("failed to create Discord session: %s", err)
	}

	// set intents (all we want is guild messages)
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// get current status
	statusMessages, err := dg.ChannelMessages(discordStatusChannel, 0, "", "", "")
	if err != nil {
		log.Fatalf("failed getting current status: %s", err)
	}
	for i := 0; i < len(statusMessages); i++ {
		m := statusMessages[i]
		if strings.HasPrefix(m.Content, "--status-set") {
			currentStatus = Status{
				Type: "warn",
				Text: strings.TrimSpace(strings.Replace(m.Content, "--status-set", "", 1)),
			}
			break
		} else if strings.HasPrefix(m.Content, "--status-remove") {
			currentStatus = Status{
				Type: "empty",
			}
			break
		}
	}

	// get latest update
	updateMessages, err := dg.ChannelMessages(discordUpdateChannel, 0, "", "", "")
	if err != nil {
		log.Fatalf("failed getting latest update: %s", err)
	}
	for i := 0; i < len(updateMessages); i++ {
		m := updateMessages[i]
		if len(m.Attachments) > 0 {
			currentUpdate = Update{
				ID:           m.ID,
				GuildID:      m.GuildID,
				ChannelID:    m.ChannelID,
				CreatedTs:    m.Timestamp.UnixMilli(),
				EditedTs:     m.Timestamp.UnixMilli(),
				AuthorID:     m.Author.ID,
				AuthorName:   m.Author.Username,
				AuthorImage:  m.Author.AvatarURL(""),
				Content:      m.Content,
				CleanContent: m.ContentWithMentionsReplaced(),
				Image:        m.Attachments[0].URL,
			}
			break
		}
	}

	// messageCreate handler
	dg.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.ChannelID == discordStatusChannel {
			log.Printf("updating status by request of %s", m.Author.Username)
			if strings.HasPrefix(m.Content, "--status-set") {
				currentStatus = Status{
					Type: "warn",
					Text: strings.TrimSpace(strings.Replace(m.Content, "--status-set", "", 1)),
				}
				s.MessageReactionAdd(m.ChannelID, m.ID, "<:good:1118293837773807657>")
			} else if strings.HasPrefix(m.Content, "--status-remove") {
				currentStatus = Status{
					Type: "empty",
				}
				s.MessageReactionAdd(m.ChannelID, m.ID, "<:good:1118293837773807657>")
			}
		} else if m.ChannelID == discordUpdateChannel && len(m.Attachments) > 0 {
			log.Printf("updating latest update by request of %s", m.Author.Username)
			currentUpdate = Update{
				ID:           m.ID,
				GuildID:      m.GuildID,
				ChannelID:    m.ChannelID,
				CreatedTs:    m.Timestamp.UnixMilli(),
				EditedTs:     m.Timestamp.UnixMilli(),
				AuthorID:     m.Author.ID,
				AuthorName:   m.Author.Username,
				AuthorImage:  m.Author.AvatarURL(""),
				Content:      m.Content,
				CleanContent: m.ContentWithMentionsReplaced(),
				Image:        m.Attachments[0].URL,
			}
		}
	})

	// connect to Discord
	err = dg.Open()
	if err != nil {
		log.Fatalf("failed to open Discord connection: %s", err)
	}
	defer dg.Close()
}
