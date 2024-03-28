#!/bin/sh

echo "requirepass $REDIS_PASSWORD" > /usr/local/etc/redis/redis.conf

docker-entrypoint.sh "$@"