package pkg

import (
	"io"
	"log/slog"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

var onExitFlushLoop func()

type ReverseProxy struct {
	Director      func(*http.Request)
	Transport     http.RoundTripper
	FlushInterval time.Duration
	BaseDomain    string
}

func NewReverseProxy(schema string, proxyDomain string, baseDomain string, port int) *ReverseProxy {
	director := func(req *http.Request) {
		host := req.Host
		if strings.Contains(host, ":") {
			host = strings.Split(host, ":")[0]
		}
		if !strings.HasSuffix(host, proxyDomain) {
			slog.Warn("proxy warn", slog.Any("warn", host), slog.Any("proxyDomain", proxyDomain))
			req.URL.Scheme = ""
			req.URL.Host = ""
			return
		}
		subDomain := strings.Split(host, ".")[0]
		subDomain = strings.ReplaceAll(strings.ReplaceAll(strings.ToLower(subDomain), "/", "-"), ".", "-")
		req.URL.Scheme = schema
		req.URL.Host = subDomain + "." + baseDomain + ":" + strconv.Itoa(port)
	}
	return &ReverseProxy{Director: director, BaseDomain: baseDomain}
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func (p *ReverseProxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	transport := p.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}

	outreq := new(http.Request)
	*outreq = *req // includes shallow copies of maps, but okay

	p.Director(outreq)
	outreq.Proto = "HTTP/1.1"
	outreq.ProtoMajor = 1
	outreq.ProtoMinor = 1
	outreq.Close = false

	if outreq.Header.Get("Connection") != "" {
		outreq.Header = make(http.Header)
		copyHeader(outreq.Header, req.Header)
		outreq.Header.Del("Connection")
	}

	if clientIP, _, err := net.SplitHostPort(req.RemoteAddr); err == nil {
		if prior, ok := outreq.Header["X-Forwarded-For"]; ok {
			clientIP = strings.Join(prior, ", ") + ", " + clientIP
		}
		outreq.Header.Set("X-Forwarded-For", clientIP)
	}

	res, err := transport.RoundTrip(outreq)
	if err != nil {
		slog.Error("proxy error", slog.Any("error", err))
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusInternalServerError)
		_, _ = rw.Write([]byte(`{"status":"error"}`))
		return
	}
	defer res.Body.Close()

	copyHeader(rw.Header(), res.Header)

	rw.WriteHeader(res.StatusCode)
	p.copyResponse(rw, res.Body)
}

func (p *ReverseProxy) copyResponse(dst io.Writer, src io.Reader) {
	if p.FlushInterval != 0 {
		if wf, ok := dst.(writeFlusher); ok {
			mlw := &maxLatencyWriter{
				dst:     wf,
				latency: p.FlushInterval,
				done:    make(chan bool),
			}
			go mlw.flushLoop()
			defer mlw.stop()
			dst = mlw
		}
	}

	_, err := io.Copy(dst, src)
	if err != nil {
		slog.Warn("copyResponse", slog.Any("error", err))
	}
}

type writeFlusher interface {
	io.Writer
	http.Flusher
}

type maxLatencyWriter struct {
	dst     writeFlusher
	latency time.Duration

	lk   sync.Mutex // protects Write + Flush
	done chan bool
}

func (m *maxLatencyWriter) Write(p []byte) (int, error) {
	m.lk.Lock()
	defer m.lk.Unlock()
	return m.dst.Write(p)
}

func (m *maxLatencyWriter) flushLoop() {
	t := time.NewTicker(m.latency)
	defer t.Stop()
	for {
		select {
		case <-m.done:
			if onExitFlushLoop != nil {
				onExitFlushLoop()
			}
			return
		case <-t.C:
			m.lk.Lock()
			m.dst.Flush()
			m.lk.Unlock()
		}
	}
}

func (m *maxLatencyWriter) stop() { m.done <- true }
