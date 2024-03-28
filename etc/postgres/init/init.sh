#!/bin/bash

set -e


function provision_svc_if_not_exists() {
	TOADD="$1"

	USERS_FOUND=$(psql -U $POSTGRES_USER -tc "SELECT 1 FROM pg_catalog.pg_roles WHERE  rolname = '$TOADD'" | wc -l)

	if [ "$USERS_FOUND" == "2" ]; then
		echo "Skipping creation for $TOADD, as they already exist"
		return 0
	fi

	psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<-EOSQL
		CREATE USER $TOADD WITH PASSWORD '$TOADD-pass';
		CREATE SCHEMA IF NOT EXISTS $TOADD AUTHORIZATION $TOADD;
		GRANT ALL PRIVILEGES ON SCHEMA $TOADD TO $TOADD;
		GRANT USAGE ON ALL SEQUENCES IN SCHEMA $TOADD TO $TOADD;

		ALTER DEFAULT PRIVILEGES IN SCHEMA $TOADD GRANT ALL PRIVILEGES ON TABLES TO $TOADD;
		ALTER DEFAULT PRIVILEGES IN SCHEMA $TOADD GRANT USAGE          ON SEQUENCES TO $TOADD;
	EOSQL
}

function main() {
	while IFS="" read -r LINE || [ -n "$LINE" ]
	do
		provision_svc_if_not_exists "$LINE"
	done < $PROVISIONING_FILE
}

main "$@"

