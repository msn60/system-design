package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/http/httptrace"
	"time"
)

type Item struct {
	ID int `json:"id"`
}

// loadItem simulates a slow data source (100ms per item).
func loadItem(id int) Item {
	time.Sleep(100 * time.Millisecond)
	return Item{ID: id}
}

// Design 1: load everything, serialize everything, then send.
func bulkHandler(w http.ResponseWriter, r *http.Request) {
	items := make([]Item, 0, 5)
	for i := 1; i <= 5; i++ {
		items = append(items, loadItem(i))
	}
	body, err := json.Marshal(items)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(body) // the first byte leaves only now
}

// Design 2: send each item (NDJSON) as soon as it is ready.
func streamHandler(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/x-ndjson")
	enc := json.NewEncoder(w)
	for i := 1; i <= 5; i++ {
		if err := enc.Encode(loadItem(i)); err != nil {
			return
		}
		flusher.Flush() // push this item to the client now
	}
}

// timingWriter records when the handler first writes to the response.
type timingWriter struct {
	http.ResponseWriter
	start      time.Time
	firstWrite time.Duration
}

func (t *timingWriter) mark() {
	if t.firstWrite == 0 {
		t.firstWrite = time.Since(t.start)
	}
}
func (t *timingWriter) WriteHeader(code int) { t.mark(); t.ResponseWriter.WriteHeader(code) }
func (t *timingWriter) Write(b []byte) (int, error) {
	t.mark()
	return t.ResponseWriter.Write(b)
}
func (t *timingWriter) Flush() { t.ResponseWriter.(http.Flusher).Flush() }

// withTiming logs server-side "time before first write" and total time.
func withTiming(name string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tw := &timingWriter{ResponseWriter: w, start: time.Now()}
		next(tw, r)
		fmt.Printf("server %-7s first-write=%4dms total=%4dms\n",
			name, tw.firstWrite.Milliseconds(), time.Since(tw.start).Milliseconds())
	}
}

// measure reports client-observed TTFB and total time (TTLB).
func measure(url string) (ttfb, total time.Duration) {
	start := time.Now()
	trace := &httptrace.ClientTrace{
		GotFirstResponseByte: func() { ttfb = time.Since(start) },
	}
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	io.Copy(io.Discard, resp.Body)
	resp.Body.Close()
	return ttfb, time.Since(start)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/bulk", withTiming("bulk", bulkHandler))
	mux.HandleFunc("/stream", withTiming("stream", streamHandler))
	srv := httptest.NewServer(mux)
	defer srv.Close()

	for _, path := range []string{"/bulk", "/stream"} {
		ttfb, total := measure(srv.URL + path)
		fmt.Printf("client %-7s TTFB=%4dms TTLB=%4dms\n", path, ttfb.Milliseconds(), total.Milliseconds())
	}
}
