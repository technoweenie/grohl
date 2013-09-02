# Grohl

Grohl is an opinionated library for outputting logs in a key=value structure.
This event data is used to drive metrics and monitoring services.

Dave Grohl is the lead singer of Foo Fighters.  I hear he's also passionate about
event driven metrics.

This is a Go version of [asenchi/scrolls](https://github.com/asenchi/scrolls).
The rest of this README is a direct rip from scrolls, with the Ruby snippets
replaced with Go.

## Usage

At Heroku we are big believers in "logs as data". We log everything so
that we can act upon that event stream of logs. Internally we use logs
to produce metrics and monitoring data that we can alert on.

Here's an example of a log you might specify in your application:

```go
grohl.Log(grohl.LogData{"fn": "trap", "signal": s, "at": "exit", "status": 0})
```

The output of which might be:

    fn=trap signal=TERM at=exit status=0

This provides a rich set of data that we can parse and act upon.

A feature of Grohl is setting contexts. Grohl has two types of
context. One is 'global_context' that prepends every log in your
application with that data and a local 'context' which can be used,
for example, to wrap requests with a request id.

In our example above, the log message is rather generic, so in order
to provide more context we might set a global context that links this
log data to our application and deployment:

```go
grohl.AddContext("app", "myapp")
grohl.AddContext("deploy", os.Getenv("DEPLOY"))
```

This would change our log output above to:

    app=myapp deploy=production fn=trap signal=TERM at=exit status=0

If we were in a file and wanted to wrap a particular point of context
we might also do something similar to:

```go
context := grohl.NewContext(grohl.LogData{"ns": "server"})
context.Log(grohl.LogData{"fn": "trap", "signal": s, "at": "exit", "status": 0})
```

This would be the output (taking into consideration our global context
above):

    app=myapp deploy=production ns=server fn=trap signal=TERM at=exit status=0

This allows us to track this log to `Server#trap` and we received a
'TERM' signal and exited 0.

As you can see we have some standard nomenclature around logging.
Here's a cheat sheet for some of the methods we use:

* `app`: Application
* `lib`: Library
* `ns`: Namespace (Class, Module or files)
* `fn`: Function
* `at`: Execution point
* `deploy`: Our deployment (typically an environment variable i.e. `DEPLOY=staging`)
* `elapsed`: Measurements (Time)
* `count`: Measurements (Counters)

Grohl makes it easy to measure the run time of a portion of code.
For example:

```go
timer := grohl.NewTimer(grohl.LogData{"fn": "test"})
grohl.Log(grohl.LogData{"status": "exec"})
// code here
timer.Log(nil)
```

This will output the following log:

    fn=test at=start
    status=exec
    fn=test at=finish elapsed=0.300

You can change the time unit that Grohl uses to "milliseconds" (the
default is "seconds"):

```go
grohl.SetTimeUnit("ms")
```

If you need multiple loggers with different contexts or time units, you can
create them instead of going through the `grohl` functions.  The functions work
identically on grohl loggers with the exception of `SetTimeUnit()`.  You can
access the TimeUnit field manually.

```go
buf := bytes.NewBuffer([]byte(""))
logger := grohl.NewLogger(buf)
logger.Log(grohl.LogData{"fn": "trap"})
```

Grohl has a rich #parse method to handle a number of cases. Here is
a look at some of the ways Grohl handles certain values.

Time and nil:

```go
grohl.Log("t": time.Date(2012, 6, 19, 11, 2, 47, 0, time.UTC), "this": nil)

t=2012-06-19T11:02:47-0400 this=nil
```

True/False:

```go
grohl.Log("that": false, "this": true)

that=false this=true
```
