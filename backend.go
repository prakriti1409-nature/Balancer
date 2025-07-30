package main

import (
    "fmt"
    "net/http"
    "os"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Please provide a port number, e.g. 8081")
        return
    }

    port := os.Args[1]

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello from backend on port %s\n", port)
    })

    fmt.Println("Backend running on port", port)
    http.ListenAndServe(":"+port, nil)
}
