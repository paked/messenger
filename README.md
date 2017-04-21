# Messenger [![GoDoc](https://godoc.org/github.com/paked/messenger?status.svg)](https://godoc.org/github.com/paked/messenger)

This is a Go library for making bots to be used on Facebook messenger. It is built on the [Messenger Platform](https://developers.facebook.com/docs/messenger-platform). One of the main goals of the project is to implement it in an idiomatic and easy to use fashion.

[You can find an example of how to use it here](https://github.com/paked/messenger/blob/master/cmd/bot/main.go)
#introduction
Facebook Messenger, sometimes abbreviated Messenger, is an instant messaging service and software application which provides text and voice communication. Integrated with Facebook's web-based Chat feature and built on the open MQTT protocol, Messenger lets Facebook users chat with friends both on mobile and on the main website.

On October 3, 2016, Facebook launched Facebook Messenger Lite to attract more users, particularly, those running Android operating system on 2G network. This light-weight app, 10MB file-size, is a stripped down version of Messenger. It skips notification sounds, voice calling and other features that increase the loading time. However, users would be able to use the basic attributes of Messenger such as emojis, stickers, and photo and link sharing.The app targets the regions or consumer-base that is devoid of high-speed Internet connection. It's available in Kenya, Tunisia, Malaysia, Sri Lanka and Venezuela, and is set to come to other countries later.

Facebook has reported that Facebook Messenger has reached 1 billion monthly active users. David A. Marcus heads Facebook Messenger and had joined Facebook on invitation of Mark Zuckerberg, CEO of Facebook.

## Tips

- Follow the [quickstart](https://developers.facebook.com/docs/messenger-platform/quickstart) guide for getting everything set up!
- You need a Facebook development app, and a Facebook page in order to build things.
- Use [ngrok](https://ngrok.com) to tunnel your locally running bot so that Facebook can reach the webhook.

## Breaking Changes

`paked/messenger` is a pretty stable library however, changes will be made which might break backwards compatibility. For the convenience of its users, these are documented here.


- [23/1/17](https://github.com/paked/messenger/commit/1145fe35249f8ce14d3c0a52544e4a4babdc15a4): Updating timezone type to `float64` in profile struct
- [12/9/16](https://github.com/paked/messenger/commit/47f193fc858e2d710c061e88b12dbd804a399e57): Removing unused parameter `text string` from function `(r *Response) GenericTemplate`.
- [20/5/16](https://github.com/paked/messenger/commit/1dc4bcc67dec50e2f58436ffbc7d61ca9da5b943): Leaving the `WebhookURL` field blank in `Options` will yield a URL of "/" instead of a panic.
- [4/5/16](https://github.com/paked/messenger/commit/eb0e72a5dcd3bfaffcfe88dced6d6ac5247f9da1): The URL to use for the webhook is changable in the `Options` struct. 

## Inspiration

Messenger takes design cues from:

- [`net/http`](https://godoc.org/net/http)
- [`github.com/nickvanw/ircx`](https://github.com/nickvanw/ircx)

## Projects

This is a list of projects use `messenger`. If you would like to add your own, submit a [Pull Request](https://github.com/paked/messenger/pulls/new) adding it below.

- [meme-maker](https://github.com/paked/meme-maker) by @paked: A bot which, given a photo and a caption, will create a macro meme.
- [drone-facebook](https://github.com/appleboy/drone-facebook) by @appleboy: [Drone.io](https://drone.io) plugin which sends Facebook notifications
