// We can't currently easily pass flags to k8s containers because of
// https://github.com/bazelbuild/rules_go/issues/1296
//
// So we're using this dummy main.go to set flags to the values we want.
//
// From:
// https://github.com/buchgr/bazel-remote (Apache 2 License)
// .. with minor modifications
package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"

	"github.com/buchgr/bazel-remote/cache"
)

func main() {
	host := flag.String("host", "", "Host to bind the http server")
	port := flag.Int("port", 8080, "The port the HTTP server listens on")
	dir := flag.String("dir", ".",
		"Directory path where to store the cache contents")
	maxSize := flag.Int64("max_size", 5,
		"The maximum size of the remote cache in GiB")
	flag.Parse()

	e := cache.NewEnsureSpacer(0.95, 0.5)
	h := cache.NewHTTPCache(*dir, *maxSize*1024*1024*1024, e)
	s := &http.Server{
		Addr:    *host + ":" + strconv.Itoa(*port),
		Handler: http.HandlerFunc(h.CacheHandler),
	}
	log.Fatal(s.ListenAndServe())
}
