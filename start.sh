#!/bin/sh

set -e

echo "run db migration on $DBSOURCE"
source /app/app.env
/app/migrate -path /app/migration -database "$DBSOURCE" -verbose up

echo "start the app"
exec "$@"