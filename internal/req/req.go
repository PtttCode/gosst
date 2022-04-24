package req

import "gosst/internal/app"

type requestHandle struct {
	fv *app.FlagVar
}

type requestPointer interface {
	FastRequest() error
	NetDial() error
}

func NewRequestHandle(fv *app.FlagVar) *requestHandle {
	return &requestHandle{fv: fv}
}

var _ requestPointer = (*requestHandle)(nil)
