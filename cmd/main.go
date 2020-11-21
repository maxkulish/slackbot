package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/maxkulish/slackbot/templates"

	"github.com/maxkulish/slackbot/config"

	"github.com/maxkulish/slackbot/localip"
)

type SlackRequestBody struct {
	Text string `json:"text"`
}

type WebHook struct {
	URL    string
	Secret string
}

func main() {

	conFile := flag.String("conf", "/etc/slackbot/config.yml", "Path to the config.yml file")
	help := flag.Bool("help", false, templates.HelpMessage)
	flag.Parse()

	if *help {
		fmt.Println(templates.HelpMessage)
		os.Exit(0)
	}

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	ips, err := localip.GetIPv4()
	if err != nil {
		log.Fatal(err)
	}

	var inputText string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		inputText += "\n" + scanner.Text()
	}

	msg := prepareMessage(hostname, inputText, ips)

	conf, err := config.NewConfig(*conFile)
	if err != nil {
		log.Fatal(err)
	}

	err = SendSlackNotification(conf.WebHook.URL+conf.WebHook.Secret, msg)
	if err != nil {
		log.Fatal(err)
	}
}

// SendSlackNotification will post to an 'Incoming Webook' url setup in Slack Apps. It accepts
// some text and the slack channel is saved within Slack.
func SendSlackNotification(webhookUrl string, msg string) error {

	slackBody, _ := json.Marshal(SlackRequestBody{Text: msg})
	req, err := http.NewRequest(http.MethodPost, webhookUrl, bytes.NewBuffer(slackBody))
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
		return errors.New("Non-ok response returned from Slack")
	}
	return nil
}

func prepareMessage(hostname, text string, ips []string) string {
	var strIPs string

	if len(ips) == 0 {
		strIPs = fmt.Sprintf("ip: `unknown`")
	} else if len(ips) == 1 {
		strIPs = fmt.Sprintf("IPv4: `%s`", ips[0])
	} else {
		strIPs += fmt.Sprintf("IPv4: ")
		for _, ip := range ips {
			strIPs += fmt.Sprintf("`%s` ", ip)
		}
	}

	return fmt.Sprintf(
		":rotating_light: hostname: `%s`, %s ```%s```",
		hostname,
		strIPs,
		text)
}
