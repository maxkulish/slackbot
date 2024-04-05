package main

import (
	"flag"

	slackbot "github.com/maxkulish/slackbot/cmd"
	"github.com/maxkulish/slackbot/templates"
)

func main() {
	c := slackbot.CMD{}
	// Use "config" as the flag instead of "conf"
	flag.StringVar(&c.ConFile, "config", "/etc/slackbot/config.yml", "Path to the config.yml file")
	flag.BoolVar(&c.Help, "help", false, templates.HelpMessage)
	flag.Parse()

	c.Run()
}
