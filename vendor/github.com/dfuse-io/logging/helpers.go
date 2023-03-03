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
	"flag"

	"go.uber.org/zap"
)

// FlagFields returns all falg as `zap.Field` element for easy logging
func FlagFields(extraFields ...zap.Field) []zap.Field {
	fields := extraFields
	flag.VisitAll(func(visitedFlag *flag.Flag) {
		fields = append(fields, zap.Any(visitedFlag.Name, visitedFlag.Value.(flag.Getter).Get()))
	})

	return fields
}
