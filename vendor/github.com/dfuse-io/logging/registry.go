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
	"fmt"
	"os"
	"regexp"
	"strings"

	"go.uber.org/zap"
)

type registerConfig struct {
	onUpdate func(newLogger *zap.Logger)
}

// RegisterOption are option parameters that you can set when registering a new logger
// in the system using `Register` function.
type RegisterOption interface {
	apply(config *registerConfig)
}

type registerOptionFunc func(config *registerConfig)

func (f registerOptionFunc) apply(config *registerConfig) {
	f(config)
}

// RegisterOnUpdate enable you to have a hook function that will receive the new logger
// that is going to be assigned to your logger instance. This is useful in some situation
// where you need to update other instances or re-configuring a bit the logger when
// a new one is attached.
//
// This is called **after** the instance has been re-assigned.
func RegisterOnUpdate(onUpdate func(newLogger *zap.Logger)) RegisterOption {
	return registerOptionFunc(func(config *registerConfig) {
		config.onUpdate = onUpdate
	})
}

type LoggerExtender func(*zap.Logger) *zap.Logger

type registryEntry struct {
	logPtr   **zap.Logger
	onUpdate func(newLogger *zap.Logger)
}

var registry = map[string]*registryEntry{}
var defaultLogger = zap.NewNop()

func Register(name string, zlogPtr **zap.Logger, options ...RegisterOption) {
	if zlogPtr == nil {
		panic("the zlog pointer (of type **zap.Logger) must be set")
	}

	if _, found := registry[name]; found {
		panic(fmt.Sprintf("name already registered: %s", name))
	}

	config := registerConfig{}
	for _, opt := range options {
		opt.apply(&config)
	}

	entry := &registryEntry{
		logPtr:   zlogPtr,
		onUpdate: config.onUpdate,
	}

	registry[name] = entry

	logger := defaultLogger
	if *zlogPtr != nil {
		logger = *zlogPtr
	}

	setLogger(entry, logger)
}

func Set(logger *zap.Logger, regexps ...string) {
	for name, entry := range registry {
		if len(regexps) == 0 {
			setLogger(entry, logger)
		} else {
			for _, re := range regexps {
				if regexp.MustCompile(re).MatchString(name) {
					setLogger(entry, logger)
				}
			}
		}
	}
}

// Extend is different than `Set` by being able to re-configure the existing logger set for
// all registered logger in the registry. This is useful for example to add a field to the
// currently set logger:
//
// ```
// logger.Extend(func (current *zap.Logger) { return current.With("name", "value") }, "github.com/dfuse-io/app.*")
// ```
func Extend(extender LoggerExtender, regexps ...string) {
	for name, entry := range registry {
		if *entry.logPtr == nil {
			continue
		}

		if len(regexps) == 0 {
			setLogger(entry, extender(*entry.logPtr))
		} else {
			for _, re := range regexps {
				if regexp.MustCompile(re).MatchString(name) {
					setLogger(entry, extender(*entry.logPtr))
				}
			}
		}
	}
}

// Override sets the given logger on previously registered and next
// registrations.  Useful in tests.
func Override(logger *zap.Logger) {
	defaultLogger = logger
	Set(logger)
}

// TestingOverride calls `Override` (or `Set`, see below) with a development
// logger setup correctly with the right level based on some environment variables.
//
// By default, override using a `zap.NewDevelopment` logger (`info`), if
// environment variable `DEBUG` is set to anything or environment variable `TRACE`
// is set to `true`, logger is set in `debug` level.
//
// If `DEBUG` is set to something else than `true` and/or if `TRACE` is set
// to something else than
func TestingOverride() {
	debug := os.Getenv("DEBUG")
	trace := os.Getenv("TRACE")
	if debug == "" && trace == "" {
		return
	}

	logger, _ := zap.NewDevelopment()

	regex := ""
	if debug != "true" {
		regex = debug
	}

	if regex == "" && trace != "true" {
		regex = trace
	}

	if regex == "" {
		Override(logger)
	} else {
		for _, regexPart := range strings.Split(regex, ",") {
			regexPart = strings.TrimSpace(regexPart)
			if regexPart != "" {
				Set(logger, regexPart)
			}
		}
	}
}

func setLogger(entry *registryEntry, logger *zap.Logger) {
	if entry == nil || logger == nil {
		return
	}

	*entry.logPtr = logger
	if entry.onUpdate != nil {
		entry.onUpdate(logger)
	}
}
