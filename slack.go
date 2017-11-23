package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/nlopes/slack"
	"github.com/tjtjtjtj/go-bot/ghe"
)

const (
	// action is used for slack attament action.
	actionSelect = "select"
	actionStart  = "start"
	actionCancel = "cancel"

	gheurl = "https://api.github.com"
)

type SlackListener struct {
	client    *slack.Client
	botID     string
	channelID string
	rtm       *slack.RTM
}

// LstenAndResponse listens slack events and response
// particular messages. It replies by slack message button.
func (s *SlackListener) ListenAndResponse() {
	s.rtm = s.client.NewRTM()

	// Start listening slack events
	go s.rtm.ManageConnection()

	// Handle slack events
	for msg := range s.rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			if err := s.handleMessageEvent(ev); err != nil {
				log.Printf("[ERROR] Failed to handle message: %s", err)
			}
		}
	}
}

// handleMesageEvent handles message events.
func (s *SlackListener) handleMessageEvent(ev *slack.MessageEvent) error {
	// Only response in specific channel. Ignore else.
	if ev.Channel != s.channelID {
		return nil
	}

	// Only response mention to bot. Ignore else.
	if !strings.HasPrefix(ev.Msg.Text, fmt.Sprintf("<@%s>", s.botID)) {
		return nil
	}

	// Parse message
	m := strings.Split(strings.TrimSpace(ev.Msg.Text), " ")[1:]
	// todo:ここにぱーす後の文字列処理を入れる
	log.Println("ここまでstart")
	if len(m) == 0 {
		s.rtm.SendMessage(s.rtm.NewOutgoingMessage("何か言ってよ", s.channelID))
		log.Println("no m")
		return fmt.Errorf("invalid message")
	}

	//switch m[0]で角処理に分岐させる
	switch m[0] {
	case "ghe":
		if len(m) == 1 {
			s.rtm.SendMessage(s.rtm.NewOutgoingMessage("gheで何したい？", s.channelID))
			return nil
		}

		c, err := ghe.NewClient(gheurl)
		if err != nil {
			return err
		}

		log.Println("ここまで1")
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, 20*time.Second) // 5秒後にキャンセル
		defer cancel()

		switch m[1] {
		case "pr":
			//ここpr の次に sreとかorg指定でもいいかも
			log.Println("ここまで2")
			repos, err := c.GetRepos(ctx, "tjtjtjtj")
			if err != nil {
				return err
			}
			for _, r := range repos {
				log.Printf("repo:%s", r.Full_name)
				pulls, err := c.GetPulls(ctx, "tjtjtjtj", r.Name)
				if err != nil {
					return err
				}
				if len(pulls) == 0 {
					continue
				}

				attachmentfields := make([]slack.AttachmentField, 0)
				for _, p := range pulls {
					//ここfielsに埋める
					attachmentfields = append(attachmentfields, slack.AttachmentField{r.Full_name, fmt.Sprint(p.Number), false})
				}

				attachment := slack.Attachment{
					Title:     r.Full_name,
					TitleLink: r.Html_url,
					ThumbURL:  "https://assets-cdn.github.com/images/modules/open_graph/github-mark.png",
					Fields:    attachmentfields,
				}
				params := slack.PostMessageParameters{
					Username:    "go-bot",
					Attachments: []slack.Attachment{attachment},
				}

				if _, _, err := s.client.PostMessage(ev.Channel, "", params); err != nil {
					return fmt.Errorf("failed to post message: %s", err)
				}
			}
		default:
			log.Println("ここまで3")
		}

	default:
		log.Println("ここまでend")
	}
	return nil
}
