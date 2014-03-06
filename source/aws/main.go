package main

import (
  "flag"
  "fmt"
)

func main() {
  port := flag.Int("port", 0, "port to listen on")
  flag.Parse()
  fmt.Printf("lkjsdf '%d'", *port)
}
