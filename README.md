# Grohl

Grohl is an opinionated library for gathering metrics and data about how your
applications are running in production.  It does this through writing logs
in a key=value structure.  It also provides interfaces for sending exceptions
or metrics to external services.

This is a Go version of [asenchi/scrolls](https://github.com/asenchi/scrolls).
The name for this library came from mashing the words "go" and "scrolls"
together.  Even Dave Grohl, lead singer of Foo Fighters, is passionate about
event driven metrics.

## Installation

    $ go get github.com/technoweenie/grohl

Then import it:

    import "github.com/technoweenie/grohl"

See the [godocs](http://godoc.org/github.com/technoweenie/grohl) for more.
