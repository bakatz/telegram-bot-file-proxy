# telegram-bot-file-proxy

## What
Proxies Telegram bot API requests to Telegram's [get file API](https://core.telegram.org/bots/api#getfile). This is meant for short-term use cases i.e. maybe you need to pass the file to a downstream system and don't want to waste money/disk space hosting the same image yourself.

If you need long term image storage, it's better to upload the file from Telegram to something like Cloudflare R2 or Amazon S3. However, Telegram file links are always accessible for "at least 1 hour" according to their docs, which is enough for many use cases.

## Why
Telegram's [get file API](https://core.telegram.org/bots/api#getfile) uses links that you cannot share or pass downstream because they have your entire private key embedded.

This simple service fixes that problem for you and includes a Dockerfile for easy deployment.

## How
First start the docker container:
```bash
docker build -t telegram-bot-file-proxy .
docker run -p 5000:5000 -e TELEGRAM_BOT_TOKEN=your_bot_token_here telegram-bot-file-proxy
```

Or run it without docker:
```bash
go build -o proxy .
./proxy
```

Then call the API and the correct file contents will be streamed to the response without leaking your private key:
```bash
wget http://localhost:5000/v1/telegram/files/your_file_id_here
```
