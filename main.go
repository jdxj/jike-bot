package main

import "fmt"

func main() {
	fmt.Printf("%+v\n", conf)
	bot := New()
	bot.Run()
}
