#!/bin/sh

set -e

echo "run db migration on $DBSOURCE"
/app/migrate -path /app/migration -database "$DBSOURCE" -verbose up

echo "start the app"
exec "$@"