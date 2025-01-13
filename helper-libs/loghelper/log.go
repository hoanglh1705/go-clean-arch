package loghelper

import (
	"context"
	"go-clean-arch/helper-libs/commonhelper"
)

func getTraceIdFromContext(ctx context.Context) any {
	return ctx.Value(commonhelper.ContextKeyType_TraceId)
}
