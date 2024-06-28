# caddy-qrcode

Caddy v2 module to generate a QR code based on the URL provided in the request.

## Installation

You can build Caddy by yourself by installing [xcaddy](https://github.com/caddyserver/xcaddy) and running:
```bash
xcaddy build --with github.com/lum8rjack/caddy-qrcode
```

If you want to clone and make any changes, you can test locally with the following command:
```bash
# Specify the location of the local build
 xcaddy build --with github.com/lum8rjack/caddy-qrcode=./caddy-qrcode
```

### Caddyfile

Below is an example Caddyfile that will return a QR code PNG file based on the text provided in the 'url' parameter of the request.
```
https://test.example.com {
  handle /img/qrcode {
    route {
      qrcode {
        param url
      }
    }
  }
}
```

The request below will generate a QR code for 'google.com'.

```bash
curl https://test.example.com/img/qrcode?url=google.com
```


## References

- [Barcode](https://github.com/boombuler/barcode)
