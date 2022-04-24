package app

import (
	"crypto/tls"
	"github.com/ptttcode/gosst/internal/socks"
	"github.com/valyala/fasthttp"
	"net/http"
)

type FlagVar struct {
	ConcurrentUsers int
	TotalRequests   int64
	ProxyAddr       string
	DstAddr         string

	Sa *socks.SocksAuth

	//tr *http.Transport
	Hc  *http.Client
	Req *http.Request

	Fc      *fasthttp.Client
	FastReq *fasthttp.Request

	Headers     map[string]string
	Body        []byte
	ContentType string
	Method      string

	TlsConfig *tls.Config
}

type FvFunctions interface {
	RequestPrepare()
	FasthttpPrepare()
	FasthttpRelease()
}

var _ FvFunctions = (*FlagVar)(nil)
