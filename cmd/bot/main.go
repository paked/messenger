package main

import (
	"fmt"
	"net/http"

	"github.com/paked/configure"
	"github.com/paked/messenger"
)

var (
	conf        = configure.New()
	verifyToken = conf.String("verify-token", "mad-skrilla", "The token used to verify facebook")
	verify      = conf.Bool("should-verify", false, "Whether or not the app should verify itself")
)

func main() {
	conf.Use(configure.NewFlag())
	conf.Use(configure.NewEnvironment())
	conf.Use(configure.NewJSONFromFile("config.json"))

	conf.Parse()

	m := messenger.New(messenger.MessengerOptions{
		Verify:      *verify,
		VerifyToken: *verifyToken,
	})

	fmt.Println("Serving messenger bot on localhost:8080")
	http.ListenAndServe("localhost:8080", m.Handler())
}
