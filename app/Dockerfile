FROM golang:1.24-bullseye as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o bin/server ./cmd/server

FROM chromedp/headless-shell:latest

RUN apt-get update && apt-get install -y \
    fonts-liberation \
    libjpeg62-turbo \
    libpng-dev \
    libx11-xcb1 \
    libxtst6 \
    libharfbuzz0b \
    libpango-1.0-0 \
    libcairo2 \
    cron \
    coreutils \
    procps \
    && apt-get clean

USER root
RUN apt-get update && apt-get install -y tini && apt-get clean

WORKDIR /app

COPY --from=builder /app/bin/server /server
COPY --from=builder /app /app

RUN mkdir -p /cron /scripts /app/screens /var/log
RUN chmod -R 777 /app/screens
RUN touch /var/log/cron.log && chmod 666 /var/log/cron.log

COPY cron/cleanup.sh /cron/cleanup.sh
RUN chmod +x /cron/cleanup.sh

COPY scripts/start.sh /scripts/start.sh
RUN chmod +x /scripts/start.sh

RUN echo "0 * * * * /cron/cleanup.sh" | crontab -

ENTRYPOINT ["/usr/bin/tini", "--"]
CMD ["/scripts/start.sh"]