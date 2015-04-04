package filters

import (
	"fmt"
	"github.com/phuslu/goproxy/context"
	"net/http"
)

type Filter interface {
	FilterName() string
	Fetch(*context.Context, *http.Request) (*http.Response, error)
}

type RegisteredFilter struct {
	New func() (Filter, error)
}

var (
	filters map[string]*RegisteredFilter
)

func init() {
	filters = make(map[string]*RegisteredFilter)
}

// Register a Filter
func Register(name string, registeredFilter *RegisteredFilter) error {
	if _, exists := filters[name]; exists {
		return fmt.Errorf("Name already registered %s", name)
	}

	filters[name] = registeredFilter
	return nil
}

// NewFilter creates a new Filter of type "name"
func NewFilter(name string) (Filter, error) {
	filter, exists := filters[name]
	if !exists {
		return nil, fmt.Errorf("hosts: Unknown filter %q", name)
	}
	return filter.New()
}
