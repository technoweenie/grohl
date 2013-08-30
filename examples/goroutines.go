package main

import (
  ".."
  "strconv"
)

// quick test for concurrent log messages
func main() {
  logger := scrolls.NewLogger(nil)
  aStopped := make(chan bool)
  bStopped := make(chan bool)
  cStopped := make(chan bool)
  dStopped := make(chan bool)
  eStopped := make(chan bool)

  go LogLikeCrazy("a", logger, aStopped)
  go LogLikeCrazy("b", logger, bStopped)
  go LogLikeCrazy("c", logger, cStopped)
  go LogLikeCrazy("d", logger, dStopped)
  go LogLikeCrazy("e", logger, eStopped)

  <- aStopped
  <- bStopped
  <- cStopped
  <- dStopped
  <- eStopped
}

func LogLikeCrazy(key string, logger scrolls.Logger, stopped chan bool) {
  for i := 0; i < 1000; i += 1 {
    logger.Log(scrolls.LogData{"key": key, "i": strconv.Itoa(i+1)})
  }
  stopped <- true
}
