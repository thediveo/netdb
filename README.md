# netdb

`netdb` provides information about TCP/IP subsystem protocols and internet
services, all this in (pure) Go. By default, it uses its built-in instead of
consulting `/etc/protocols` and `/etc/services`. If needed, it can also consult
these files, please see the examples in the documentation.

The built-in database has been auto-generated from the etc/protocols and
etc/services files courtesy of the [netdb](https://salsa.debian.org/md/netbase)
package of the Debian project.

Our `netdb` package does not even try to slavishly replicate the POSIX C API;
instead, it tries to be Go-ish. For instance, the C type "servent" has simply
become the netdb.Service type in order to avoid arcane type names.

## Acknowledgement

In some sense, this `netdb` package picks up the baton from the
[@dominikh/go-netdb](https://pkg.go.dev/honnef.co/go/netdb) package. However, it
is not a fork, but was written from scratch, considering (at least some of) the
advice in issue #1 of the go-netdb package.

## ⚖️ Copyright and License

`netdb` is Copyright 2021 Harald Albrecht, and licensed under the Apache License,
Version 2.0.
