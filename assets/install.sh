sudo cp ./unread-bot.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable --now unread-bot
sudo systemctl status unread-bot
