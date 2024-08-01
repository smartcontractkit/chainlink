package protocol

import (
	"context"
	"fmt"
	"time"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/loghelper"
)

const ReportingPluginTimeoutWarningGracePeriod = 100 * time.Millisecond

func callPlugin[T any](
	ctx context.Context,
	logger loghelper.LoggerWithContext,
	logFields commontypes.LogFields,
	name string,
	maxDuration time.Duration,
	f func(context.Context) (T, error),
) (T, bool) {
	pluginCtx, cancel := context.WithTimeout(ctx, maxDuration)
	defer cancel()

	ins := loghelper.NewIfNotStopped(
		maxDuration+ReportingPluginTimeoutWarningGracePeriod,
		func() {
			logger.MakeChild(logFields).Error(fmt.Sprintf("call to ReportingPlugin.%s is taking too long", name), commontypes.LogFields{
				"maxDuration": maxDuration.String(),
				"gracePeriod": ReportingPluginTimeoutWarningGracePeriod.String(),
			})
		},
	)

	result, err := f(pluginCtx)

	ins.Stop()

	if err != nil {
		logger.MakeChild(logFields).ErrorIfNotCanceled(fmt.Sprintf("call to ReportingPlugin.%s errored", name), ctx, commontypes.LogFields{
			"error": err,
		})
		// failed to get data, nothing to be done
		var zero T
		return zero, false
	}

	return result, true
}
