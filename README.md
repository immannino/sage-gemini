# Sage Gemini

> crawls sitemaps and forwards html based articles to recipients.
> primarily build to send blog posts to my sage colored kindle, Sage Gemini.

## TODO

- Migrate email sending to SES since Kindle blocks outgoing SMS

## Usage

```
go get github.com/immannino/sage-gemini
```

## Env

Expects the following ENV variables:

```
EMAIL_ADDRESS=
EMAIL_PASSWORD=
EMAIL_HOST=
EMAIL_PORT=
RECIPIENT=
```