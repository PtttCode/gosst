# gosst
a web requests benchmark tool built by golang.

## Feature
- Easy to use 
- High performance
- Common communications protocol support
- Extensive metrics

## Support Communications Protocol
- Http
- Https
- Socks5 proxy

---

## Go version

`1.14.15`

## Install
`go install github.com/ptttcode/gosst`

## Tutorial

---

### Request to http server

`gosst -c 1000 -n 100000 --dst http://127.0.0.1:9900/hi`

---

### Request to https server

**1. without tls verification**

`gosst -c 1000 -n 100000 --dst https://127.0.0.1:9900/hi`

**2. with tls verification**

on this occasion, you need to provide `ca cert`, `client.crt` and `client.key` files.

`gosst -c 1000 -n 100000 --dst https://127.0.0.1:9900/hi --ca tls/ca.crt --crt tls/client.crt --key tls/client.key`

---

### Request to socks5 proxy server

`gosst -c 1000 -n 100000 --dst http://127.0.0.1:9900/hi --proxy 127.0.0.1:6789`

Also, you can request https server through proxy, just add `--proxy 127.0.0.1:6789` 

---

## Benchmark
The interface of http server just return an easy response with string 'Hey'.  

### gosst
`go run main.go -c 1000 -n 100000 --dst 'http://127.0.0.1:9900'`


<details>
<summary>click to check gosst Test Result</summary>

```
Benchmark 100000 times to http://127.0.0.1:9900 by 1000 concurrency:

Server Address:         127.0.0.1:9900
Api Path:               /
Total Concurrency: 1000
Total Requests: 100000
Failed Requests: 0

Request per second: 96211.94 [req/sec]
Time per request: 8.00 ms
Time taken for benchmark: 1.039372 s

Total Sent: 8800000 bytes
Total read: 11900000 bytes

Percentage of the requests served within a certain time (ms)
50%   7ms
65%   9ms
75%   11ms
85%   14ms
95%   22ms
98%   29ms
99%   34ms
100%  143ms

```
</details>

### apache benchmark
`ab -c 1000 -n 100000 http://127.0.0.1:9900/`

<details>
<summary>click to check ab Test Result</summary>


```
This is ApacheBench, Version 2.3 <$Revision: 1843412 $>                                                                                                                                                    [3/1882]
Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
Licensed to The Apache Software Foundation, http://www.apache.org/

Benchmarking 127.0.0.1 (be patient)
Completed 10000 requests
Completed 20000 requests
Completed 30000 requests
Completed 40000 requests
Completed 50000 requests
Completed 60000 requests
Completed 70000 requests
Completed 80000 requests
Completed 90000 requests
Completed 100000 requests
Finished 100000 requests


Server Software:        
Server Hostname:        127.0.0.1
Server Port:            9900

Document Path:          /
Document Length:        3 bytes

Concurrency Level:      1000
Time taken for tests:   4.482 seconds
Complete requests:      100000
Failed requests:        0
Total transferred:      11900000 bytes
HTML transferred:       300000 bytes
Requests per second:    22311.96 [#/sec] (mean)
Time per request:       44.819 [ms] (mean)
Time per request:       0.045 [ms] (mean, across all concurrent requests)
Transfer rate:          2592.89 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0   21   5.0     21      37
Processing:     5   23   5.2     23      39
Waiting:        1   14   6.1     13      37
Total:         27   45   2.4     44      59

Percentage of the requests served within a certain time (ms)
  50%     44
  66%     45
  75%     45
  80%     46
  90%     47
  95%     49
  98%     50
  99%     55
 100%     59 (longest request)
```
</details>


## Help
if you have some tls certificate questions, you can view https://github.com/jcbsmpsn/golang-https-example. Maybe it can help you!