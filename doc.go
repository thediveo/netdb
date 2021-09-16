/*

Package netdb provides information about TCP/IP subsystem protocols and internet
services, as commonly stored in files /etc/protocols and /etc/services. It is a
pure Go implementation.

This netdb package does not even try to slavishly replicate the POSIX C API;
instead, it attempts to be Go-ish. For instance, the C type "servent" has simply
become the netdb.Service type in order to avoid the arcane POSIX-rooted type
names.

The netdb package defaults to looking up protocols and services information from
its built-in database instead of consulting /etc/protocols and /etc/services.
Additionally, it also supports reading the protocol and service descriptions
from the well-known /etc/protocols and /etc/services files; please see the
examples for how to access these sources.

Notes

This package bases on the file format descriptions for protocols(5) and
services(5), as documented in
https://man7.org/linux/man-pages/man5/protocols.5.html and
https://man7.org/linux/man-pages/man5/services.5.html.

The built-in database has been auto-generated from the etc/protocols and
etc/services files courtesy of the netbase package of the Debian project
(https://salsa.debian.org/md/netbase).

In some sense, this netdb package picks up the baton from the
https://github.com/dominikh/go-netdb package. However, it is not a fork but was
written from scratch, considering (at least some of) the advice in issue #1 of
the go-netdb package.

*/
package netdb
