# SlackBot 

Created to redirect shell and bash script messages from a server to a Slack Channel

## Example how to use
To run with the config `/etc/slackbot/config.yml`
````shell script
echo "[INFO] Info message\n[ERROR] Error message" | /usr/local/bin/slackbot
````

With config path
```shell script
echo "[ERROR] Error Message" | /usr/local/bin/slackbot -conf ./config.yml
```

## Config file
Create config file
```shell script
mkdir /etc/slackbot

touch /etc/slackbot/config.yml
```

Add to the file
```yaml
webhook:
  url: "https://hooks.slack.com/services/"
  secret: "XXXXXXXXXXX/YYYYYYYYY/Azv7tLC9y0yYHiB80AZMG"
```