package httpproxy

import (
	"github.com/golang/glog"
	"github.com/phuslu/goproxy/context"
	"github.com/phuslu/goproxy/httpproxy/plugins"
	"io"
	"net"
	"net/http"
)

type Handler struct {
	http.Handler
	Listener        net.Listener
	Transport       *http.Transport
	RequestFilters  []RequestFilter
	ResponseFilters []ResponseFilter
}

type RequestFilter interface {
	Filter(ctx *context.Context, req *http.Request) (*context.Context, plugins.Plugin, error)
}

type ResponseFilter interface {
	Filter(ctx *context.Context, res *http.Response) (*context.Context, *http.Response, error)
}

func (h Handler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if req.Method != "CONNECT" && !req.URL.IsAbs() {
		if req.TLS != nil {
			req.URL.Scheme = "https"
			if req.Host != "" {
				req.URL.Host = req.Host
			} else {
				req.URL.Host = req.TLS.ServerName
			}
		} else {
			req.URL.Scheme = "http"
			req.URL.Host = req.Host
		}
	}

	ctx := &context.Context{}
	for _, f1 := range h.RequestFilters {
		ctx, plugin, err := f1.Filter(ctx, req)
		if err != nil {
			glog.Infof("ServeHTTP RequestFilter error: %v", err)
			return
		}
		if plugin != nil {
			resp, err := plugin.Fetch(ctx, req)
			if err != nil {
				glog.Infof("ServeHTTP HandleRequest error: %v", err)
				return
			}
			if resp == nil {
				return
			}
			for _, f2 := range h.ResponseFilters {
				ctx, resp, err = f2.Filter(ctx, resp)
				if err != nil {
					glog.Infof("ServeHTTP ResponseFilter error: %v", err)
					return
				}
			}
			if resp != nil {
				for key, values := range resp.Header {
					for _, value := range values {
						rw.Header().Add(key, value)
					}
				}
				rw.WriteHeader(resp.StatusCode)
				io.Copy(rw, resp.Body)
			}
			break
		}
	}
}
