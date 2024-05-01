package main

import "demo-api/pkg/app"

func main() {
	demo := app.NewApp()
	demo.Run()
}
