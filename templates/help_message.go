package templates

const HelpMessage = `SlackBot sends message to the Slack channel

echo "[ERROR] Some error details" | /usr/local/bin/slackbot

cat file.txt | slackbot

echo "Text message" | slackbot -config ./config.yml
`
