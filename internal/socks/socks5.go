package socks

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"golang.org/x/net/proxy"
	"net"
	"net/url"
)

type SocksAuth struct {
	Username string
	Password string
}

func FasthttpSocksDialer(proxyAddr string) fasthttp.DialFunc {
	var (
		u      *url.URL
		err    error
		dialer proxy.Dialer
	)
	if u, err = url.Parse(proxyAddr); err == nil {
		dialer, err = proxy.FromURL(u, proxy.Direct)
	}
	// It would be nice if we could return the error here. But we can't
	// change our API so just keep returning it in the returned Dial function.
	// Besides the implementation of proxy.SOCKS5() at the time of writing this
	// will always return nil as error.

	return func(addr string) (net.Conn, error) {
		if err != nil {
			return nil, err
		}
		return dialer.Dial("tcp", addr)
	}
}

func BuildTcp(conn net.Conn) (err error) {
	sendData := []byte{0x05, 0x01, 0x00}
	if _, err := conn.Write(sendData); err != nil {
		return err
	}
	getData := make([]byte, 3)
	if _, err = conn.Read(getData[:]); err != nil {
		return err
	}
	if getData[0] != 0x05 || getData[1] == 0xFF {
		return errors.New("socks5 first handshake failed!")
	}
	//if getData[1] == 0x02 {
	//	sendData := append(
	//		append(
	//			append(
	//				[]byte{0x01, byte(len(socks5client.Username))},
	//				[]byte(socks5client.Username)...),
	//			byte(len(socks5client.Password))),
	//		[]byte(socks5client.Password)...)
	//	_, _ = conn.Write(sendData)
	//	getData := make([]byte, 3)
	//	if _, err = conn.Read(getData[:]); err != nil {
	//		return err
	//	}
	//	if getData[1] == 0x01 {
	//		return errors.New("username or password not correct,socks5 handshake failed!")
	//	}
	//}
	return nil
}

func RequestToProxyServer(conn net.Conn, reqIp string, reqPort int) (err error) {
	sendData := append(
		append(
			[]byte{0x5, 0x01, 0x00, 0x03, byte(len(reqIp))},
			[]byte(reqIp)...), byte(reqPort>>8),
		byte(reqPort&255))
	if _, err = conn.Write(sendData); err != nil {
		fmt.Println(err)
		return
	}
	getData := make([]byte, 1024)
	if _, err = conn.Read(getData[:]); err != nil {
		return
	}
	if getData[0] != 0x05 || getData[1] != 0x00 {
		return errors.New("socks5 second handshake failed!")
	}
	return
}
