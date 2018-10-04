# simpleGoHTTPServer

## Requirements

1. Only import go's standard library
1. Receive a POST request to /hash whose body contains a FORM with a key _password_ and value, respectively.
    1. Immediately return a monotonically increasing number (by 1) of stored keys
    1. Wait 5s, perform a sha512 on the password, base64 encode it
1. Client can sent a GET request to /hash/{ID} where ID represents a previously returned ID. Base64 encoded, hashed
   password is returned.
1. Client can send a GET to /stats endpoint that returns a JSON object representing the total number of POST requests to
   /hash and and the average time spent, in microseconds, servicing that endpoint.
1. Client can send a POST request to /shutdown which issues a graceful shutdown to the server, satisfying existing
   requests and not allowing new requests.

## How to run 
HTTP_LISTEN_PORT must be an environmental variable passed into the system to determine the port by which it runs on:

> HTTP_LISTEN_PORT=3000 go run main.go

or the name of executable if you choose to build it.

## Assumptions
1. The system must be able to bind to the port passed into the program. If the port is <1000, the user must be
   a privileged user. Otherwise, the system will exit.
1. It is fine that there is no persistent data store among restarts of the service. Calls to /hash/{ID} between
   subsequent runs may not yield the same result. An idea to mitigate this would be to use a file, another service or a
   database.
1. For the JSON returned from the call to /stats, the _average_ field will be the floor of the median (e.g. 12.8
   microseconds will return 12).
1. If a call to /hash/{ID} is called with an unknown ID, do not inform the user it does not exist. Just return 202
   without a status request.
   1. If a user makes this call with an ID whose 5s timer has not expired the behavior will be the same.

## HTTP Response Codes
1. /hash
    1. 202 (Status Accepted) if the post call succeeds
    1. 405 (Method not allowed) if not a POST call
    1. 400 (Status Bad Request) if a form with _password_ doesn't exist
1. /hash/{ID}
    1. 202 if the GET call succeeds
        1. 202 is also returned if ID does not exist (yet) in the in-memory data store
    1. 405 (Method Not Allowed) if the call is not a GET
    1. 400 if the ID  is not an integer
1. /stats
    1. 202 (Status Accepted) if the GET call succeeds
    1. 405 (Method not allowed) if not a GET call
    1. 500 (Internal Server Error) if a serialization error occurs
1. /shutdown
    1. 202 (Status Accepted) if the POST call succeeds
    1. 405 (Method not allowed) if not a POST call
1. any other route will return 400 (Status Bad Request)

## Things to do
1. Performance tests using golang's benchmark tooling (pprof/Benchmark)! This is to test the mutex capability of the
   http server. Mutexes were used since it is shared state and is faster than channels. However, typically channels are
   the go-to (hah!) of go.
    1. Investigate the possible use of sync.Map, go channels and mutexes (see above!).
1. Write to a persistent file store to retain information between restarts, potentially.
1. Write tests around graceful shutdown. Unit tests and automation--using docker--to be able to perform regression
   tests around graceful shutdowns. As a work around I put in a sleep in one of the handler calls. Obviously, this is
   not sufficient for production, but in the absence of time it is permissible.
1. Timing attacks against /hash/{ID} and thinking through how to securely maintain an API, if client side facing.