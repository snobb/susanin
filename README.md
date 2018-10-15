Susanin - Go HTTP router
========================

`Susanin` is a lightweight HTTP router for building Go HTTP services. It is
built on the new `context` package introduced in Go 1.7 to handle
request-scoped values across a handler chain, etc.

The router is inspired by [http-hash](https://www.npmjs.com/package/http-hash) module in nodejs.

The focus of the project has been to seek out an elegant and comfortable design for writing
general purpose REST API servers.

## Install

`go get -u github.com/snobb/susanin/pkg/framework`

Optionally there are some middleware available:

`go get -u github.com/snobb/susanin/pkg/middleware`


## Features

* **Lightweight** - tiny in size ~400SLOC.
* **100% compatible with net/http** - use any http or middleware pkg in the ecosystem that is also compatible with `net/http`
* **Context control** - built on new `context` package
* **No external dependencies** - plain ol' Go 1.7+ stdlib + net/http


## Examples

* `examples/server.go` - REST APIs made easy, productive and maintainable
