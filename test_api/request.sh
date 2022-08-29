curl --location --request GET 'http://localhost:3000/api/v1/book/1'

curl --location --request POST 'http://localhost:3000/api/v1/book/' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "harry potter 2",
    "price": 600,
    "author": "JK rowling",
    "description": "harry potter 2",
    "imageUrl": "http://www.adviceforyou.co.th/blog/wp-content/uploads/2011/12/harry-potter.jpeg"
}'

curl --location --request DELETE 'http://localhost:3000/api/v1/book/1'

curl --location --request PUT 'http://localhost:3000/api/v1/book/1' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "harry potter 1.5",
    "price": 350,
    "author": "JK rowling",
    "description": "harry porter 1.5",
    "imageUrl": "http://www.adviceforyou.co.th/blog/wp-content/uploads/2011/12/harry-potter.jpeg"
}'