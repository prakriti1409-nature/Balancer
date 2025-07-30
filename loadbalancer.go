package main

import (
    "fmt"
    "log"
    "net/http"
    "net/http/httputil"
    "net/url"
    "os"
    "sync"
    "sync/atomic"
    "time"
)

var backends = []string{
    "http://localhost:8081",
    "http://localhost:8082",
    "http://localhost:8083",
}

var counter uint64

// Track request counts for each backend
var stats = make(map[string]int)
var mu sync.Mutex // to safely update stats

func getNextBackend() string {
    i := atomic.AddUint64(&counter, 1)
    return backends[(int(i)-1)%len(backends)]
}

// proxyHandler forwards the request to a backend
func proxyHandler(w http.ResponseWriter, r *http.Request) {
    target := getNextBackend()
    backendURL, _ := url.Parse(target)
    w.Header().Set("X-Backend-Server", target)

    // Log to console (green text)
    log.Printf("\033[32mForwarding request to %s\033[0m", target)

    // Log request to file
    logToFile(r, target)

    // Update stats
    mu.Lock()
    stats[target]++
    mu.Unlock()

    proxy := httputil.NewSingleHostReverseProxy(backendURL)
    proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
        log.Printf("Error proxying to %s: %v", target, err)
        http.Error(w, "Backend unavailable", http.StatusBadGateway)
    }
    proxy.ServeHTTP(w, r)
}

// logToFile appends request info to requests.log
func logToFile(r *http.Request, backend string) {
    f, err := os.OpenFile("requests.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Printf("Could not open log file: %v", err)
        return
    }
    defer f.Close()
    line := fmt.Sprintf("%s - %s %s -> %s\n",
        time.Now().Format(time.RFC3339),
        r.Method,
        r.URL.Path,
        backend,
    )
    f.WriteString(line)
}

func main() {
    log.Println("Load balancer running on port 8080")

    // Health-check page
    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "text/html")
        w.Write([]byte(`
            <html>
            <head><title>Health Status</title></head>
            <body style="font-family:sans-serif; text-align:center; margin-top:80px; background-color:#f0f2f5;">
                <div style="display:inline-block; padding:30px; background:white; border-radius:12px; box-shadow:0 4px 12px rgba(0,0,0,0.1);">
                    <h1 style="color:green;">Healthy</h1>
                    <p style="font-size:18px; color:#555;">Load balancer is running and reachable.</p>
                    <button onclick="window.location.href='/'"
                            style="margin-top:20px; padding:10px 20px; font-size:16px; cursor:pointer; background-color:#007BFF; color:white; border:none; border-radius:6px;">
                        Back to Home
                    </button>
                </div>
            </body>
            </html>
        `))
    })

    // Stats page
    http.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "text/html")
        mu.Lock()
        defer mu.Unlock()

        html := `<html><head><title>Load Balancer Stats</title></head>
                 <body style="font-family:sans-serif; text-align:center; margin-top:50px;">
                 <h2>Request Count per Backend</h2>
                 <table border="1" style="margin:auto; border-collapse: collapse;">
                 <tr><th style="padding:8px;">Backend</th><th style="padding:8px;">Requests</th></tr>`
        for backend, count := range stats {
            html += fmt.Sprintf("<tr><td style='padding:8px;'>%s</td><td style='padding:8px;'>%d</td></tr>", backend, count)
        }
        html += "</table></body></html>"

        w.Write([]byte(html))
    })

    // Homepage
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path == "/" {
            w.Header().Set("Content-Type", "text/html")
            w.Write([]byte(`
                <html>
                <head><title>Go Load Balancer</title></head>
                <body style="font-family:sans-serif; text-align:center; margin-top:50px;">
                    <h2>Go Load Balancer</h2>
                    <button onclick="window.location.href='/health'"
                            style="padding:10px 20px; font-size:16px; cursor:pointer;">
                        Check Health
                    </button>
                    <button onclick="window.location.href='/stats'"
                            style="padding:10px 20px; font-size:16px; margin-left:10px; cursor:pointer;">
                        View Stats
                    </button>
                    <p style="margin-top:20px;">
                        Requests to other paths will be forwarded to backend servers automatically.
                    </p>
                </body>
                </html>
            `))
            return
        }
        proxyHandler(w, r)
    })

    port := os.Getenv("PORT")
if port == "" {
    port = "8080" // local default
}
log.Fatal(http.ListenAndServe(":"+port, nil))

}
