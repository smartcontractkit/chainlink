// Copyright 2021 github.com/gagliardetto
// This file has been modified by github.com/gagliardetto
//
// Copyright 2020 dfuse Platform Inc.
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

package system

import (
	"github.com/streamingfast/logging"
	"go.uber.org/zap"
)

var traceEnabled = logging.IsTraceEnabled("solana-go", "github.com/gagliardetto/solana-go/system")
var zlog = zap.NewNop()

func init() {
	logging.Register("github.com/gagliardetto/solana-go/system", &zlog)
}
