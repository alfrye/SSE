package main

import (
	"net/http"
	"time"

	"fmt"
)


func main(){
	http.HandleFunc("/stream", eventsHandler)
	http.ListenAndServe(":9000", nil)
}

func eventsHandler(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers to allow all origins. You may want to restrict this to specific origins in a production environment.
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Type")
   
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
   
	// Simulate sending events (you can replace this with real data)
	for i := 0; i < 10; i++ {
	 fmt.Fprintf(w, "data: %s\n\n", fmt.Sprintf("Event %d", i))
	 time.Sleep(2 * time.Second)
	 w.(http.Flusher).Flush()
	}
   
	// Simulate closing the connection
	closeNotify := w.(http.CloseNotifier).CloseNotify()
	<-closeNotify
   }