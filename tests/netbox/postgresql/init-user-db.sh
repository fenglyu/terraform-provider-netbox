#!/bin/bash
set -e

#psql --username "$POSTGRES_USER" -tc "SELECT 1 FROM pg_database WHERE datname = 'netbox'" | grep -q 1 || psql -U postgres -c "CREATE DATABASE my_db"

psql -v ON_ERROR_STOP=0 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE USER docker;
    CREATE DATABASE netbox;
    CREATE USER netbox WITH PASSWORD 'J5brHrAXFLQSif0K';
    GRANT ALL PRIVILEGES ON DATABASE netbox TO netbox;
EOSQL