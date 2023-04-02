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

To call the server

```sh
$ curl -u "user_token:" http://localhost:8080/tickers/PFE/history
$ curl -u "user_token:" http://localhost:8080/tickers
```

user_token can be whatever.