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
	log.Printf("channelid:%v", s.channelID)
	log.Printf("botid:%v", s.botID)
	log.Printf("channelid:%v", ev.Channel)
	log.Printf("botid:%v", ev.BotID)
	if ev.Channel != s.channelID {
		log.Printf("%s %s", ev.Channel, ev.Msg.Text)
		log.Println("ここまでchannelid")
		return nil
	}

	// Only response mention to bot. Ignore else.
	if !strings.HasPrefix(ev.Msg.Text, fmt.Sprintf("<@%s> ", s.botID)) {
		log.Println("ここまでbotid")
		return nil
	}

	// Parse message
	m := strings.Split(strings.TrimSpace(ev.Msg.Text), " ")[1:]
	// todo:ここにぱーす後の文字列処理を入れる
	log.Println("ここまでstart")
	if len(m) == 0 {
		// 何か言ってよと返す
		log.Println("no m")
		return fmt.Errorf("invalid message")
	}

	//switch m[0]で角処理に分岐させる
	switch m[0] {
	case "ghe":
		c, err := ghe.NewClient(gheurl)
		if err != nil {
			return err
		}

		log.Println("ここまで1")
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second) // 5秒後にキャンセル
		defer cancel()

		switch m[1] {
		case "pr":
			log.Println("ここまで2")
			repos, err := c.GetRepos(ctx, "tjtjtjtj")
			if err != nil {
				return err
			}
			for _, r := range repos {
				log.Printf("repo:%s", r.Full_name)
				//s.rtm.SendMessage(s.rtm.NewOutgoingMessage(r.Full_name, s.channelID))

				attachment := slack.Attachment{
					Title:    "Test",
					ImageURL: "https://www.google.com.tw/images/branding/googlelogo/2x/googlelogo_color_120x44dp.png",
				}
				params := slack.PostMessageParameters{
					Username:    "Log Reporter",
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
	/*
		// value is passed to message handler when request is approved.
		attachment := slack.Attachment{
			Text:       "Which beer do you want? :beer:",
			Color:      "#f9a41b",
			CallbackID: "beer",
			Actions: []slack.AttachmentAction{
				Text: "Which beer do you want? :beer:",
				Name: actionSelect,
				Type: "select",
				Options: []slack.AttachmentActionOption{
					{
						Text:  "Asahi Super Dry",
						Value: "Asahi Super Dry",
					},
					{
						Text:  "Kirin Lager Beer",
						Value: "Kirin Lager Beer",
					},
					{
						Text:  "Sapporo Black Label",
						Value: "Sapporo Black Label",
					},
					{
						Text:  "Suntory Malts",
						Value: "Suntory Malts",
					},
					{
						Text:  "Yona Yona Ale",
						Value: "Yona Yona Ale",
					},
				},
			},

			{
				Name:  actionCancel,
				Text:  "Cancel",
				Type:  "button",
				Style: "danger",
			},
		}

		params := slack.PostMessageParameters{
			Attachments: []slack.Attachment{
				attachment,
			},
		}

		if _, _, err := s.client.PostMessage(ev.Channel, "", params); err != nil {
			return fmt.Errorf("failed to post message: %s", err)
		}
	*/

	return nil
}
