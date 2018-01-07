package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/kelseyhightower/envconfig"
	"github.com/nlopes/slack"
	"gopkg.in/robfig/cron.v2"
)

// https://api.slack.com/slack-apps
// https://api.slack.com/internal-integrations
type envConfig struct {

	// BotToken is bot user token to access to slack API.
	BotToken string `envconfig:"BOT_TOKEN" required:"true"`

	// BotID is bot user ID.
	BotID string `envconfig:"BOT_ID" required:"true"`

	// ChannelID is slack channel ID where bot is working.
	// Bot responses to the mention in this channel.
	ChannelID string `envconfig:"CHANNEL_ID" required:"true"`

	// GHEToken is bot user token to access to GHE API.
	GHEToken string `envconfig:"GHE_TOKEN" required:"false"`

	// ZabbixUser is user to access to Zabbix API.
	ZabbixUser string `envconfig:"ZABBIX_USER" required:"false"`

	// ZabbixPasswd is password for zabbix user.
	ZabbixPasswd string `envconfig:"ZABBIX_PASSWD" required:"false"`
}

var env envConfig

func main() {
	os.Exit(_main(os.Args[1:]))
}

func _main(args []string) int {
	if err := envconfig.Process("", &env); err != nil {
		log.Printf("[ERROR] Failed to process env var: %s", err)
		return 1
	}

	log.Printf("env:%v", env)

	// Listening slack event and response
	log.Printf("[INFO] Start slack event listening")
	client := slack.New(env.BotToken)
	slackListener := &SlackListener{
		client:    client,
		botID:     env.BotID,
		channelID: env.ChannelID,
	}

	go slackListener.ListenAndResponse()

	c := cron.New()
	if _, err := c.AddJob("0 30 * * * *", slackListener); err != nil {
		log.Printf("[ERROR] Failed to AddJob: %s", err)
		return 1
	}
	c.Start()
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig

	return 0
}
