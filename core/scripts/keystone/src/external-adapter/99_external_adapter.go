package main

// Taken from https://github.com/smartcontractkit/chainlink/blob/4d5fc1943bd6a60b49cbc3d263c0aa47dc3cecb7/core/services/ocr2/plugins/mercury/integration_test.go#L1055
import (
	"fmt"
	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	"math/rand"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
)

func main() {
	// get initial value from env
	btcUsdInitialValue := 0.0
	btcUsdInitialValueEnv := os.Getenv("BTCUSD_INITIAL_VALUE")
	linkInitialValue := 0.0
	linkInitialValueEnv := os.Getenv("LINK_INITIAL_VALUE")
	nativeInitialValue := 0.0
	nativeInitialValueEnv := os.Getenv("NATIVE_INITIAL_VALUE")

	if btcUsdInitialValueEnv == "" {
		fmt.Println("INITIAL_VALUE not set, using default value")
		btcUsdInitialValue = 1000
	} else {
		fmt.Println("INITIAL_VALUE set to ", btcUsdInitialValueEnv)
		val, err := strconv.ParseFloat(btcUsdInitialValueEnv, 64)
		helpers.PanicErr(err)
		btcUsdInitialValue = val
	}

	if linkInitialValueEnv == "" {
		fmt.Println("LINK_INITIAL_VALUE not set, using default value")
		linkInitialValue = 11.0
	} else {
		fmt.Println("LINK_INITIAL_VALUE set to ", linkInitialValueEnv)
		val, err := strconv.ParseFloat(linkInitialValueEnv, 64)
		helpers.PanicErr(err)
		linkInitialValue = val
	}

	if nativeInitialValueEnv == "" {
		fmt.Println("NATIVE_INITIAL_VALUE not set, using default value")
		nativeInitialValue = 2400.0
	} else {
		fmt.Println("NATIVE_INITIAL_VALUE set to ", nativeInitialValueEnv)
		val, err := strconv.ParseFloat(nativeInitialValueEnv, 64)
		helpers.PanicErr(err)
		nativeInitialValue = val
	}

	pctBounds := 0.3
	externalAdapter(btcUsdInitialValue, "4001", pctBounds)
	externalAdapter(linkInitialValue, "4002", pctBounds)
	externalAdapter(nativeInitialValue, "4003", pctBounds)

	select {}
}

func externalAdapter(initialValue float64, port string, pctBounds float64) *httptest.Server {
	// Create a custom listener on the specified port
	listener, err := net.Listen("tcp", "0.0.0.0:"+port)
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
