package imagez

import (
	"bytes"
	"github.com/chai2010/webp"
	"github.com/golang/glog"
	"github.com/phuslu/goproxy/context"
	"github.com/phuslu/goproxy/httpproxy/plugins"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
)

type Plugin struct {
	UnderlayPlugin string
}

func init() {
	plugins.Register("imagez", &plugins.RegisteredPlugin{
		New: NewPlugin,
	})
}

func NewPlugin() (plugins.Plugin, error) {
	return &Plugin{
		UnderlayPlugin: "direct",
	}, nil
}

func (p *Plugin) PluginName() string {
	return "imagez"
}

func (p *Plugin) Fetch(ctx *context.Context, req *http.Request) (*http.Response, error) {
	underlay, err := plugins.NewPlugin(p.UnderlayPlugin)
	if err != nil {
		return nil, err
	}

	resp, err := underlay.Fetch(ctx, req)
	if err != nil {
		return nil, err
	}

	switch resp.Header.Get("Content-Type") {
	case "image/gif":
	case "image/png":
	case "image/jpeg":
		break
	default:
		return resp, nil
	}

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		glog.Errorf("%s decode %#v error: %v", p.PluginName(), resp, err)
		return nil, err
	}

	var b bytes.Buffer
	err = webp.Encode(&b, img, &webp.Options{Lossless: true})
	if err != nil {
		glog.Errorf("%s encode %#v error: %v", p.PluginName(), img, err)
		return nil, err
	}

	resp.Header.Set("Content-Type", "image/webp")
	resp.ContentLength = int64(b.Len())

	return resp, nil
}
