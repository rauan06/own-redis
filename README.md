# own-redis

## Abstract

A custom key-value database. Such databases are essential for efficiently storing and retrieving data using unique keys. This application will operate in RAM for fast access, utilizing the UDP protocol to facilitate communication between the storage and the client. Key operations will include SET, GET, and PING.

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