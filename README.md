# go-shortener [![Build Status](https://travis-ci.org/guilhermef/go-shortener.svg?branch=master)](https://travis-ci.org/guilhermef/go-shortener)
A simple redis based url shortener

## Installation

Get the package with
```
go get github.com/guilhermef/go-shortener
```

then run `godep restore` to fetch necessary dependencies.

## Usage

Run with `make run`. Set parameters with environment variables or a __settings.yml__, as described in the next section.

__go-shortener__ reads the origin and target URIs from a REDIS server. REDIS keys correspond to origin URLs with a `go-shortener:` prefix and the values will be the target URIs. So, for instance, a redirect from _/origin_ to _/target_ will require the following entry in REDIS:
```
go-shortener:/origin  ->  /target
```

This can be set in REDIS via the command `SET go-shortener:/origin /target`.

A counter is incremented whenever its corresponding key is hit. Counter key names have the `go-shortener-count:` prefix so, considering the previous example, if __15__ requests are made to `/origin`, REDIS will have the following entry:
```
go-shortener-count:/origin  ->  15
```

## Settings
Following are the necessary settings and their default values:

* `PORT`: __go-shortener__ listening port. Default: _1234_
* `LOG_PATH`: Path to file where log entries are written. Default: _STDOUT_
* `REDIS_HOST`: Address to Redis server. Default: _127.0.0.1:6379_
* `REDIS_PASS`: Password to Redis server. Default: ""
* `REDIS_DB`: Redis database to which __go-shortener__ should connect. Default: 0

They can be set via environment variables (following the above nomenclature) or through a _settings.yml_ file in the following form

```YAML
logpath: path_to_log_file
port: 1234
redishost: 127.0.0.1:6379
redispass: pass
redisdb: 0
```

## Custom Redirect

If the key on Redis database isn't found, that's possible to make a custom redirect instead of a simple response, like Not Found.
To do this, set the following environment variable

* `REDIRECT_HOST`: Host to redirect (Required, if wanna make a custom redirect)
* `REDIRECT_CODE`: Response status code for that redirect. Default: 302

## Test

Test by running `make test`.
