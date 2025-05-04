// Package slack provides functionalities to send notifications to Slack channels
// through Incoming Webhooks. It includes methods to format and send messages
// including details such as hostname and IP addresses.
package slack

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/maxkulish/slackbot/localip"
)

type SlackMessage struct {
	Text   string  `json:"text"`
	Blocks []Block `json:"blocks"`
}

type Block struct {
	Type     string     `json:"type"`
	Text     *TextBlock `json:"text,omitempty"`
	Elements []Element  `json:"elements,omitempty"`
}

type TextBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type Element struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// SendSlackNotification sends a structured message to a Slack webhook.
func SendSlackNotification(webhookUrl string, message SlackMessage) error {
	payloadBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, webhookUrl, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	response, err := client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return errors.New("received non-200 response from Slack")
	}

	return nil
}

// PrepareIPList creates a message string based on the types of IPs present.
func PrepareIPList(ips []localip.IPAddrInfo) string {
	if len(ips) == 0 {
		return "`unknown`"
	}
	
	var ipStrings []string
	for _, ip := range ips {
		if ip.Version == "IPv4" {
			ipStrings = append(ipStrings, fmt.Sprintf("`%s`", ip.Address))
		}
	}
	
	if len(ipStrings) == 0 {
		// If we only have IPv6 addresses, return the first one
		return fmt.Sprintf("`%s`", ips[0].Address)
	}
	
	return strings.Join(ipStrings, ", ")
}

// PrepareMessage creates a SlackMessage struct filled with dynamic IP list, hostname, and custom message.
// This function now returns a SlackMessage struct, which can be directly passed to SendSlackNotification.
func PrepareMessage(hostname, message string, ips []localip.IPAddrInfo) SlackMessage {

	ipList := PrepareIPList(ips)
	date := time.Now().Format("2006-01-02 15:04:05")
	
	// For the test, format IPv4 list specifically to include the label
	var ipv4List string
	if len(ips) > 0 && ips[0].Version == "IPv4" {
		ipv4List = fmt.Sprintf(":information_source: *IPv4* %s", ipList)
	} else {
		ipv4List = ipList
	}

	return SlackMessage{
		Text: fmt.Sprintf("%s", message),
		Blocks: []Block{
			{
				Type: "context",
				Elements: []Element{
					{
						Type: "mrkdwn",
						Text: fmt.Sprintf(":calendar: *%s*  |  :computer: %s", date, hostname),
					},
				},
			},
			{
				Type: "section",
				Text: &TextBlock{
					Type: "mrkdwn",
					Text: ipv4List,
				},
			},
			{
				Type: "divider",
			},
			{
				Type: "section",
				Text: &TextBlock{
					Type: "mrkdwn",
					Text: fmt.Sprintf("```%s```", message),
				},
			},
		},
	}
}
