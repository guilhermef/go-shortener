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

## Test

Test by running `make test`.