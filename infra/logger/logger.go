package logger

import (
	"context"
	"time"

	"github.com/rs/zerolog"
)

// Global is the global logger.
//
//nolint:gochecknoglobals // this is acceptable for a logger
var Global = zerolog.New(zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
	w.FormatTimestamp = func(i interface{}) string {
		return i.(string) //nolint:forcetypeassert // we know the type
	}
})).With().Timestamp().Logger()

func init() { //nolint:gochecknoinits // this is acceptable for a logger
	zerolog.TimestampFunc = func() time.Time {
		return time.Now().UTC()
	}
}

// Info logs a message at the info level.
//
//nolint:zerologlint // the caller will dispatch the event
func Info() *zerolog.Event {
	return Global.Info()
}

// Debug logs a message at the debug level.
//
//nolint:zerologlint // the caller will dispatch the event
func Debug() *zerolog.Event {
	return Global.Debug()
}

// Warn logs a message at the warn level.
//
//nolint:zerologlint // the caller will dispatch the event
func Warn() *zerolog.Event {
	return Global.Warn()
}

// Error logs a message at the error level.
//
//nolint:zerologlint // the caller will dispatch the event
func Error() *zerolog.Event {
	return Global.Error()
}

// WithContext creates a child logger with the field added to its context.
func WithContext(ctx context.Context) context.Context {
	return Global.WithContext(ctx)
}

// FromContext returns the logger stored in ctx.
// If no logger is stored in ctx, it returns a disabled logger.
func FromContext(ctx context.Context) *zerolog.Logger {
	return zerolog.Ctx(ctx)
}

// UpdateContext updates the logger stored in ctx.
func UpdateContext(ctx context.Context, cb func(logger *zerolog.Logger) zerolog.Logger) context.Context {
	existing := FromContext(ctx)
	updated := cb(existing)

	return updated.WithContext(ctx)
}
