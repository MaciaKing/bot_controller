
sudo apt install golang-go -y
sudo apt install git -y

cd /opt

git clone https://github.com/MaciaKing/bot_controller.git
cd bot_controller
go build -o bot
chmod +x bot

# Install the bot service
echo "[Unit]
Description=Bot service

[Service]
ExecStart=/opt/bot_controller/bot -ip 192.168.0.11

[Install]
WantedBy=multi-user.target
" > /etc/systemd/system/bot.service

# Start the new service
sudo systemctl daemon-reload
sudo systemctl start bot.service
sudo systemctl enable bot.service
