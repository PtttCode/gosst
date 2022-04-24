package req

import (
	"github.com/ptttcode/gosst/internal/pkg/ilog"
	"github.com/valyala/fasthttp"
)

func (rh *requestHandle) FastRequest() (err error) {
	fv := rh.fv
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	err = fv.Fc.Do(fv.FastReq, resp)
	if err != nil {
		ilog.GetLogger().Error("请求失败:", string(fv.FastReq.Host()), err.Error())
	}

	return
}

func (rh *requestHandle) NetDial() (err error) {
	fv := rh.fv
	resp, err := fv.Hc.Do(fv.Req)
	if err != nil {
		ilog.GetLogger().Error("Failto to Request server", err)
		return
	}

	err = resp.Body.Close()
	if err != nil {
		ilog.GetLogger().Error("Failed to close tcp connection!", err)
	}

	return
}
