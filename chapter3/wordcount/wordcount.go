package main

import (
  "chapter3/words"
  "fmt"
  "os"
)


func main() {
  text := os.Args[1]
  count := words.CountWords(text)
  fmt.Printf("一共: %d\n", count)
}
