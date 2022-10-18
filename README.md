# GoPrivate

## What & Why

GoPrivate is inspired by [privnote.com](https://privnote.com/).
The main concept is the same: Create notes that are encrypted
and stored on a server. Use the given link to access and decrypt
the notes.

Another goal of GoPrivate is easy deployment. This is achieved by
having a single executable making up the whole service. This executable
acts as the server and also hosts the frontend.
The goal has also affected some architectural choices, one being
SQLite as database.

## How

GoPrivate follows a simple client-server architecture. The server is
written in Golang and the client is a simple web app (named webfront)
created with SvelteJS. The web app is served as a single page application
and interacts with the server using a purpose specific API. The Golang package
[Gin-Gonic](https://github.com/gin-gonic) is used to create the API.

### Encrypting notes

The API and the server is responsible for storing notes and deleting
them upon a note being read. However, it does not perform any encryption
by itself. This task is expected to be handled by the client using
the API.

The web app (webfront) encrypts the note using AES-CBC before sending
it to the backend. IV and encryption key is generated at encryption time.
The IV is stored together with the note. The encryption key is transformed
to a hex-string and appended to the "read note" URL after `#`.
**NB:** This is a POC and I have tried following best practices but
the implementation might be lacking in some aspects. This application is not
intended for production use.
