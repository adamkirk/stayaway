FROM golang:1.22 as dev

RUN go install github.com/air-verse/air@latest

COPY ./.local/docker/watch.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENV PATH="$PATH:/app/build"

# This is where you should mount your code
VOLUME /app

ENTRYPOINT ["/entrypoint.sh"]
