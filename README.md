# go-shortener [![Build Status](https://travis-ci.org/guilhermef/go-shortener.svg?branch=master)](https://travis-ci.org/guilhermef/go-shortener)
A simple redis based url shortener

## Settings
Settings are read from a _settings.yml_ file, which should follow the structure below

```YAML
logpath: path_to_log_file
port: 1234
redishost: localhost:6379
```

`port` and `redishost` have default values of `1234` and `localhost:6379`, respectively. If `logpath` is not informed, log entries will be sent to STDOUT.