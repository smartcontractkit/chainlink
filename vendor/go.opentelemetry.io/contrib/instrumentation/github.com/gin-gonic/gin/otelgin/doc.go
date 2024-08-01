// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package otelgin instruments the github.com/gin-gonic/gin package.
//
// Currently there are two ways the code can be instrumented. One is
// instrumenting the routing of a received message (the Middleware function)
// and instrumenting the response generation through template evaluation (the
// HTML function).
package otelgin // import "go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
