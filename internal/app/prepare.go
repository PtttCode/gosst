package app

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/ptttcode/gosst/internal/pkg/ilog"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"
	"golang.org/x/net/proxy"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

func GetTLSConfig(caCrtPath, clientCrtPath, clientKeyPath string) *tls.Config {
	tc := &tls.Config{
		//RootCAs:            pool,
		InsecureSkipVerify: true,
		//Certificates:       []tls.Certificate{crt},
	}
	pool := x509.NewCertPool()
	caCrt, err := ioutil.ReadFile(caCrtPath)
	if err != nil {
		ilog.GetLogger().Warning(fmt.Sprintf("read %s error! failed to build https config!", caCrtPath), err)
		return tc
	}
	pool.AppendCertsFromPEM(caCrt)

	crt, err := tls.LoadX509KeyPair(clientCrtPath, clientKeyPath)
	if err != nil {
		ilog.GetLogger().Warning(fmt.Sprintf("read %s and %s error! failed to build https config!", clientCrtPath, clientKeyPath), err)
		return tc
	}

	tc.RootCAs = pool
	tc.Certificates = []tls.Certificate{crt}
	tc.InsecureSkipVerify = false

	return tc
}

func (fv *FlagVar) RequestPrepare() {

	var tr = &http.Transport{
		TLSClientConfig:       fv.TlsConfig,
		TLSHandshakeTimeout:   5 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
	}

	var dialer proxy.Dialer
	var err error
	if fv.ProxyAddr != "" {
		dialer, err = proxy.SOCKS5(
			"tcp",
			fv.ProxyAddr,
			//"socks5://127.0.0.1:6789",
			&proxy.Auth{User: fv.Sa.Username, Password: fv.Sa.Password},
			&net.Dialer{
				Timeout: 30 * time.Second,
				//KeepAlive: 30 * time.Second,
			},
		)
		if err != nil {
			ilog.GetLogger().Error("Get Socks Dialer Error!", err)
			return
		}
		tr.DialContext = func(ctx context.Context, network, addr string) (conn net.Conn, err error) {
			conn, err = dialer.Dial(network, addr)
			if err != nil {
				ilog.GetLogger().Error("Failed to Dial", err)
			}
			return
		}
	}

	fv.Hc = &http.Client{
		Timeout:   30 * time.Second,
		Transport: tr,
	}

	b := bytes.NewBuffer(fv.Body)
	fv.Req, err = http.NewRequest(fv.Method, fv.DstAddr, b)
	if err != nil {
		ilog.GetLogger().Error("Failed to New Request!", err)
		return
	}

	// set headers
	fv.Req.Header.Set("Content-Type", fv.ContentType)
	for k, v := range fv.Headers {
		fv.Req.Header.Set(k, v)
	}

}

func (fv *FlagVar) FasthttpPrepare() {
	c := &fasthttp.Client{
		TLSConfig: fv.TlsConfig,
		//MaxConnsPerHost:     20000,
		ReadTimeout:         5 * time.Second,
		WriteTimeout:        5 * time.Second,
		MaxConnWaitTimeout:  5 * time.Second,
		MaxConnDuration:     5 * time.Second,
		MaxIdleConnDuration: 5 * time.Second,
	}

	if fv.ProxyAddr != "" {
		c.Dial = fasthttpproxy.FasthttpSocksDialer("socks5://" + fv.ProxyAddr)
	}

	fv.Fc = c
	fv.FastReq = fasthttp.AcquireRequest()
	fv.FastReq.SetRequestURI(fv.DstAddr)
	fv.FastReq.SetBody(fv.Body)

	// set headers
	fv.FastReq.Header.SetMethod(fv.Method)
	fv.FastReq.Header.Set("Content-Type", fv.ContentType)
	for k, v := range fv.Headers {
		fv.FastReq.Header.Set(k, v)
	}

}
