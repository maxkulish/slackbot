# SlackBot

Created to redirect shell and bash script messages from a server to a Slack Channel

## Example how to use

To run with the config `/etc/slackbot/config.yml`

```shell script
echo "[INFO] Info message\n[ERROR] Error message" | slackbot
```

With config path

```shell script
echo "[ERROR] Error Message" | slackbot -conf ./config.yml
```

## Config file

Create config file

```shell script
mkdir /etc/slackbot

touch /etc/slackbot/config.yml
```

Add to the file

```yaml
webhook: "https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX"
```
