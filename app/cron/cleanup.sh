#!/bin/bash

# Delete all screens .jpg older than 60 minutes в /app/screens
find /app/screens -type f -name "*.jpg" -mmin +60 -exec rm {} \;

echo "$(date '+%Y-%m-%d %H:%M:%S') - Old screenshots deleted" >> /var/log/cron.log