package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/paked/configure"
	"github.com/paked/messenger"
)

var (
	conf        = configure.New()
	verifyToken = conf.String("verify-token", "mad-skrilla", "The token used to verify facebook")
	verify      = conf.Bool("should-verify", false, "Whether or not the app should verify itself")
	pageToken   = conf.String("page-token", "not skrilla", "The token that is used to verify the page on facebook")
)

func main() {
	conf.Use(configure.NewFlag())
	conf.Use(configure.NewEnvironment())
	conf.Use(configure.NewJSONFromFile("config.json"))

	conf.Parse()

	m := messenger.New(messenger.MessengerOptions{
		Verify:      *verify,
		VerifyToken: *verifyToken,
		Token:       *pageToken,
	})

	m.HandleMessage(func(m messenger.Message, r *messenger.Response) {
		fmt.Printf("%v (Sent, %v)\n", m.Text, m.Time.Format(time.UnixDate))
		fmt.Println(r.Text("Hello, World!"))
	})

	m.HandleDelivery(func(d messenger.Delivery, r *messenger.Response) {
		fmt.Println(d.Watermark().Format(time.UnixDate))
	})

	fmt.Println("Serving messenger bot on localhost:8080")

	http.ListenAndServe("localhost:8080", m.Handler())
}
