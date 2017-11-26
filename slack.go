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
	gheurl = "https://api.github.com"
)

type SlackListener struct {
	client    *slack.Client
	botID     string
	channelID string
	rtm       *slack.RTM
}

// LstenAndResponse listens slack events and response particular messages.
func (s *SlackListener) ListenAndResponse() {
	s.rtm = s.client.NewRTM()

	go s.rtm.ManageConnection()

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
	if len(m) == 0 {
		s.rtm.SendMessage(s.rtm.NewOutgoingMessage("何か言ってよ", s.channelID))
		return nil
	}

	switch m[0] {
	case "ghe":
		if len(m) == 1 {
			s.rtm.SendMessage(s.rtm.NewOutgoingMessage("gheで何したい？", s.channelID))
			return nil
		}

		c, err := ghe.NewClient(gheurl, env.GHEToken)
		if err != nil {
			return err
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second) // 30秒後にキャンセル
		defer cancel()

		switch m[1] {
		case "pr-orgs":
			fallthrough
		case "pr-users":
			repos, err := c.GetRepos(ctx, strings.TrimLeft(m[1], "pr-"), m[2])
			if err != nil {
				return err
			}
			for _, r := range repos {
				log.Printf("repo:%s", r.Full_name)
				pulls, err := c.GetPulls(ctx, r.Owner.Login, r.Name)
				if err != nil {
					return err
				}
				if len(pulls) == 0 {
					continue
				}

				attachmentfields := make([]slack.AttachmentField, 2)
				for _, p := range pulls {
					reviews, err := c.GetReviews(ctx, p.Base.Repo.Owner.Login, p.Base.Repo.Name, fmt.Sprint(p.Number))
					if err != nil {
						return err
					}

					var assigneelist string
					for _, a := range p.Assignees {
						assigneelist = assigneelist + a.User + "\n"
					}
					attachmentfields[0] = slack.AttachmentField{"Assignees", assigneelist, true}

					var reviewlist string
					for _, r := range reviews {
						reviewlist = reviewlist + r.User.Login + ":" + r.State + "\n"
					}
					attachmentfields[1] = slack.AttachmentField{"Reviews", reviewlist, true}

					attachment := slack.Attachment{
						Color:      "#33bbff",
						AuthorName: p.Base.Repo.Full_name,
						AuthorLink: p.Base.Repo.Html_url,
						Title:      p.Title,
						TitleLink:  p.Html_url,
						ThumbURL:   "https://assets-cdn.github.com/images/modules/open_graph/github-mark.png",
						Fields:     attachmentfields,
					}
					params := slack.PostMessageParameters{
						Username:    "go-bot",
						Attachments: []slack.Attachment{attachment},
					}

					if _, _, err := s.client.PostMessage(ev.Channel, "", params); err != nil {
						return fmt.Errorf("failed to post message: %s", err)
					}
				}
			}
		default:
			s.rtm.SendMessage(s.rtm.NewOutgoingMessage("gheでその機能ないです", s.channelID))
		}

	default:
		s.rtm.SendMessage(s.rtm.NewOutgoingMessage("I don't understand", s.channelID))
	}
	return nil
}
