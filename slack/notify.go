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
	var publicIPs, localIPs []string

	for _, ip := range ips {
		if ip.Version == "IPv4" {
			if ip.Local {
				localIPs = append(localIPs, fmt.Sprintf("`%s`", ip.Address))
			} else {
				publicIPs = append(publicIPs, fmt.Sprintf("`%s`", ip.Address))
			}
		}
	}

	var message strings.Builder
	if len(publicIPs) > 0 {
		message.WriteString(fmt.Sprintf(":globe_with_meridians: %s\n", strings.Join(publicIPs, ", ")))
	}
	if len(localIPs) > 0 {
		if message.Len() > 0 {
			message.WriteString(" :house: ")
		}
		message.WriteString(strings.Join(localIPs, ", "))
	}

	if message.Len() == 0 {
		return "`unknown`" // Fallback message if no IPs are provided or they're all IPv6.
	}
	return message.String()
}

// PrepareMessage creates a SlackMessage struct filled with dynamic IP list, hostname, and custom message.
// This function now returns a SlackMessage struct, which can be directly passed to SendSlackNotification.
func PrepareMessage(hostname, message string, ips []localip.IPAddrInfo) SlackMessage {

	ipList := PrepareIPList(ips)
	date := time.Now().Format("2006-01-02 15:04:05")

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
					Text: ipList,
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
