package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type requestDetails struct {
	Start                time.Time
	DNSStart             time.Duration
	DNSDone              time.Duration
	GotFirstResponseByte time.Duration
	GotLastResponseByte  time.Duration
	WroteRequest         time.Duration

	Canceled bool
	Deadline bool
	Failed   bool
}

func (r *requestDetails) String() string {
	return fmt.Sprintf("DNSDone: %v, TTFB: %v, TTLB: %v, C: %v, D: %v, F: %v", r.DNSDone, r.GotFirstResponseByte, r.GotLastResponseByte, r.Canceled, r.Deadline, r.Failed)
}

func newHTTPClient(ctx context.Context, wg *sync.WaitGroup, flush time.Duration) *http.Client {
	// create a transport here to facilitate setting MaxIdleConnsPerHost
	// if we don't set this we could exhaust ports because of TIME_WAIT
	// default timeout for TIME_WAIT on linux is 120s
	// SO_LINGER 0 isn't really an option with connection: close semantics
	// but could be useful when we exit with a cancel if we allow thousands
	// of connections in the pool
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	go func(ctx context.Context, wg *sync.WaitGroup, transport *http.Transport, flush time.Duration) {
		wg.Add(1)
		defer wg.Done()

		loop := time.After(flush)
		log.Println("INFO: starting http.Transport.CloseIdleConnections flush routine")
		for {
			select {
			case <-loop:
				log.Println("INFO: executing http.Transport.CloseIdleConnections flush")
				transport.CloseIdleConnections()
				loop = time.After(flush)
			case <-ctx.Done():
				log.Println("INFO: http.Transport.CloseIdleConnections flush routine received completion signal")
				return
			default:
			}
		}
	}(ctx, wg, transport, flush)

	return &http.Client{Transport: transport}
}

func do(ctx context.Context, client *http.Client, endpoint string) {

	// allow argument for url to call
	// probably include channel to send timing results
	// use timeout context

	details := &requestDetails{}
	trace := &httptrace.ClientTrace{
		DNSStart: func(_ httptrace.DNSStartInfo) {
			details.DNSStart = time.Since(details.Start)
		},
		DNSDone: func(_ httptrace.DNSDoneInfo) {
			details.DNSDone = time.Since(details.Start)
		},
		GotFirstResponseByte: func() {
			details.GotFirstResponseByte = time.Since(details.Start)
		},
		WroteRequest: func(_ httptrace.WroteRequestInfo) {
			details.WroteRequest = time.Since(details.Start)
		},
	}

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		log.Fatal(err)
	}

	// test to see what happens
	rctx := httptrace.WithClientTrace(ctx, trace)
	rctx, cancel := context.WithTimeout(rctx, 1*time.Second)
	req = req.WithContext(rctx)
	defer cancel()

	start := time.Now()
	details.Start = start

	// request starts
	//resp, err := http.DefaultClient.Do(req)
	resp, err := client.Do(req)
	if err != nil {
		details.Failed = true
		if rctx.Err() != nil {
			if rctx.Err() == context.Canceled {
				details.Canceled = true
			}

			if rctx.Err() == context.DeadlineExceeded {
				details.Deadline = true
			}
		}
		log.Printf("INFO: %v", details)
		return
	}
	defer resp.Body.Close()

	// ensure we read the body so we can reuse connections
	_, err = io.Copy(ioutil.Discard, resp.Body)
	if err != nil {
		log.Fatal(err) // TODO do better
	}

	// at this point we should be at ttlb for our request
	// from the client pov
	details.GotLastResponseByte = time.Since(details.Start)

	// send me via a channel
	log.Printf("INFO: %v", details)
}

func executor(ctx context.Context, wg *sync.WaitGroup, client *http.Client, endpoint string) {
	wg.Add(1)
	defer wg.Done()

	loop := time.After(1 * time.Second)

	for {
		select {
		case <-loop:

			// logically this doesn't make a request happen every second
			// it makes a request happen one second after the previous request has finished
			// we could probably just spawn the request in the background here, but for now
			// requests are quick enough we're feeling okay
			// if we propagate our context maybe it will be okay?
			do(ctx, client, endpoint)
			loop = time.After(1 * time.Second)
		case <-ctx.Done():
			log.Println("INFO: executor routine received completion signal")
			return
		default:
		}
	}
}

func commandLine(endpoint *string, flush *time.Duration, rps *int) error {

	flag.DurationVar(flush, "flush", 60*time.Second, "specifies how often to flush using a duration; default 60s")
	flag.IntVar(rps, "rps", 2, "specifies requests per second; default 2")
	flag.StringVar(endpoint, "url", "", "specifies the url to query")

	flag.Parse()

	if *endpoint == "" {
		return fmt.Errorf("url argument cannot be empty")
	}

	// TODO
	// this seems to do a crappy job parsing the url, look into it more
	_, err := url.Parse(*endpoint)
	if err != nil {
		return fmt.Errorf("url argument cannot be parsed: %v", err)
	}

	// arbitrary right now
	if *rps < 1 || *rps > 10 {
		return fmt.Errorf("rps argument must be between 1 and 10, got %v", *rps)
	}

	return nil
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.LUTC)

	var endpoint string
	var flush time.Duration
	var rps int

	err := commandLine(&endpoint, &flush, &rps)
	if err != nil {
		log.Fatalf("ERROR: %v", err)
	}

	log.Println("INFO: Starting application...")

	// signal handling
	s := make(chan os.Signal, 1)
	signal.Notify(s,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	// create context
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	// create WaitGroup to keep track of go routines
	var wg sync.WaitGroup

	// create client
	client := newHTTPClient(ctx, &wg, flush)

	// start executing
	for i := 1; i <= rps; i++ {
		log.Printf("INFO: starting executor %v", i)
		go executor(ctx, &wg, client, endpoint)
	}

	// wait for signal
	<-s

	//
	log.Println("INFO: signal received calling cancel")
	cancel()

	// wait until go routines are complete
	log.Println("INFO: waiting for executors to cleanup")
	wg.Wait()
}
