package main

import (
	"flag"
	"github.com/golang/glog"
	"github.com/phuslu/goproxy/certutil"
	"github.com/phuslu/goproxy/httpproxy"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"
)

func getCA() (certutil.CA, error) {
	filename := "CA.crt"
	_, err := os.Stat(filename)
	var ca certutil.CA
	if err == nil {
		ca, err = certutil.NewStdCAFromFile(filename)
		if err != nil {
			return nil, err
		}
	} else {
		ca, err = certutil.NewStdCA("GoAgent", 3*365*24*time.Hour, 2048)
		if err != nil {
			return nil, err
		}
		if err = ca.Dump("CA.crt"); err != nil {
			return nil, err
		}
	}
	return ca, nil
}

func main() {
	flag.Set("logtostderr", "true")
	flag.Parse()

	ca, err := getCA()
	if err != nil {
		glog.Fatalf("getCA() failed: %s", err)
	}

	common, err := ReadConfigFile("proxy.ini")
	if err != nil {
		glog.Fatalf("ReadConfigFile() failed: %s", err)
	}

	addr := net.JoinHostPort(common.ListenIp, strconv.Itoa(common.ListenPort))
	ln, err := httpproxy.Listen("tcp4", addr)
	if err != nil {
		glog.Fatalf("Listen(\"tcp4\", %s) failed: %s", addr, err)
	}

	resolver := httpproxy.NewResolver(nil)
	for name, iplist := range common.IplistMap {
		resolver.SetHost(name, iplist)
	}
	for host, name := range common.HostMap {
		resolver.SetCNAME(host, name)
	}

	dialer := &httpproxy.Dialer{
		Timeout:     30 * time.Second,
		KeepAlive:   30 * time.Second,
		DNSResolver: resolver,
	}
	h := httpproxy.Handler{
		Listener: ln,
		Transport: &http.Transport{
			Dial:                  dialer.Dial,
			DialTLS:               dialer.DialTLS,
			TLSHandshakeTimeout:   2 * time.Second,
			ResponseHeaderTimeout: 2 * time.Second,
			DisableKeepAlives:     true,
			DisableCompression:    true,
			Proxy:                 nil,
		},
		RequestFilters: []httpproxy.RequestFilter{
			&httpproxy.ForcehttpsRequestFilter{
				ForcehttpsSites:   common.ForcehttpsSites,
				NoforcehttpsSites: common.NoforcehttpsSites,
			},
			&httpproxy.FakehttpsRequestFilter{
				CA:               ca,
				FakehttpsSites:   common.FakehttpsSites,
				NofakehttpsSites: common.NofakehttpsSites,
			},
			&httpproxy.StripRequestFilter{
				CA: ca,
			},
			&GAERequestFilter{
				AppIDs: common.GaeAppids,
				Scheme: common.GaeMode,
			},
		},
		ResponseFilters: []httpproxy.ResponseFilter{
			&httpproxy.ImageResponseFilter{},
			&httpproxy.RawResponseFilter{},
		},
	}
	s := &http.Server{
		Handler:        h,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	common.WriteSummary(os.Stderr)
	glog.Infof("ListenAndServe on %s\n", h.Listener.Addr().String())
	glog.Exitln(s.Serve(h.Listener))
}
