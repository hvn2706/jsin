# jsin
A telegram chatbot to query your own hosted image

## Common flow
- Receive command -> query db and get s3 url -> get object -> send object to telegram

## Setup
- Create config.yml from config.tmp.yml and replace with your configuration

### S3
- Currently using [cloudflare r2](https://www.cloudflare.com/developer-platform/products/r2/)
- Replace s3 config in config.yml
- Use any s3 provider if you want, and add your own interface in ./external/s3

### Mysql
- Init mysql and create folder ./objects inside repo, put your images inside that folder
- Replace mysql config in config.yml
- Run: ```go run ./cmd/main.go jsin-migration``` and it will migrate your images to bucket

### Telegram bot
- Create your telegram bot with [BotFather](https://core.telegram.org/bots/tutorial), remember to allow it to send message, group settings, ...
- Replace your bot config in config.yml
- Run: ```go run ./cmd/main.go jsin-telegram``` and your bot is up and running
