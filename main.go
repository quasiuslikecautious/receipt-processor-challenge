package main

import (
  "quasiuslikecautious/receipt-processor-challenge/app"
  "quasiuslikecautious/receipt-processor-challenge/config"
)

func main() {
  config := config.DefaultConfig()

  app := &app.App{}
  app.Initialize(config)
  app.Run()
}
