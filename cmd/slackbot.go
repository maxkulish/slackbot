// Package slackbot provides functionality for sending notifications to Slack.
// It includes support for command line interaction, configuration management,
// message preparation based on the local system's information, and the actual
// sending of messages to a Slack webhook.
package slackbot

import (
	"bufio"
	"fmt"
	"os"

	"github.com/maxkulish/slackbot/config"
	"github.com/maxkulish/slackbot/localip"
	"github.com/maxkulish/slackbot/slack"
	"github.com/maxkulish/slackbot/templates"
)

type CMD struct {
	ConfigFile string
	Help       bool
}

func (c *CMD) Run() error {
	if c.Help {
		fmt.Println(templates.HelpMessage)
		return nil
	}

	hostname, err := c.getHostname()
	if err != nil {
		return fmt.Errorf("failed to get hostname: %w", err)
	}

	ips, err := localip.GetLocalIPAddr()
	if err != nil {
		return fmt.Errorf("failed to get IP addresses: %w", err)
	}

	inputText, err := c.readInputText()
	if err != nil {
		return fmt.Errorf("failed to read input text: %w", err)
	}

	msg := slack.PrepareMessage(hostname, inputText, ips)

	conf, err := config.NewConfig(c.ConfigFile)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	if err = slack.SendSlackNotification(conf.WebHook, msg); err != nil {
		return fmt.Errorf("failed to send Slack notification: %w", err)
	}

	return nil
}

func (c *CMD) getHostname() (string, error) {
	return os.Hostname()
}

func (c *CMD) readInputText() (string, error) {
	fileInfo, err := os.Stdin.Stat()
	if err != nil {
		return "", fmt.Errorf("failed to stat stdin: %w", err)
	}

	// Check if data is available on stdin (e.g., piped input or redirect)
	if fileInfo.Mode()&os.ModeCharDevice == 0 {
		var inputText string
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			inputText += "\n" + scanner.Text()
		}
		if err := scanner.Err(); err != nil {
			return "", fmt.Errorf("error reading stdin: %w", err)
		}
		return inputText, nil
	}

	// No data available on stdin; don't block waiting for input
	return "", nil
}
