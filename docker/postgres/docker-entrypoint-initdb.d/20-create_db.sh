#!/usr/bin/env sh

set -e

echo "==> Creating databases…"

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<-EOSQL
  CREATE DATABASE alf_data_development;
  CREATE DATABASE alf_store_development;
EOSQL
