#!/bin/sh

# Creates the pass file for the servers.json file
echo "$PGPASS" > /etc/pgpass

/entrypoint.sh "$@"