#!/bin/bash

for id in {1..500}
do
   echo "Create book-$id"
   curl --request POST 'http://arsenal.default.svc.cluster.local:8080/v1/dummy-books' \
   --header 'Content-Type: application/json' \
   --data-raw "{
       \"title\": \"book-$id\",
       \"author\": \"author-$id\",
       \"pages\": 100
   }"

   echo "Read book-$id"
   curl --request GET "http://arsenal.default.svc.cluster.local:8080/v1/dummy-books/book-$id"

   echo "Read with api delay, book-$id"
   curl --request GET "http://arsenal.default.svc.cluster.local:8080/v1/dummy-books/book-$id?api-delay=true"

   echo "Read with db delay, book-$id"
   curl --request GET "http://arsenal.default.svc.cluster.local:8080/v1/dummy-books/book-$id?db-delay=true"

   echo "Read with error and api delay, book-$id"
   curl --request GET "http://arsenal.default.svc.cluster.local:8080/v1/dummy-books/book-$id?error=true&api-delay=true"

   sleep 60
done





