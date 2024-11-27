# own-redis

## Mandatory Part

Writing a key-value store via REST API would be too easy, wouldn't it? Let's make it a bit more complicated and let the client and your application communicate using the UDP protocol, i.e. each request and response is a single UDP packet. In our key-value store implementation you have to implement three methods SET, GET and PING. SET puts a key-value and GET gets and returns the given value back to the client. The PING command verifies that the storage is working.

NOTE:
- Command names, command arguments are  case-insensitive. So `PING`, `ping` and `Ping` are all valid and denote the same command.

### PING

`PING` is one of the simplest Redis commands. It's used to check whether a Redis server is healthy. The response for the `PING` command is `PONG`.
```sh
$ nc 0.0.0.0 8080
PING
PONG
```

### SET

Any request that has the string SET as the first argument in its message will be considered an insert request.

Example:
- `SET foo bar` will insert a key `foo` with value “`bar`”.
- `SET foo bar baz` will insert a key `foo` with value “`bar baz`”.

If the number of arguments is not enough to save the key, the server should return an error.

Example:
- `SET KEYVAL` will return error with text “`(error) ERR wrong number of arguments for 'SET' command`”
- `SET` will return error with text “`(error) ERR wrong number of arguments for 'SET' command`”.

SET should return `OK`.

```sh
$ nc 0.0.0.0 8080
SET Foo Bar
OK
```
#### Options
The `SET` command supports option `PX` that modify its behavior:
- `PX` _milliseconds_ - Set the specified expire time, in milliseconds (a positive integer).

Example:
```sh
$ nc 0.0.0.0 8080
SET foo bar px 10000
OK
GET foo
bar
```
A request within 10000 milliseconds will produce a bar response, but once the time is up, the server should clear the value and should return `(nil)` when client attempting to retrieve the value.
```sh
$ nc 0.0.0.0 8080
GET foo
(nil)
```


### GET

A GET request is any request in which the first argument contains the `GET` command. When attempting to query an existing key, the server must return the previously stored value in response.

Example:
```sh
$ nc 0.0.0.0 8080
SET Foo Bar
OK
GET Foo
Bar
```


Example:
```sh
$ nc 0.0.0.0 8080
GET RandomKey
(nil)
```

If the client tries to insert a value into an existing key, your application should update the value to the last value specified by the client.

Example:
```sh
$ nc 0.0.0.0 8080
SET Foo Bar
OK
GET Foo
Bar
SET Foo Buz
OK
GET Foo
Buz
```

### Usage
Outcomes:

- Program prints usage text.

```shell
$ ./own-redis --help
Own Redis

Usage:
  own-redis [--port <N>]
  own-redis --help

Options:
  --help       Show this screen.
  --port N     Port number.
  --env S      Environment variable ('local', 'dev', 'prod').
```