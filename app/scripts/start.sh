#!/bin/bash
echo "Run start.sh"
mkdir -p /app/screens
mkdir -p /var/log
touch /var/log/cron.log
chmod 666 /var/log/cron.log
/usr/sbin/cron -f -L 15 &
sleep 2
if pgrep cron > /dev/null; then
    echo "Cron is running"
else
    echo "Cron is not running"
    exit 1
fi
#tail -f /var/log/cron.log &
echo "Application running..."
/server