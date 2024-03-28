FROM dpage/pgadmin4:8.4

# Cause it's locked down annoyingly for a dev environment
# I can't create the pgpass file in the entrypoint
USER root

COPY ./config/servers.json /pgadmin4/servers.json
COPY ./custom-entrypoint.sh /custom-entrypoint.sh
RUN chmod +x /custom-entrypoint.sh

ENTRYPOINT [ "/custom-entrypoint.sh" ]