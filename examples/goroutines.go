package main

import (
  ".."
  "strconv"
  "runtime"
  "fmt"
  "bytes"
)

// quick test for concurrent log messages
func main() {
  fmt.Printf("%d -> %d\n", runtime.NumCPU, runtime.GOMAXPROCS(1000))

  buf := bytes.NewBuffer([]byte(""))
  logger := scrolls.NewLogger(buf)
  chans := make([]chan bool, 1000)

  for i := range chans {
    chans[i] = make(chan bool)
    go LogLikeCrazy(strconv.Itoa(i), logger, chans[i])
  }

  for _, ch := range chans {
    <- ch
  }
  fmt.Println(buf.String())
}

func LogLikeCrazy(key string, logger scrolls.Logger, stopped chan bool) {
  for i := 0; i < 10; i += 1 {
    logger.Log(scrolls.LogData{"key": key, "i": strconv.Itoa(i+1)})
  }
  stopped <- true
}
