# netdb

[![pkg.go.dev](https://img.shields.io/badge/-reference-blue?logo=go&logoColor=white&labelColor=505050)](https://pkg.go.dev/github.com/thediveo/netdb)
[![GitHub](https://img.shields.io/github/license/thediveo/netdb)](https://img.shields.io/github/license/thediveo/netdb)
![build and test](https://github.com/thediveo/netdb/actions/workflows/buildandtest.yaml/badge.svg?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/thediveo/lxkns)](https://goreportcard.com/report/github.com/thediveo/netdb)
![Coverage](https://img.shields.io/badge/Coverage-97.2%25-brightgreen)

`netdb` provides information about TCP/IP subsystem protocols and internet
services, all this in (pure) Go. By default, it uses its built-in database
instead of consulting `/etc/protocols`, `/etc/services`, and `/etc/ethertypes`.
If needed, it can also consult these files, please see the examples in the
[documentation](https://pkg.go.dev/github.com/thediveo/netdb).

The built-in database has been auto-generated from the `etc/protocols`,
`etc/ethertypes`, and `etc/services` files courtesy of the
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

## DevContainer

> [!CAUTION]
>
> Do **not** use VSCode's "~~Dev Containers: Clone Repository in Container
> Volume~~" command, as it is utterly broken by design, ignoring
> `.devcontainer/devcontainer.json`.

1. `git clone https://github.com/thediveo/netdb`
2. in VSCode: Ctrl+Shift+P, "Dev Containers: Open Workspace in Container..."
3. select `netdb.code-workspace` and off you go...

## Supported Go Versions

`netdb` supports versions of Go that are noted by the [Go release
policy](https://golang.org/doc/devel/release.html#policy), that is, major
versions _N_ and _N_-1 (where _N_ is the current major version).

## Contributing

Please see [CONTRIBUTING.md](CONTRIBUTING.md).

## Copyright and License

`netdb` is Copyright 2021-26 Harald Albrecht, and licensed under the Apache License,
Version 2.0.
