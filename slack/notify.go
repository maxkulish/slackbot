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

type SlackRequestBody struct {
	Text string `json:"text"`
}

// SendSlackNotification will post to an 'Incoming Webhook' url setup in Slack Apps. It accepts
// some text and the slack channel is saved within Slack.
func SendSlackNotification(webhook string, msg string) error {

	slackBody, _ := json.Marshal(SlackRequestBody{Text: msg})
	req, err := http.NewRequest(http.MethodPost, webhook, bytes.NewBuffer(slackBody))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	if buf.String() != "ok" {
		return errors.New("non-ok response returned from Slack")
	}
	return nil
}

// PrepareMessage formats a message string including the hostname, IP addresses, and provided text.
// It returns a string ready for sending as a Slack message.
// The function handles various cases for IP address availability and formats them accordingly.
func PrepareMessage(hostname, text string, ips []localip.IPAddrInfo) string {
	var strIPs strings.Builder

	switch len(ips) {
	case 0:
		strIPs.WriteString("ip: `unknown`")
	case 1:
		strIPs.WriteString(fmt.Sprintf("IP: `%s`", ips[0]))
	default:
		strIPs.WriteString("IPs: ")
		for i, ip := range ips {
			if i > 0 && ip.Version == "IPv4" {
				strIPs.WriteString(", ") // Add comma and space before each entry except the first
				strIPs.WriteString(fmt.Sprintf("`%s`", ip.Address))
			}
		}
	}

	return fmt.Sprintf(
		":rotating_light: hostname: `%s`, %s ```%s```",
		hostname,
		strIPs.String(),
		text,
	)
}
