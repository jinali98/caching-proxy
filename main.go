package main

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func main() {

	port := flag.Int("port", 0, "port to run the caching proxy on")
	origin := flag.String("origin", "", "origin server URL to forward requests to")
	clearCache := flag.Bool("clear-cache", false, "clear the cache and exit")

	flag.Parse()

	if *clearCache {
		fmt.Println("cache cleared")
		return
	}

	if err := run(*port, *origin); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// run checks the settings, builds the proxy, and starts listening.
// It returns an error instead of exiting, so main stays in charge of quitting.
func run(port int, origin string) error {

	if err := validateOrigin(origin); err != nil {
		return err
	}

	if err := validatePort(port); err != nil {
		return err
	}

	p := &proxy{
		origin: origin,
		client: &http.Client{Timeout: 10 * time.Second},
	}

	// ":3000" means "listen on port 3000 on every network interface".
	addr := fmt.Sprintf(":%d", port)

	fmt.Printf("proxying %s -> %s\n", addr, origin)

	// Hand our proxy to the server and block here forever.
	// ListenAndServe only returns if something goes wrong
	return http.ListenAndServe(addr, p)
}

// proxy holds everything a request needs in order to be forwarded.
// It satisfies http.Handler because it has a ServeHTTP method
type proxy struct {
	origin string       // where to forward requests, e.g. "http://dummyjson.com"
	client *http.Client // used to make the outbound call
}

func (p *proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Where this request should really go.
	target := p.origin + r.URL.RequestURI()

	fmt.Printf("target %s\n", target)

	// Build the OUTBOUND request using same method, same body, new address.
	req, err := http.NewRequestWithContext(r.Context(), r.Method, target, r.Body)
	if err != nil {
		http.Error(w, "bad request", http.StatusInternalServerError)
		return
	}

	res, err := p.client.Do(req)
	if err != nil {
		http.Error(w, "origin unreachable", http.StatusBadGateway)
		return
	}

	defer res.Body.Close()

}

func validatePort(p int) error {

	if p < 1 || p > 65535 {
		return fmt.Errorf("invalid port %d: must be between 1 and 65535", p)
	}

	return nil
}

func validateOrigin(o string) error {
	if strings.TrimSpace(o) == "" {
		return errors.New("origin cannot be empty")
	}

	u, err := url.Parse(o)
	if err != nil {

		return fmt.Errorf("invalid origin %q: %w", o, err)
	}

	if strings.TrimSpace(u.Host) == "" {
		return errors.New("invalid host")
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("origin %q must use http or https, got %q", o, u.Scheme)
	}

	return nil
}
