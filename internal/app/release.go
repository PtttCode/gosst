package app

import (
	"github.com/valyala/fasthttp"
)

func (fv *FlagVar) FasthttpRelease() {
	fasthttp.ReleaseRequest(fv.FastReq)
}
