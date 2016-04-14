# Messenger

This is a Go library for making bots to be used on Facebook messenger. It is built on the [Messenger Platform](https://developers.facebook.com/docs/messenger-platform). One of the main goals of the project is to implement it in an idiomatic and easy to use fashion.

**It is currently in very early development. [Discussion](https://github.com/paked/messenger/issues/new) is much appreciated!**

## Tips

- You need a Facebook development app, and a Facebook page in order to build things.
- Use [ngrok](https://ngrok.com) to tunnel your locally runnning bot so that Facebook can reach the webhook.

## Inspiration

Messenger takes design cues from:

- [`net/http`](https://godoc.org/net/http)
- [`github.com/nickvanw/ircx`](https://github.com/nickvanw/ircx)
