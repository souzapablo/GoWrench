package handlers

import (
	"context"
	contexts "wrench/app/contexts"
	"wrench/app/manifest/connection_settings"
)

type DynamoDbHandler struct {
	Next  Handler
	Table *connection_settings.DynamodbTableSettings
}

func (handler *DynamoDbHandler) Do(ctx context.Context, wrenchContext *contexts.WrenchContext, bodyContext *contexts.BodyContext) {

}

func (handler *DynamoDbHandler) SetNext(next Handler) {
	handler.Next = next
}
