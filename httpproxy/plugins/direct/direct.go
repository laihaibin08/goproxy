package direct

import (
	"crypto/tls"
	"fmt"
	"github.com/golang/glog"
	"github.com/phuslu/goproxy/context"
	"github.com/phuslu/goproxy/httpproxy/plugins"
	"io"
	"net"
	"net/http"
	"time"
)

type Plugin struct {
	transport http.RoundTripper
	dialer    *net.Dialer
}

func init() {
	plugins.Register("direct", &plugins.RegisteredPlugin{
		New: NewPlugin,
	})
}

func NewPlugin() (plugins.Plugin, error) {
	return &Plugin{
		transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
			DisableCompression: true,
		},
		dialer: &net.Dialer{
			KeepAlive: 5 * time.Minute,
		},
	}, nil
}

func (p *Plugin) PluginName() string {
	return "direct"
}

func (p *Plugin) Fetch(ctx *context.Context, req *http.Request) (*http.Response, error) {
	if req.Method != "CONNECT" {
		req1, err := http.NewRequest(req.Method, req.URL.String(), req.Body)
		if err != nil {
			return nil, fmt.Errorf("DIRECT Fetch %#v error: %#v", req, err)
		}
		req1.Header = req.Header
		res, err := p.transport.RoundTrip(req1)
		if err == nil {
			glog.Infof("%s \"DIRECT %s %s %s\" %d %s", req.RemoteAddr, req.Method, req.URL.String(), req.Proto, res.StatusCode, res.Header.Get("Content-Length"))
		}
		return res, err
	} else {
		glog.Infof("%s \"DIRECT %s %s %s\" - -", req.RemoteAddr, req.Method, req.Host, req.Proto)
		remote, err := p.dialer.Dial("tcp", req.Host)
		if err != nil {
			return nil, err
		}
		rw, err := ctx.GetResponseWriter("rw")
		if err != nil {
			return nil, err
		}
		hijacker, ok := rw.(http.Hijacker)
		if !ok {
			return nil, fmt.Errorf("http.ResponseWriter(%#v) does not implments Hijacker", rw)
		}
		local, _, err := hijacker.Hijack()
		local.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		go io.Copy(remote, local)
		io.Copy(local, remote)
		return nil, nil
	}
}
