package caddy_qrcode

import (
	"fmt"
	"image/png"
	"net/http"
	"strconv"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"go.uber.org/zap"
)

func init() {
	caddy.RegisterModule(QRCode{})
	httpcaddyfile.RegisterHandlerDirective("qrcode", parseCaddyfile)
}

var (
	DEFAULT_SIZE int = 200
)

// QRCode implements an HTTP handler
type QRCode struct {
	Param  string
	Size   int
	logger *zap.Logger
}

// CaddyModule returns the Caddy module information.
func (QRCode) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.qrcode",
		New: func() caddy.Module { return new(QRCode) },
	}
}

// Provision implements caddy.Provisioner.
func (q *QRCode) Provision(ctx caddy.Context) error {
	q.logger = ctx.Logger(q)

	if q.Param == "" {
		return fmt.Errorf("you must specify a qrcode 'param'")
	}

	if q.Size <= 100 {
		return fmt.Errorf("qrcode 'size' must be larger than 100: %d", q.Size)
	}
	return nil
}

// Validate implements caddy.Validator.
func (q *QRCode) Validate() error {
	return nil
}

// ServeHTTP implements caddyhttp.QRCodeHandler.
func (q QRCode) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	p := r.URL.Query().Get(q.Param)
	if p != "" {
		code, err := qr.Encode(p, qr.M, qr.Auto)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return err
		}
		bc, err := barcode.Scale(code, q.Size, q.Size)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return err
		}
		//w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "imge/png")
		png.Encode(w, bc)
		return err
	}

	return next.ServeHTTP(w, r)
}

// UnmarshalCaddyfile implements caddyfile.Unmarshaler.
func (q *QRCode) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		for nesting := d.Nesting(); d.NextBlock(nesting); {
			switch d.Val() {
			case "param":
				if !d.NextArg() {
					return d.ArgErr()
				}
				q.Param = d.Val()
			case "size":
				var err error
				if !d.NextArg() {
					return d.ArgErr()
				}
				q.Size, err = strconv.Atoi(d.Val())
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// parseCaddyfile unmarshals tokens from h into a new QRCode.
func parseCaddyfile(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	q := QRCode{
		Size: DEFAULT_SIZE,
	}
	err := q.UnmarshalCaddyfile(h.Dispenser)
	return q, err
}

// Interface guards
var (
	_ caddy.Provisioner           = (*QRCode)(nil)
	_ caddy.Validator             = (*QRCode)(nil)
	_ caddyhttp.MiddlewareHandler = (*QRCode)(nil)
	_ caddyfile.Unmarshaler       = (*QRCode)(nil)
)
