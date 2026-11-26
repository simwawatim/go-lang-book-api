#!/bin/sh
set -e

# Run migrations before starting (like `python manage.py migrate`).
# Set RUN_MIGRATIONS=false if you run them as a separate k8s Job instead.
if [ "${RUN_MIGRATIONS:-true}" = "true" ]; then
  ./migrate up
fi

exec ./server
