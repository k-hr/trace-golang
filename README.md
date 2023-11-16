## TRACE APP GOLANG

Sample app for emitting traces.

Emits traces for API calls and associated DB calls.

### Implementation
* gorilla mux router
* go pg orm for connection to postgres
* POST and GET APIs for books
* Simulated delays and errors using query params

### Deployment
* Can be deployed on k8s using the files in templates folder

### APIs supported
* `POST http://localhost:8080/v1/books` with the following request body to create book
```
{
    "title": "book-1",
    "author": "author-1",
    "pages": 100
}
```
* `GET http://localhost:8080/v1/books/book-1` to get the book
* `GET http://localhost:8080/v1/books/book-1?api-delay=true` to get the book with simulated delay of 10s in api workflow
* `GET http://localhost:8080/v1/books/book-1?db-delay=true` to get the book with simulated delay of 10s in db call 
* `GET http://localhost:8080/v1/books/book-1?error=true` to simulate runtime error
