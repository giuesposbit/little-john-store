To run the server on 8080 just clone the repo and run

```sh
$ go build
$ ./little-john-store
```

To test

```sh
$ go test
```

Prerequisite, Go-lang binary

To run with Docker

```sh
$ docker build -t little-john-store .
$ docker run -d -p 8080:8080 little-john-store
```

To call the server

```sh
$ curl -u "user_token:" http://localhost:8080/tickers/PFE/history
$ curl -u "user_token:" http://localhost:8080/tickers
```

user_token can be whatever.
