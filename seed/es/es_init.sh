#!/bin/bash

set -e

ES=http://elasticsearch:9200
until $(curl --output /dev/null --silent --head --fail $ES); do
    echo 'waiting for elasticsearch...'
    sleep 5
done


curl -X POST -u elastic:changeme -H 'Content-Type: application/json' $ES/_template/books_temp?include_type_name=true --data-binary @/usr/local/bin/books_temp.json
curl -X DELETE -u elastic:changeme -H 'Content-Type: application/json' $ES/books_temp
curl -X PUT -u elastic:changeme -H 'Content-Type: application/json' $ES/books_temp
curl -X POST -u elastic:changeme -H 'Content-Type: application/json' $ES/books_temp/_alias/books

curl -XPOST http://elasticsearch:9200/books/_bulk?pretty -H "Content-Type: application/json" --data-binary @/usr/local/bin/data.json
