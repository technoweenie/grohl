package main

import (
  ".."
  "strconv"
  "runtime"
  "fmt"
)

// quick test for concurrent log messages
func main() {
  fmt.Printf("%d -> %d\n", runtime.NumCPU, runtime.GOMAXPROCS(1000))

  logger := grohl.NewLogger(nil)
  chans := make([]chan bool, 1000)

  for i := range chans {
    chans[i] = make(chan bool)
    go LogLikeCrazy(strconv.Itoa(i), logger, chans[i])
  }

  for _, ch := range chans {
    <- ch
  }
  timer := grohl.NewTimer(grohl.LogData{"fn": "test"})
  grohl.Log(grohl.LogData{"status": "exec"})
  // code here
  timer.Log(nil)
}

func LogLikeCrazy(key string, logger grohl.Logger, stopped chan bool) {
  for i := 0; i < 10; i += 1 {
    logger.Log(grohl.LogData{"key": key, "i": strconv.Itoa(i+1)})
  }
  stopped <- true
}
