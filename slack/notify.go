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

// SendSlackNotification will post to an 'Incoming Webhook' url setup in Slack Apps. It accepts
// some text and the slack channel is saved within Slack.
// func SendSlackNotificationOld(webhook string, msg string) error {

// 	slackBody, _ := json.Marshal(SlackRequestBody{Text: msg})
// 	req, err := http.NewRequest(http.MethodPost, webhook, bytes.NewBuffer(slackBody))
// 	if err != nil {
// 		return err
// 	}

// 	req.Header.Add("Content-Type", "application/json")

// 	client := &http.Client{Timeout: 10 * time.Second}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return err
// 	}

// 	buf := new(bytes.Buffer)
// 	buf.ReadFrom(resp.Body)
// 	if buf.String() != "ok" {
// 		return errors.New("non-ok response returned from Slack")
// 	}
// 	return nil
// }

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

// PrepareIPList formats a slice of IPAddrInfo into a string.
// For an empty slice, it returns "`unknown`". For a single IP address, it returns that IP
// wrapped in backticks. For multiple IPs, it filters and formats only the IPv4 addresses,
// separating them with commas and enclosing each in backticks. IPv6 addresses are ignored.
func PrepareIPList(ips []localip.IPAddrInfo) string {
	var strIPs strings.Builder

	switch len(ips) {
	case 0:
		strIPs.WriteString("`unknown`")
	case 1:
		strIPs.WriteString(fmt.Sprintf("`%s`", ips[0].Address))
	default:
		firstIPv4Added := false
		for _, ip := range ips {
			if ip.Version == "IPv4" {
				if firstIPv4Added {
					strIPs.WriteString(", ") // Add comma and space before each entry except the first
				}
				strIPs.WriteString(fmt.Sprintf("`%s`", ip.Address))
				firstIPv4Added = true
			}
		}
	}

	return strIPs.String()
}

// PrepareMessage creates a SlackMessage struct filled with dynamic IP list, hostname, and custom message.
// This function now returns a SlackMessage struct, which can be directly passed to SendSlackNotification.
func PrepareMessage(hostname, message string, ips []localip.IPAddrInfo) SlackMessage {

	ipList := PrepareIPList(ips)

	return SlackMessage{
		Text: "Danny Torrence left a 1 star review for your property.",
		Blocks: []Block{
			{
				Type: "context",
				Elements: []Element{
					{
						Type: "mrkdwn",
						Text: fmt.Sprintf(":calendar: *November 12, 2019*  |  :computer: %s", hostname),
					},
				},
			},
			{
				Type: "divider",
			},
			{
				Type: "section",
				Text: &TextBlock{
					Type: "mrkdwn",
					Text: fmt.Sprintf(":information_source: *IPv4* %s", ipList),
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
