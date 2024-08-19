# telegram-bot-file-proxy

## What
Proxies Telegram bot API requests to Telegram's file API.

## Why
Telegram's get file API (https://core.telegram.org/bots/api#getfile) uses links that you cannot share because they have your entire private key embedded.

This simple service fixes that problem for you and includes a Dockerfile for easy deployment.

## How
First start the docker container:
```bash
docker build -t telegram-bot-file-proxy .
docker run -p 5000:5000 -e TELEGRAM_BOT_TOKEN=your_bot_token_here telegram-bot-file-proxy
```

Or run it without docker:
```bash
go build -o api .
./api
```

Then call the API and the correct file contents will be streamed to the response without leaking your private key:
```bash
wget "http://localhost:5000/v1/telegram/file/your_file_id_here"
```
