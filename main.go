package main

import (
	"flag"
	"fmt"
	"github.com/ptttcode/gosst/internal/app"
	"github.com/ptttcode/gosst/internal/pkg/ilog"
	"github.com/ptttcode/gosst/internal/req"
	"github.com/ptttcode/gosst/internal/socks"
	"github.com/valyala/fasthttp"
	"sort"
	"strings"
	"sync"
	"time"
)

type Count struct {
	count int
	m     sync.Mutex
}

var ct = &Count{
	count: 0,
	m:     sync.Mutex{},
}

type sliceString []string

func (s *sliceString) String() string {
	return strings.Join(*s, "\n")
}

func (s *sliceString) Set(str string) error {
	*s = append(*s, str)
	return nil
}

func getSum(s []int64) int64 {
	var sum int64
	for _, i := range s {
		sum += i
	}

	return sum
}

func getRespSize(fv *app.FlagVar) (respSize int) {
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	err := fv.Fc.Do(fv.FastReq, resp)
	if err != nil {
		ilog.GetLogger().Error("请求失败:", string(fv.FastReq.Host()), err.Error())
	} else {
		respSize = len(resp.Body()) + len(resp.Header.Header())
	}

	return
}

func main() {
	var concurrency int
	var totalRequests int64
	var proxyAddr, dstAddr, username, passwd string
	var headers sliceString
	var contentType, body, method string
	var clientCrtPath, clientKeyPath, caCrtPath string

	flag.Var(&headers, "H", "Input the Headers you need. exp: -H 'appid:fdg231^%#a1' -H 'auth:xxxxxxxx' ")
	flag.IntVar(&concurrency, "c", 10, "The number of simulated concurrent users")
	flag.Int64Var(&totalRequests, "n", 1000, "The total number of requests")
	flag.StringVar(&proxyAddr, "proxy", "", "Input the address of proxy server.")
	flag.StringVar(&dstAddr, "dst", "", "Input the address of destination sever.")
	flag.StringVar(&username, "u", "", "Input the username of auth in proxy server.")
	flag.StringVar(&passwd, "p", "", "Input the password of auth in proxy server.")
	flag.StringVar(&body, "b", "", "input body as string during request if you need.")
	flag.StringVar(&contentType, "T", "text/plain", "Input the content type you need.")
	flag.StringVar(&clientCrtPath, "crt", "", "Input the path of 'client.crt' file.")
	flag.StringVar(&clientKeyPath, "key", "", "Input the path of 'client.key' file.")
	flag.StringVar(&caCrtPath, "ca", "", "Input the path of 'ca.crt' file.")
	flag.StringVar(&method, "m", "GET", "Input the request method. exp: GET, POST, PUT, DELETE...")
	flag.Parse()

	dstAddr = strings.Replace(dstAddr, "localhost", "http://127.0.0.1", 1)
	b := []byte(body)

	ilog.InitLogger("Gosst")

	var reqTimeMap = sync.Map{}
	var requestTimeSlice = []int64{}

	ch := make(chan int, totalRequests/2)
	wg := &sync.WaitGroup{}

	headersMap := make(map[string]string)
	for _, s := range headers {
		v := strings.Split(s, ":")
		if len(v) != 2 {
			ilog.GetLogger().Warning("The given header value is not valid: ", s)
		} else {
			headersMap[v[0]] = v[1]
		}
	}

	sa := &socks.SocksAuth{
		Username: username,
		Password: passwd,
	}
	fv := &app.FlagVar{
		ConcurrentUsers: concurrency,
		TotalRequests:   totalRequests,
		ProxyAddr:       proxyAddr,
		DstAddr:         dstAddr,

		Sa: sa,

		Body:        b,
		Headers:     headersMap,
		ContentType: contentType,
		Method:      strings.ToUpper(method),

		TlsConfig: app.GetTLSConfig(caCrtPath, clientCrtPath, clientKeyPath),
	}
	fv.RequestPrepare()
	fv.FasthttpPrepare()
	defer fv.FasthttpRelease()

	var rh = req.NewRequestHandle(fv)
	rh.FastRequest()

	ccyPerRequest := totalRequests / int64(concurrency)
	s := time.Now()

	go func() {
		//wg.Add(1)
		//defer wg.Done()
		for {
			if n, ok := <-ch; ok == true {
				ct.count += n
			} else if len(ch) == 0 {
				return
			}

		}

		//for rep := range ch {
		//	ct.count += rep
		//}
	}()

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		if i == concurrency-1 {
			ccyPerRequest = totalRequests - (ccyPerRequest * (int64(concurrency) - 1))
		}
		go func(cc int, max int64) {
			n := 0
			var tmpSlice []int64
			for j := int64(0); j < max; j++ {
				s := time.Now()
				err := rh.FastRequest()
				// err = rh.NetDial()
				if err != nil {
					n++
				}
				consume := time.Now().Sub(s).Milliseconds()
				tmpSlice = append(tmpSlice, consume)
			}
			reqTimeMap.Store(cc, tmpSlice)
			ch <- n
			wg.Done()
		}(i, ccyPerRequest)
	}

	wg.Wait()
	close(ch)
	t := time.Now().Sub(s)

	reqTimeMap.Range(func(key, value interface{}) bool {
		v := value.([]int64)
		requestTimeSlice = append(requestTimeSlice, v...)
		return true
	})
	//fmt.Println(len(requestTimeSlice), t.Seconds(), concurrency, float64(totalRequests)/t.Seconds(), ct.count)

	sort.Slice(requestTimeSlice, func(i, j int) bool {
		return requestTimeSlice[i] < requestTimeSlice[j]
	})

	distribExp := []float32{0.50, 0.65, 0.75, 0.85, 0.95, 0.98, 0.99, 1.00}
	length := len(requestTimeSlice)

	fmt.Println(fmt.Sprintf("\n\nBenchmark %d times to %s by %d concurrency:\n", totalRequests, dstAddr, concurrency))

	if proxyAddr != "" {
		fmt.Println(fmt.Sprintf("Proxy Server Addr:\t%s", proxyAddr))
	}

	// 请求服务具体信息
	serverAddrSplit := strings.Split(strings.Replace(dstAddr, "http://", "", 1), "/")
	if len(serverAddrSplit) == 1 {
		serverAddrSplit = append(serverAddrSplit, "/")
	}
	fmt.Println(fmt.Sprintf(
		"Server Address:\t\t%s\nApi Path:\t\t%s",
		serverAddrSplit[0], serverAddrSplit[1]),
	)

	// Requests详情
	fmt.Println(fmt.Sprintf(
		"Total Concurrency: %d\nTotal Requests: %d\nFailed Requests: %d\n",
		concurrency, totalRequests, ct.count),
	)

	fmt.Println(fmt.Sprintf(
		"Request per second: %.2f [req/sec]\nTime per request: %.2f ms",
		float64(totalRequests)/t.Seconds(),
		float64(getSum(requestTimeSlice)/int64(len(requestTimeSlice)))),
	)

	// Time详情
	fmt.Println(fmt.Sprintf(
		"Time taken for benchmark: %f s\n", t.Seconds()))

	fmt.Println(
		fmt.Sprintf(
			"Total Sent: %d bytes\nTotal read: %d bytes\n",
			int64(len(fv.FastReq.Header.Header())+len(fv.FastReq.Body()))*totalRequests,
			int64(getRespSize(fv))*totalRequests,
		),
	)

	fmt.Println("Percentage of the requests served within a certain time (ms)")
	for _, v := range distribExp {
		idx := int(float32(length)*v) - 1
		if idx < 0 {
			idx = 0
		}
		fmt.Println(fmt.Sprintf("  %d%% \t%dms", int(v*100), int(requestTimeSlice[idx])))
	}
}
