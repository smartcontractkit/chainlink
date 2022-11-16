package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
	"time"
)

var reportsCount uint64

func handleReports(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		fmt.Printf("%s is not supported\n", req.Method)
		return
	}

	atomic.AddUint64(&reportsCount, 1)

	// fmt.Println("POST /reports called")

	// b, err := io.ReadAll(req.Body)
	// if err != nil {
	//     fmt.Println("error reading body", err)
	// }
	// fmt.Printf("RAW BINARY received body: 0x%x\n", b)
	// fmt.Printf("STRING received body %s\n", b)
}

func main() {
	http.HandleFunc("/reports", handleReports)

	fmt.Println("running server on :3000")
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		for {
			<-ticker.C
			c := atomic.SwapUint64(&reportsCount, 0)
			fmt.Printf("%s - POST /reports called %d times\n", time.Now().String(), c)
		}
	}()

	// nolint
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}
