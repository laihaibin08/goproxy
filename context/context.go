package context

import (
	"fmt"
	"net/http"
)

type Context map[string]interface{}

func (f *Context) Get(name string) (interface{}, error) {
	v, ok := (*f)[name]
	if !ok {
		return nil, fmt.Errorf("Context(%#v) cannot Get(%#v)", f, name)
	}
	return v, nil
}

func (f *Context) GetString(name string) (string, error) {
	v, ok := (*f)[name]
	if !ok {
		return "", fmt.Errorf("Context(%#v) cannot GetString(%#v)", f, name)
	}
	s, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("Context(%#v) cannot convert %#v to string", f, v)
	}
	return s, nil
}

func (f *Context) GetInt(name string) (int, error) {
	v, ok := (*f)[name]
	if !ok {
		return 0, fmt.Errorf("Context(%#v) cannot GetInt(%#v)", f, name)
	}
	s, ok := v.(int)
	if !ok {
		return 0, fmt.Errorf("Context(%#v) cannot convert %#v to int", f, v)
	}
	return s, nil
}

func (f *Context) GetStringMap(name string) (map[string]string, error) {
	v, ok := (*f)[name]
	if !ok {
		return nil, fmt.Errorf("Context(%#v) cannot GetStringMap(%#v)", f, name)
	}
	s, ok := v.(map[string]string)
	if !ok {
		return nil, fmt.Errorf("Context(%#v) cannot convert %#v to map[string]string", f, v)
	}
	return s, nil
}

func (f *Context) GetHeader(name string) (*http.Header, error) {
	v, ok := (*f)[name]
	if !ok {
		return nil, fmt.Errorf("Context(%#v) cannot GetHeader(%#v)", f, name)
	}
	s, ok := v.(*http.Header)
	if !ok {
		return nil, fmt.Errorf("Context(%#v) cannot convert %#v to *http.Header", f, v)
	}
	return s, nil
}

func (f *Context) GetResponseWriter(name string) (http.ResponseWriter, error) {
	v, ok := (*f)[name]
	if !ok {
		return nil, fmt.Errorf("Context(%#v) cannot GetHeader(%#v)", f, name)
	}
	s, ok := v.(http.ResponseWriter)
	if !ok {
		return nil, fmt.Errorf("Context(%#v) cannot convert %#v to *http.Header", f, v)
	}
	return s, nil
}
