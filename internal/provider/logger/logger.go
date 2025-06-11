package logger

import (
	"context"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/richseviora/huego/pkg/logger"
)

type ContextLogger struct {
	ctx context.Context
}

func (c ContextLogger) Debug(message string, fields ...map[string]interface{}) {
	tflog.Debug(c.ctx, message, fields...)
}

func (c ContextLogger) Error(message string, fields ...map[string]interface{}) {
	tflog.Error(c.ctx, message, fields...)
}

func (c ContextLogger) Info(message string, fields ...map[string]interface{}) {
	tflog.Info(c.ctx, message, fields...)
}

func (c ContextLogger) Trace(message string, fields ...map[string]interface{}) {
	tflog.Trace(c.ctx, message, fields...)
}

func (c ContextLogger) Warn(message string, fields ...map[string]interface{}) {
	tflog.Warn(c.ctx, message, fields...)
}

var _ logger.Logger = &ContextLogger{}

func NewContextLogger(ctx context.Context) *ContextLogger {
	return &ContextLogger{ctx: ctx}
}
