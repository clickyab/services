package xlog

import (
	"context"

	"github.com/sirupsen/logrus"
)

type contextKey int

const ctxKey contextKey = iota

// Get return the logger and initialize it based on context ctxKey
func Get(ctx context.Context) *logrus.Entry {
	fields, ok := ctx.Value(ctxKey).(logrus.Fields)
	entry := logrus.NewEntry(logrus.StandardLogger())
	if ok {
		return entry.WithFields(fields)
	}

	return entry
}

// GetWithError is a shorthand for Get(ctx).WithError(err)
func GetWithError(ctx context.Context, err error) *logrus.Entry {
	return Get(ctx).WithError(err)
}

// GetWithField is a shorthand for Get().WithField()
func GetWithField(ctx context.Context, key string, val interface{}) *logrus.Entry {
	return Get(ctx).WithField(key, val)
}

// GetWithFields is a shorthand for Get().WithFields()
func GetWithFields(ctx context.Context, f logrus.Fields) *logrus.Entry {
	return Get(ctx).WithFields(f)
}

// SetField in the context
func SetField(ctx context.Context, key string, val interface{}) context.Context {
	fields, ok := ctx.Value(ctxKey).(logrus.Fields)
	if !ok {
		fields = make(logrus.Fields)
	}
	fields[key] = val
	return context.WithValue(ctx, ctxKey, fields)
}

// SetFields set the fields for logger
func SetFields(ctx context.Context, fl logrus.Fields) context.Context {
	fields, ok := ctx.Value(ctxKey).(logrus.Fields)
	if !ok {
		fields = make(logrus.Fields)
	}
	for i := range fl {
		fields[i] = fl[i]
	}
	return context.WithValue(ctx, ctxKey, fields)
}
