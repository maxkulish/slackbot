package main

import (
	"flag"
	"log"
	"os"

	slackbot "github.com/maxkulish/slackbot/cmd"
	"github.com/maxkulish/slackbot/templates"
)

func main() {
	// Initialize logging
	log.SetPrefix("slackbot: ")
	log.SetFlags(0)

	// Try to read config path from environment variable first
	configPath := os.Getenv("SLACKBOT_CONFIG")

	c := slackbot.CMD{}

	if configPath == "" {
		// If the environment variable is not set, inform and use the flag or default value
		log.Println("SLACKBOT_CONFIG environment variable is empty; reading from the config file")
		flag.StringVar(&c.ConfigFile, "config", "/etc/yourproject/default.yml", "Path to the config file")
	} else {
		// Environment variable is set, use it directly and log the usage
		c.ConfigFile = configPath
		log.Printf("Reading config from SLACKBOT_CONFIG environment variable: %s", configPath)
	}

	flag.BoolVar(&c.Help, "help", false, templates.HelpMessage)
	flag.Parse()

	if c.ConfigFile == "" {
		log.Println("No config file specified. Using default settings.")
	} else {
		log.Printf("Using config file: %s", c.ConfigFile)
	}

	if err := c.Run(); err != nil {
		log.Fatalf("Failed to run: %v", err)
	}
}
