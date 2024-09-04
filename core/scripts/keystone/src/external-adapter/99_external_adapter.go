package main

// Taken from https://github.com/smartcontractkit/chainlink/blob/4d5fc1943bd6a60b49cbc3d263c0aa47dc3cecb7/core/services/ocr2/plugins/mercury/integration_test.go#L1055
import (
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"net/http/httptest"
)

func main() {
	// Simulating MATIC/USD 
	initialValue := 0.4
	pctBounds := 0.3

	externalAdapter(initialValue, "4001", pctBounds)
	externalAdapter(initialValue, "4002", pctBounds)
	externalAdapter(initialValue, "4003", pctBounds)

	select {}
}

func externalAdapter(initialValue float64, port string, pctBounds float64) *httptest.Server {
	// Create a custom listener on the specified port
	listener, err := net.Listen("tcp", "localhost:"+port)
	if err != nil {
		panic(err)
	}

	mid := initialValue
	// we make step a tenth of the pctBounds to give ample room for the value to move
	step := mid * pctBounds / 10
	bid := mid - step
	ask := mid + step
	// Calculate the floor and ceiling based on percentages of the initial value
	ceiling := float64(initialValue) * (1 + pctBounds)
	floor := float64(initialValue) * (1 - pctBounds)

	handler := http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusOK)
		// [floor <= bid <= mid <= ask <= ceiling]
		mid = adjustValue(mid, step, floor, ceiling)
		bid = adjustValue(mid, step, floor, mid)
		ask = adjustValue(mid, step, mid, ceiling)

		resp := fmt.Sprintf(`{"result": {"bid": %.4f, "mid": %.4f, "ask": %.4f}}`, bid, mid, ask)

		_, herr := res.Write([]byte(resp))
		if herr != nil {
			fmt.Printf("failed to write response: %v", herr)
		}
	})

	ea := &httptest.Server{
		Listener: listener,
		Config:   &http.Server{Handler: handler},
	}
	ea.Start()

	fmt.Print("Mock external adapter started at ", ea.URL, "\n")
	fmt.Printf("Initial value: %.4f, Floor: %.4f, Ceiling: %.4f\n", initialValue, floor, ceiling)
	return ea
}

// adjustValue takes a starting value and randomly shifts it up or down by a step.
// It ensures that the value stays within the specified bounds.
func adjustValue(start, step, floor, ceiling float64) float64 {
	// Randomly choose to increase or decrease the value
	if rand.Intn(2) == 0 {
		step = -step
	}

	// Apply the step to the starting value
	newValue := start + step

	// Ensure the value is within the bounds
	if newValue < floor {
		newValue = floor
	} else if newValue > ceiling {
		newValue = ceiling
	}

	return newValue
}
