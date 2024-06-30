# RSS feeds to your discord webhook URL

## Setup

Install

```
go install github.com/arran4/discord-rss-webhook/cmd/discord-rss-cron
```

Setup Scheduled task:
```bash
% cd ~/.config/systemd/user/

% cat > discord-rss-cron.timer 
[Unit]
Description=Run discord-rss-cron daily

[Timer]
OnCalendar=daily
Persistent=true

[Install]
WantedBy=timers.target
^D

% cat ~/.config/systemd/user/discord-rss-cron.service
[Unit]
Description=Run discord-rss-cron
DefaultDependencies=no
Conflicts=shutdown.target
After=local-fs.target time-sync.target
Before=shutdown.target

[Service]
Type=oneshot
ExecStart=/home/$USER/go/bin/discord-rss-cron
IOSchedulingClass=idle
^D

% systemctl enable --user discord-rss-cron.timer
Created symlink ~/.config/systemd/user/timers.target.wants/discord-rss-cron.service → ~/.config/systemd/user/discord-rss-cron.service.

```

Ensure you update $USER with your own user name

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

## Run and test in systemd

```
% systemd-run --user discord-rss-cron        
Running as unit: run-r4acdd1003f7d4de4be2be64ae1be9ad4.service; invocation ID: fa4917cee43c463e9317c32548291cf1
% journalctl -xe --user
```
You should see:
```
Jun 30 11:49:57 discord-rss-cron[287802]: 2024/06/30 11:49:57 Config: /home/$USER/.config/discord-rss-webhook/channel-1
Jun 30 11:49:57 discord-rss-cron[287802]: 2024/06/30 11:49:57 Done nothing found
```