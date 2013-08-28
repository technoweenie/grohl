# Scrolls

Like asenchi/scrolls, but in go.  Let's rip some ruby code samples from scrolls'
README and convert them to go.

Here's an example of a log you might specify in your application:

```go
import "github.com/technoweenie/go-scrolls"

scrolls.Log(map[string]string{
  "fn": "trap", "signal": s, "at": "exit", "status": "0"
})
```

The output of which might be:

```
fn=trap signal=TERM at=exit status=0
```

In our example above, the log message is rather generic, so in order to provide
more context we might set a global context that links this log data to our
application and deployment:

```go
scrolls.GlobalContext["app"] = "myapp"
scrolls.GlobalContext["deploy"] = os.Getenv("DEPLOY")
```

This would change our log output above to:

```
app=myapp deploy=production fn=trap signal=TERM at=exit status=0
```

If we were in a file and wanted to wrap a particular point of context we might also do something similar to:

```go
context := scrolls.NewContext(map[string]string{"ns": "server"})
context.Log(map[string]string{
    "fn": "trap", "signal": s, "at": "exit", "status": "0"
})
```
This would be the output (taking into consideration our global context above):

```
app=myapp deploy=production ns=server fn=trap signal=TERM at=exit status=0
```

Scrolls makes it easy to measure the run time of a portion of code. For example:

```go
timer := scrolls.NewTimer(map[string]string{"fn": "test"})
scrolls.Log("status": "exec")
// code here
timer.Log()
```

This will output the following log:

```
fn=test at=start
status=exec
fn=test at=finish elapsed=0.300
```

You can change the time unit that Scrolls uses to "milliseconds" (the default is "seconds"):

```go
scrolls.SetTimeUnit("ms")
```
