// Package templates provides predefined templates for SlackBot messages.
// It includes templates for help messages and configurations.
package templates

const HelpMessage = `SlackBot sends message to the Slack channel

echo "[ERROR] Some error details" | slackbot

cat file.txt | slackbot

echo "Text message" | slackbot -config ./config.yml`
