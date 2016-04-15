# Messenger [![GoDoc](https://godoc.org/github.com/paked/messenger?status.svg)](https://godoc.org/github.com/paked/messenger)

This is a Go library for making bots to be used on Facebook messenger. It is built on the [Messenger Platform](https://developers.facebook.com/docs/messenger-platform). One of the main goals of the project is to implement it in an idiomatic and easy to use fashion.

[You can find an example of how to use it here](https://github.com/paked/messenger/blob/master/cmd/bot/main.go)

## Tips

- Follow the [quickstart](https://developers.facebook.com/docs/messenger-platform/quickstart) guide for getting everything set up!
- You need a Facebook development app, and a Facebook page in order to build things.
- Use [ngrok](https://ngrok.com) to tunnel your locally runnning bot so that Facebook can reach the webhook.

## Inspiration

Messenger takes design cues from:

- [`net/http`](https://godoc.org/net/http)
- [`github.com/nickvanw/ircx`](https://github.com/nickvanw/ircx)
