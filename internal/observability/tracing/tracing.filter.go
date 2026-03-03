package tracing

import "strings"

type TraceFilter interface {
	ShouldTrace(path string) bool
}

type traceFilterImpl struct {
	exact    map[string]struct{}
	prefixes []string
}

func NewTraceFilter() TraceFilter {
	return &traceFilterImpl{
		exact: map[string]struct{}{
			"/api/v1/healthz": {},
			"/api/v1/metrics": {},
		},
		prefixes: []string{
			"/swagger/",
			"/static/",
			"/favicon.ico",
		},
	}
}

func (f *traceFilterImpl) ShouldTrace(path string) bool {
	if _, ok := f.exact[path]; ok {
		return false
	}

	for _, p := range f.prefixes {
		if strings.HasPrefix(path, p) {
			return false
		}
	}

	return true
}
