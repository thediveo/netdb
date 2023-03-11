# netdb

[![pkg.go.dev](https://img.shields.io/badge/-reference-blue?logo=go&logoColor=white&labelColor=505050)](https://pkg.go.dev/github.com/thediveo/netdb)
[![GitHub](https://img.shields.io/github/license/thediveo/lxkns)](https://img.shields.io/github/license/thediveo/netdb)
![build and test](https://github.com/thediveo/netdb/workflows/build%20and%20test/badge.svg?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/thediveo/lxkns)](https://goreportcard.com/report/github.com/thediveo/netdb)
![Coverage](https://img.shields.io/badge/Coverage-97.9%25-brightgreen)

`netdb` provides information about TCP/IP subsystem protocols and internet
services, all this in (pure) Go. By default, it uses its built-in database
instead of consulting `/etc/protocols` and `/etc/services`. If needed, it can
also consult these files, please see the examples in the
[documentation](https://pkg.go.dev/github.com/thediveo/netdb).

The built-in database has been auto-generated from the `etc/protocols` and
`etc/services` files courtesy of the
[netbase](https://salsa.debian.org/md/netbase) package of the Debian project.

This `netdb` package does not even try to slavishly replicate the POSIX C API;
instead, it attempts to be Go-ish. For instance, the C type `servent` has simply
become the `netdb.Service` type in order to avoid arcane POSIX-rooted type
names.

Please refer to the [reference
documentation](https://pkg.go.dev/github.com/thediveo/netdb) for usage examples.

## Acknowledgement

In some sense, this `netdb` package picks up the baton from the
[@dominikh/go-netdb](https://github.com/dominikh/go-netdb) package. However, it
is not a fork but was written from scratch, considering (at least some of) the
advice in [issue #1](https://github.com/dominikh/go-netdb/issues/1) of the
go-netdb package.

## Copyright and License

`netdb` is Copyright 2021-23 Harald Albrecht, and licensed under the Apache License,
Version 2.0.
