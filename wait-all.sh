#!/bin/sh

./wait-for-it.sh -t 180 db:5432 -- echo "DB is up"
./wait-for-it.sh -t 180 rabbitmq:5672 -- echo "Rabbitmq 5672 is up"
./wait-for-it.sh -t 180 rabbitmq:15672 -- echo "Rabbitmq 15672 is up"
exec "$@"