FROM redis:7.2-alpine

COPY ./entrypoint.sh /custom-entrypoint.sh
RUN chmod +x /custom-entrypoint.sh

RUN mkdir -p /usr/local/etc/redis/

ENTRYPOINT [ "/custom-entrypoint.sh" ]
CMD ["redis-server", "/usr/local/etc/redis/redis.conf"]