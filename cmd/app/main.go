package main

import (
	a "an3softbot/internal/app"
)

var App a.Application

func main() {
	App = a.Application{}
	App.Run()
}
