# RSS feeds to your discord webhook URL

## Setup

Install

```
go install github.com/arran4/discord-rss-webhook
```

Setup Scheduled task:
```bash
% cd ~/.config/systemd/user/

% cat > rss-discord-cron.timer 
[Unit]
Description=Run rss-discord-cron daily

[Timer]
OnCalendar=daily
Persistent=true
ExecStart=/home/$USER/go/bin/discord-rss-cron

[Install]
WantedBy=timers.target
^D

% systemctl enable --user rss-discord-cron.timer
Created symlink ~/.config/systemd/user/timers.target.wants/rss-discord-cron.service → ~/.config/systemd/user/rss-discord-cron.service.

```

Configure, get the webhook URL: https://support.discord.com/hc/en-us/articles/228383668-Intro-to-Webhooks
```
vim ~/.config/discord-npr-tdc-rss-webhook/channel-1
```

With

```
{
  "HookUrl": "https://discord.com/api/webhooks/......",
  "FeedUrl": "https://www.youtube.com/feeds/videos.xml?playlist_id=....."
}
```

`^D` means press `CTRL+D`

That should be it.