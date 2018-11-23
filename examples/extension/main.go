package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/paked/messenger"
)

var (
	serverURL   = flag.String("serverURL", "", "The server (webview) URL, must be https (required)")
	verifyToken = flag.String("verify-token", "mad-skrilla", "The token used to verify facebook (required)")
	verify      = flag.Bool("should-verify", false, "Whether or not the app should verify itself")
	pageToken   = flag.String("page-token", "not skrilla", "The token that is used to verify the page on facebook")
	appSecret   = flag.String("app-secret", "", "The app secret from the facebook developer portal (required)")
	host        = flag.String("host", "localhost", "The host used to serve the messenger bot")
	port        = flag.Int("port", 8080, "The port used to serve the messenger bot")
)

func main() {
	flag.Parse()

	if *verifyToken == "" || *appSecret == "" || *pageToken == "" {
		fmt.Printf("missing arguments\n\n")
		flag.Usage()

		os.Exit(-1)
	}

	client := messenger.New(messenger.Options{
		Verify:      *verify,
		AppSecret:   *appSecret,
		VerifyToken: *verifyToken,
		Token:       *pageToken,
	})

	err := client.EnableChatExtension(messenger.HomeURL{
		URL:                *serverURL,
		WebviewHeightRatio: "tall",
		WebviewShareButton: "show",
		InTest:             true,
	})
	if err != nil {
		fmt.Println("Failed to EnableChatExtension, err=", err)
	}

	// Setup a handler to be triggered when a message is received
	client.HandleMessage(func(m messenger.Message, r *messenger.Response) {
		fmt.Printf("%v (Sent, %v)\n", m.Text, m.Time.Format(time.UnixDate))

		p, err := client.ProfileByID(m.Sender.ID)
		if err != nil {
			fmt.Println("Something went wrong!", err)
		}

		r.Text(fmt.Sprintf("Hello, %v!", p.FirstName), messenger.ResponseType)
	})

	addr := fmt.Sprintf("%s:%d", *host, *port)
	log.Println("Serving messenger bot on", addr)
	log.Fatal(http.ListenAndServe(addr, client.Handler()))
}
