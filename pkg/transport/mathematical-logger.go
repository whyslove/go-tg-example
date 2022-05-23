// GENERATED BY 'T'ransport 'G'enerator. DO NOT EDIT.
package transport

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/seniorGolang/dumper/viewer"
	"test-tg/pkg/interfaces"
	"time"
)

type loggerMathematical struct {
	next interfaces.Mathematical
	log  zerolog.Logger
}

func loggerMiddlewareMathematical(log zerolog.Logger) MiddlewareMathematical {
	return func(next interfaces.Mathematical) interfaces.Mathematical {
		return &loggerMathematical{
			log:  log,
			next: next,
		}
	}
}

func (m loggerMathematical) Add(ctx context.Context, a int, b int) (result int, err error) {
	log := m.log.With().Str("service", "Mathematical").Str("method", "add").Logger()
	if ctx.Value(headerRequestID) != nil {
		log = log.With().Interface("requestID", ctx.Value(headerRequestID)).Logger()
	}
	defer func(begin time.Time) {
		fields := map[string]interface{}{
			"request": viewer.Sprintf("%+v", requestMathematicalAdd{
				A: a,
				B: b,
			}),
			"response": viewer.Sprintf("%+v", responseMathematicalAdd{Result: result}),
			"took":     time.Since(begin).String(),
		}
		if err != nil {
			log.Error().Err(err).Fields(fields).Msg("call add")
			return
		}
		log.Info().Fields(fields).Msg("call add")
	}(time.Now())
	return m.next.Add(ctx, a, b)
}