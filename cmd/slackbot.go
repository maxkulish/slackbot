// Package slackbot provides functionality for sending notifications to Slack.
// It includes support for command line interaction, configuration management,
// message preparation based on the local system's information, and the actual
// sending of messages to a Slack webhook.
package slackbot

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/maxkulish/slackbot/config"
	"github.com/maxkulish/slackbot/localip"
	"github.com/maxkulish/slackbot/slack"
	"github.com/maxkulish/slackbot/templates"
)

// CMD represents the command line structure.
type CMD struct {
	ConFile string
	Help    bool
}

// Run executes the main application logic.
func (c *CMD) Run() {
	if c.Help {
		fmt.Println(templates.HelpMessage)
		os.Exit(0)
	}

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	ips, err := localip.GetLocalIPAddr()
	if err != nil {
		log.Fatal(err)
	}

	var inputText string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		inputText += "\n" + scanner.Text()
	}

	msg := slack.PrepareMessage(hostname, inputText, ips)

	conf, err := config.NewConfig(c.ConFile)
	if err != nil {
		log.Fatal(err)
	}

	err = slack.SendSlackNotification(conf.WebHook, msg)
	if err != nil {
		log.Fatal(err)
	}
}
