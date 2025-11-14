package handlers

import (
	"context"
	contexts "wrench/app/contexts"
	settings "wrench/app/manifest/action_settings"
)

type FuncSignatureHandler struct {
	ActionSettings *settings.ActionSettings
	Next           Handler
}

func (handler *FuncSignatureHandler) Do(ctx context.Context, wrenchContext *contexts.WrenchContext, bodyContext *contexts.BodyContext) {

	if !wrenchContext.HasError && !wrenchContext.HasCache {

		ctxSpan, span := wrenchContext.GetSpan(ctx, *handler.ActionSettings)
		ctx = ctxSpan
		defer span.End()

		//signSetting := handler.ActionSettings.Func.Sign

		// rsa, error :=

		// 	bodyContext.SetBodyAction(handler.ActionSettings, []byte(signatureValue))
	}
	if handler.Next != nil {
		handler.Next.Do(ctx, wrenchContext, bodyContext)
	}
}

func (handler *FuncSignatureHandler) SetNext(next Handler) {
	handler.Next = next
}
