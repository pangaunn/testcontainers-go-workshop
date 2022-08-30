#!/bin/bash

set -e

ES=http://localhost:9200
until $(curl --output /dev/null --silent --head --fail $ES); do
    echo 'waiting for elasticsearch...'
    sleep 5
done

cd /pre-test-script

curl -X POST -u elastic:changeme -H 'Content-Type: application/json' $ES/_template/books_temp?include_type_name=true --data-binary @books_temp.json
curl -X DELETE -u elastic:changeme -H 'Content-Type: application/json' $ES/books_temp
curl -X PUT -u elastic:changeme -H 'Content-Type: application/json' $ES/books_temp
curl -X POST -u elastic:changeme -H 'Content-Type: application/json' $ES/books_temp/_alias/books

curl -X POST -u elastic:changeme $ES/books/_bulk?pretty -H "Content-Type: application/json" --data-binary @data.json
