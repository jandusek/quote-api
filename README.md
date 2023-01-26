# quote-api

REST API for fetching a random quote from MongoDB

## Example usage

```sh
$ curl http://server/get_quote_ascii
{"_id":"63a3063b51d143b33e91af10","quote":"The fact that some geniuses were laughed at does not imply that all who are laughed at are geniuses. They laughed at Columbus, they laughed at Fulton, they laughed at the Wright brothers. But they also laughed at Bozo the Clown.","author":"Carl Sagan"}
```

## Installation

1. `go build`
2. `docker build . -t koiru/quotes-api`
3. `docker push koiru/quotes-api:latest`

Available on [Docker Hub](https://hub.docker.com/r/koiru/quotes-api).
