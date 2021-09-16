/*

Package netdb provides information about TCP/IP subsystem protocols and internet
services, as commonly stored in files /etc/protocols and /etc/services.

This package is a pure Go implementation and defaults to looking up protocols
and services information in its built-in database instead of consulting
/etc/protocols and /etc/services. Please see the examples for how to also take
the well-known /etc/protocols and /etc/services files into consideration.

The netdb package does not even try to slavishly replicate the POSIX C API;
instead, it tries to be Go-ish. For instance, the C type "servent" has simply
become the netdb.Service type in order to avoid arcane type names.

Notes

This package bases on the file format descriptions for protocols(5) and
services(5), as documented in
https://man7.org/linux/man-pages/man5/protocols.5.html and
https://man7.org/linux/man-pages/man5/services.5.html.

The built-in database has been auto-generated from the etc/protocols and
etc/services files courtesy of the netdb package of the Debian project.

In some sense, this netdb package picks up the baton from the
https://pkg.go.dev/honnef.co/go/netdb package. However, it is not a fork, but
was written from scratch, considering (at least some of) the advice in issue #1
of the go-netdb package.

*/
package netdb
