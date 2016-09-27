# Messenger [![GoDoc](https://godoc.org/github.com/paked/messenger?status.svg)](https://godoc.org/github.com/paked/messenger)

This is a Go library for making bots to be used on Facebook messenger. It is built on the [Messenger Platform](https://developers.facebook.com/docs/messenger-platform). One of the main goals of the project is to implement it in an idiomatic and easy to use fashion.

[You can find an example of how to use it here](https://github.com/paked/messenger/blob/master/cmd/bot/main.go)

## Tips

- Follow the [quickstart](https://developers.facebook.com/docs/messenger-platform/quickstart) guide for getting everything set up!
- You need a Facebook development app, and a Facebook page in order to build things.
- Use [ngrok](https://ngrok.com) to tunnel your locally runnning bot so that Facebook can reach the webhook.

## Breaking Changes

`paked/messenger` is a pretty stable library however, changes will be made which might break backwards compatibility. For the convenience of its users, these are documented here.


- [12/9/16](https://github.com/paked/messenger/commit/47f193fc858e2d710c061e88b12dbd804a399e57): Removing unused parameter `text string` from function `(r *Response) GenericTemplate`.
- [20/5/16](https://github.com/paked/messenger/commit/1dc4bcc67dec50e2f58436ffbc7d61ca9da5b943): Leaving the `WebhookURL` field blank in `Options` will yield a URL of "/" instead of a panic.
- [4/5/16](https://github.com/paked/messenger/commit/eb0e72a5dcd3bfaffcfe88dced6d6ac5247f9da1): The URL to use for the webhook is changable in the `Options` struct. 

## Inspiration

Messenger takes design cues from:

- [`net/http`](https://godoc.org/net/http)
- [`github.com/nickvanw/ircx`](https://github.com/nickvanw/ircx)

## Example

Awesome project lists using [Messenger](https://github.com/paked/messenger) library.

* [drone-facebook](https://github.com/appleboy/drone-facebook): Drone plugin for sending Facebook notifications
