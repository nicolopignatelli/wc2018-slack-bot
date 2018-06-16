# World Cup 2018 Slack Bot
A simple Slack bot for notifications about the current WC2018 match.

It will notify about the start and end of a match, yellow and red cards, half times and, naturally, goals.

# Setup and installation

- Create a Slack bot in your account and note down the token
- Invite the bot in your channel
- Clone the repo
- Run `make build`. An x64 executable is created in the bin/ directory.
  Adjust the target according to your runtime architecture.
- Run the binary with the required env variables.
  Eg: `WC2018_SLACK_TOKEN=<your-token> WC2018_SLACK_CHANNEL=<your-channel> ./wc2018-slack-bot`
