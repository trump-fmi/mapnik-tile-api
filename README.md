# mapnik-tile-api [![Build Status](https://travis-ci.org/trump-fmi/mapnik-tile-api.svg?branch=master)](https://travis-ci.org/trump-fmi/mapnik-tile-api)

Go tool that parses a provided `renderd.conf` and serves the available tile endpoints with a JSON endpoint.
The JSON response is a list of endpoints, each of them is a object with the string keys `name`, `uri` and `description`.
The values are exactly the same as the ones written in the provided `renderd.conf`.
The endpoint cna be found under `/tileEndpoints` using the specified port.

## Installation
To build this project yourself you will need a [correctly configured](https://golang.org/doc/install#testing) go setup and run `go get -u github.com/trump-fmi/mapnik-tile-api`.

## Flags
Flags can be set by `-flag <value>` or `-flag=value`.
If you want to keep the defaults, you do not have to supply anything. 

| Flag   | Default value       | Explanation                             |
|--------|---------------------|-----------------------------------------|
| `port` | `8081`              | Port where the socket is opened.        |
| `path` | `/etc/renderd.conf` | Path to the `renderd.conf` to be parsed |
