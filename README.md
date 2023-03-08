# What is Turbo?

Turbo is my personal request library used for bypassing cybersecurity checks.
I was using CycleTLS for a while, but I found it has issues and wasn't maintained the best.
I also read things talking about high proxy usage from CycleTLS and it had issues working on my mac.
Since I want DalphanAIO to be working on both PC and Mac, building my own request library was the best route.
♥ 🐬 🚀

# How does it work?

Well pretty much we are trying to edit the client hello packet of our requests. There is a library called uTLS which
makes it easy to edit the client hello packet and http client. The issue is although we can create an http client that
fits our needs, it still uses go's default http library by default. This http library does not have support for actually
using the uTLS client, and causes errors if we try to do so. Therefore, we need to actually have a copy of the http library,
and edit it to use uTLS instead of go's default crypto/tls library.
