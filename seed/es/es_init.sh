#!/bin/bash

set -e

ES=http://elasticsearch:9200
until $(curl --output /dev/null --silent --head --fail $ES); do
    echo 'waiting for elasticsearch...'
    sleep 3
done


curl -X POST -u elastic:changeme -H 'Content-Type: application/json' $ES/_template/books_temp?include_type_name=true -d "{\"index_patterns\":[\"books_*\"],\"mappings\":{\"_doc\":{\"_source\":{\"enabled\":true},\"dynamic\":\"false\",\"properties\":{\"id\":{\"type\":\"integer\"},\"name\":{\"type\":\"text\",\"analyzer\":\"rebuilt_thai\"},\"price\":{\"type\":\"integer\"},\"author\":{\"type\":\"text\",\"analyzer\":\"rebuilt_thai\"},\"description\":{\"type\":\"text\",\"analyzer\":\"rebuilt_thai\"},\"imageURL\":{\"type\":\"text\"}}}},\"settings\":{\"number_of_shards\":\"1\",\"auto_expand_replicas\":\"0-all\",\"max_result_window\":\"5000\",\"analysis\":{\"filter\":{\"thai_stop\":{\"type\":\"stop\",\"stopwords\":\"_thai_\"}},\"analyzer\":{\"rebuilt_thai\":{\"tokenizer\":\"thai\",\"filter\":[\"lowercase\",\"asciifolding\",\"decimal_digit\",\"thai_stop\"]}}}}}"
curl -X DELETE -u elastic:changeme -H 'Content-Type: application/json' $ES/books_temp
curl -X PUT -u elastic:changeme -H 'Content-Type: application/json' $ES/books_temp
curl -X POST -u elastic:changeme -H 'Content-Type: application/json' $ES/books_temp/_alias/books