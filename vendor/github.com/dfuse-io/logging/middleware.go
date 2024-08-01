// Copyright 2019 dfuse Platform Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package logging

import (
	crand "crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"math/rand"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	strackdriverPropagation "contrib.go.opencensus.io/exporter/stackdriver/propagation"
	"github.com/teris-io/shortid"
	"go.opencensus.io/trace"
	"go.opencensus.io/trace/propagation"
	"go.uber.org/zap"
)

var defaultFormat propagation.HTTPFormat = &strackdriverPropagation.HTTPFormat{}
var shortIDGenerator *shortid.Shortid
var traceIDGenerator *defaultIDGenerator

func init() {
	// A new generator using the default alphabet set
	shortIDGenerator = shortid.MustNew(1, shortid.DefaultABC, uint64(time.Now().UnixNano()))

	traceIDGenerator = &defaultIDGenerator{}
	// initialize traceID and spanID generators.
	var rngSeed int64
	for _, p := range []interface{}{
		&rngSeed, &traceIDGenerator.traceIDAdd, &traceIDGenerator.nextSpanID, &traceIDGenerator.spanIDInc,
	} {
		binary.Read(crand.Reader, binary.LittleEndian, p)
	}
	traceIDGenerator.traceIDRand = rand.New(rand.NewSource(rngSeed))
	traceIDGenerator.spanIDInc |= 1

}

// Handler is an http.Handler wrapper to instrument your HTTP server with
// an automatic `zap.Logger` per request (i.e. context).
//
// Logging
//
// This handler is aware of the incoming request's trace id, reading it
// from request headers as configured using the Propagation field. The extracted
// trace id if present is used to configure the actual logger with the field
// `trace_id`.
//
// If the trace id cannot be extracted from the request, an random request id is
// generated and used under the field `req_id`.
type Handler struct {
	// Handler is the handler used to handle the incoming request.
	Next http.Handler

	// Propagation defines how traces are propagated. If unspecified,
	// Stackdriver propagation will be used.
	Propagation propagation.HTTPFormat

	// Actual root logger to instrument with request information
	RootLogger *zap.Logger
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	spanContext, ok := h.extractSpanContext(r)
	var logger *zap.Logger
	if !ok {
		// Not found in the header, check from the context directly than
		span := trace.FromContext(r.Context())
		if span == nil {
			traceIDField := zap.Stringer("trace_id", traceIDGenerator.NewTraceID())
			logger = h.RootLogger.With(traceIDField)
		} else {
			spanContext := span.SpanContext()
			traceID := hex.EncodeToString(spanContext.TraceID[:])
			logger = h.RootLogger.With(zap.String("trace_id", traceID))
		}
	} else {
		traceID := hex.EncodeToString(spanContext.TraceID[:])
		logger = h.RootLogger.With(zap.String("trace_id", traceID))
	}

	ctx := WithLogger(r.Context(), logger)
	h.Next.ServeHTTP(w, r.WithContext(ctx))
}

func (h *Handler) extractSpanContext(r *http.Request) (trace.SpanContext, bool) {
	if h.Propagation == nil {
		return defaultFormat.SpanContextFromRequest(r)
	}

	return h.Propagation.SpanContextFromRequest(r)
}

type defaultIDGenerator struct {
	sync.Mutex

	// Please keep these as the first fields
	// so that these 8 byte fields will be aligned on addresses
	// divisible by 8, on both 32-bit and 64-bit machines when
	// performing atomic increments and accesses.
	// See:
	// * https://github.com/census-instrumentation/opencensus-go/issues/587
	// * https://github.com/census-instrumentation/opencensus-go/issues/865
	// * https://golang.org/pkg/sync/atomic/#pkg-note-BUG
	nextSpanID uint64
	spanIDInc  uint64

	traceIDAdd  [2]uint64
	traceIDRand *rand.Rand
}

// NewSpanID returns a non-zero span ID from a randomly-chosen sequence.
func (gen *defaultIDGenerator) NewSpanID() [8]byte {
	var id uint64
	for id == 0 {
		id = atomic.AddUint64(&gen.nextSpanID, gen.spanIDInc)
	}
	var sid [8]byte
	binary.LittleEndian.PutUint64(sid[:], id)
	return sid
}

// NewTraceID returns a non-zero trace ID from a randomly-chosen sequence.
// mu should be held while this function is called.
func (gen *defaultIDGenerator) NewTraceID() trace.TraceID {
	var tid [16]byte
	// Construct the trace ID from two outputs of traceIDRand, with a constant
	// added to each half for additional entropy.
	gen.Lock()
	binary.LittleEndian.PutUint64(tid[0:8], gen.traceIDRand.Uint64()+gen.traceIDAdd[0])
	binary.LittleEndian.PutUint64(tid[8:16], gen.traceIDRand.Uint64()+gen.traceIDAdd[1])
	gen.Unlock()
	return tid
}
