package negronidog

import (
	//"fmt"
	"github.com/DataDog/datadog-go/statsd"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type NegroniDog struct {
	Cli *statsd.Client
}

func NewMiddleWare(statsdHost string, namespace string, tags []string) *NegroniDog {
	c, err := statsd.New(statsdHost)
	if err != nil {
		panic(err)
	}
	nd := NegroniDog{
		Cli: c,
	}
	// prefix every metric with the namespace
	c.Namespace = namespace + "."
	c.Tags = tags
	return &nd
}

func (nd *NegroniDog) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()
	next(rw, r)

	statPrefix := "http"
	if mux.CurrentRoute(r) != nil && mux.CurrentRoute(r).GetName() != "" {
		statPrefix += "." + mux.CurrentRoute(r).GetName()
	}

	//resp time
	responseTime := time.Since(start)
	statName := strings.Join([]string{statPrefix, "resp_time"}, ".")
	nd.Cli.Histogram(statName, responseTime.Seconds(), nil, 1)
	//resp code
	statName = strings.Join([]string{statPrefix, "status_code", strconv.Itoa(rw.(negroni.ResponseWriter).Status())}, ".")
	nd.Cli.Count(statName, 1, nil, 1)

}
