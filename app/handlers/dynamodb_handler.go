package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"
	"wrench/app"
	contexts "wrench/app/contexts"
	"wrench/app/manifest/action_settings"
	"wrench/app/manifest/action_settings/dynamodb_settings"
	"wrench/app/manifest/connection_settings"
	"wrench/app/startup/connections"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type DynamoDbHandler struct {
	Next            Handler
	TableConnection *connections.DynamoDbTableConnection
	TableSettings   *connection_settings.DynamodbTableSettings
	ActionSettings  *action_settings.ActionSettings
}

type dynamodbCommandResult struct {
	HttpStatusCode int
	ErrorMessage   string
	Error          error
	Body           []byte
}

func (result *dynamodbCommandResult) IsSuccess() bool {
	return result.Error == nil && len(result.ErrorMessage) == 0
}

func (handler *DynamoDbHandler) Do(ctx context.Context, wrenchContext *contexts.WrenchContext, bodyContext *contexts.BodyContext) {

	if !wrenchContext.HasError &&
		!wrenchContext.HasCache {

		start := time.Now()
		ctx, span := wrenchContext.GetSpan(ctx, *handler.ActionSettings)
		defer span.End()

		item, err := attributevalue.MarshalMap(bodyContext.GetBodyMap(handler.ActionSettings))

		if err != nil {
			handler.setError(wrenchContext, bodyContext, span, 500, err.Error(), err)
		} else {
			var result dynamodbCommandResult
			if handler.ActionSettings.DynamoDb.Command == dynamodb_settings.DynamoDbCommandCreate {
				result = handler.createCommand(ctx, wrenchContext, bodyContext, item)
			} else if handler.ActionSettings.DynamoDb.Command == dynamodb_settings.DynamoDbCommandUpdate {
				result = handler.updateCommand(ctx, wrenchContext, bodyContext, item)
			} else if handler.ActionSettings.DynamoDb.Command == dynamodb_settings.DynamoDbCommandCreateOrUpdate {
				result = handler.createOrUpdateCommand(ctx, item)
			} else if handler.ActionSettings.DynamoDb.Command == dynamodb_settings.DynamoDbCommandDelete {
				result = handler.deleteCommand(ctx, wrenchContext, bodyContext)
			} else if handler.ActionSettings.DynamoDb.Command == dynamodb_settings.DynamoDbCommandGet {
				result = handler.getCommand(ctx, wrenchContext, bodyContext)
			} else {
				result.ErrorMessage = fmt.Sprintf("The command %v is not implemented", handler.ActionSettings.DynamoDb.Command)
				result.Error = errors.New(result.ErrorMessage)
				result.HttpStatusCode = 500
			}

			if result.IsSuccess() {
				bodyContext.HttpStatusCode = result.HttpStatusCode
				bodyContext.SetBodyAction(handler.ActionSettings, result.Body)
			} else {
				handler.setError(wrenchContext, bodyContext, span, result.HttpStatusCode, result.ErrorMessage, result.Error)
			}
		}

		duration := time.Since(start).Seconds() * 1000
		handler.metricRecord(ctx, duration, 200, string(handler.ActionSettings.DynamoDb.Command), handler.TableConnection.TableName)
	}

	if handler.Next != nil {
		handler.Next.Do(ctx, wrenchContext, bodyContext)
	}
}

func (handler *DynamoDbHandler) SetNext(next Handler) {
	handler.Next = next
}

func (handler *DynamoDbHandler) metricRecord(ctx context.Context, duration float64, statusCode int, command string, tableName string) {
	app.DynamoDbDuration.Record(ctx, duration,
		metric.WithAttributes(
			attribute.Int("dynamodb_status_code", statusCode),
			attribute.String("dynamodb_command", command),
			attribute.String("dynamodb_table_name", tableName),
		),
	)
}

func (handler *DynamoDbHandler) createCommand(ctx context.Context, wrenchContext *contexts.WrenchContext, bodyContext *contexts.BodyContext, item map[string]types.AttributeValue) dynamodbCommandResult {
	var result dynamodbCommandResult

	keys, err := handler.getKeyFromItem(item)

	if err != nil {
		result.ErrorMessage = fmt.Sprintf("Error to get key. Here's why: %v\n", err)
		result.Error = err
		result.HttpStatusCode = 500
		return result
	}

	itemExist, err := handler.getItem(ctx, wrenchContext, bodyContext, keys)

	if (itemExist != nil && itemExist.Item != nil) || err != nil {
		if err != nil {
			result.ErrorMessage = err.Error()
			result.Error = err
			result.HttpStatusCode = 500
		} else {
			result.ErrorMessage = fmt.Sprintf("Conflit! The document already exist in table %v", handler.TableConnection.TableName)
			result.Error = errors.New(result.ErrorMessage)
			result.HttpStatusCode = 409
		}
	} else {

		_, err := handler.TableConnection.DynamoDbClient.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(handler.TableConnection.TableName), Item: item,
		})

		if err == nil {
			result.HttpStatusCode = 201
		} else {
			result.ErrorMessage = fmt.Sprintf("Couldn't add item in table %v. Here's why: %v\n", handler.TableConnection.TableName, err)
			result.Error = err
			result.HttpStatusCode = 500
		}
	}

	return result
}

func (handler *DynamoDbHandler) updateCommand(ctx context.Context, wrenchContext *contexts.WrenchContext, bodyContext *contexts.BodyContext, item map[string]types.AttributeValue) dynamodbCommandResult {
	var result dynamodbCommandResult

	keys, err := handler.getKeyFromItem(item)

	if err != nil {
		result.ErrorMessage = fmt.Sprintf("Error to get key. Here's why: %v\n", err)
		result.Error = err
		result.HttpStatusCode = 500
		return result
	}

	itemExist, err := handler.getItem(ctx, wrenchContext, bodyContext, keys)

	if (itemExist != nil && itemExist.Item == nil) || err != nil {
		if err != nil {
			result.ErrorMessage = err.Error()
			result.Error = err
			result.HttpStatusCode = 500
		} else {
			result.ErrorMessage = fmt.Sprintf("Not found! The document don't exist in table %v", handler.TableConnection.TableName)
			result.Error = errors.New(result.ErrorMessage)
			result.HttpStatusCode = 404
		}
	} else {

		_, err := handler.TableConnection.DynamoDbClient.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(handler.TableConnection.TableName), Item: item,
		})

		if err == nil {
			result.HttpStatusCode = 200
		} else {
			result.ErrorMessage = fmt.Sprintf("Couldn't update item in table %v. Here's why: %v\n", handler.TableConnection.TableName, err)
			result.Error = err
			result.HttpStatusCode = 500
		}
	}

	return result
}

func (handler *DynamoDbHandler) createOrUpdateCommand(ctx context.Context, item map[string]types.AttributeValue) dynamodbCommandResult {
	var result dynamodbCommandResult

	_, err := handler.TableConnection.DynamoDbClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(handler.TableConnection.TableName), Item: item,
	})

	if err == nil {
		result.HttpStatusCode = 200
	} else {
		result.ErrorMessage = fmt.Sprintf("Couldn't update item in table %v. Here's why: %v\n", handler.TableConnection.TableName, err)
		result.Error = err
		result.HttpStatusCode = 500
	}

	return result
}

func (handler *DynamoDbHandler) deleteCommand(ctx context.Context, wrenchContext *contexts.WrenchContext, bodyContext *contexts.BodyContext) dynamodbCommandResult {
	var result dynamodbCommandResult

	keys, err := handler.getKey(wrenchContext, bodyContext)

	if err != nil {
		result.ErrorMessage = fmt.Sprintf("Error to get key. Here's why: %v\n", err)
		result.Error = err
		result.HttpStatusCode = 500
		return result
	}

	itemExist, err := handler.getItem(ctx, wrenchContext, bodyContext, keys)

	if (itemExist != nil && itemExist.Item == nil) || err != nil {
		if err != nil {
			result.ErrorMessage = err.Error()
			result.Error = err
			result.HttpStatusCode = 500
		} else {
			result.ErrorMessage = fmt.Sprintf("Not found! The document don't exist in table %v", handler.TableConnection.TableName)
			result.Error = errors.New(result.ErrorMessage)
			result.HttpStatusCode = 404
		}
	} else {
		key, err := handler.getKey(wrenchContext, bodyContext)
		if err == nil {
			_, err = handler.TableConnection.DynamoDbClient.DeleteItem(ctx, &dynamodb.DeleteItemInput{
				TableName: aws.String(handler.TableConnection.TableName), Key: key,
			})
		}

		if err == nil {
			result.HttpStatusCode = 200
			result.Body = []byte("{}")
		} else {
			result.ErrorMessage = fmt.Sprintf("Couldn't update item in table %v. Here's why: %v\n", handler.TableConnection.TableName, err)
			result.Error = err
			result.HttpStatusCode = 500
		}
	}

	return result
}

func (handler *DynamoDbHandler) getCommand(ctx context.Context, wrenchContext *contexts.WrenchContext, bodyContext *contexts.BodyContext) dynamodbCommandResult {
	var result dynamodbCommandResult

	keys, err := handler.getKey(wrenchContext, bodyContext)

	if err != nil {
		result.ErrorMessage = fmt.Sprintf("Error to get key. Here's why: %v\n", err)
		result.Error = err
		result.HttpStatusCode = 500
		return result
	}

	itemExist, err := handler.getItem(ctx, wrenchContext, bodyContext, keys)

	if (itemExist != nil && itemExist.Item == nil) || err != nil {
		if err != nil {
			result.ErrorMessage = err.Error()
			result.Error = err
			result.HttpStatusCode = 500
		} else {
			result.ErrorMessage = fmt.Sprintf("Not found! The document don't exist in table %v", handler.TableConnection.TableName)
			result.Error = errors.New(result.ErrorMessage)
			result.HttpStatusCode = 404
		}
	} else {
		var itemResult map[string]interface{}
		err = attributevalue.UnmarshalMap(itemExist.Item, &itemResult)
		var jsonArray []byte

		if err == nil {
			jsonArray, err = json.Marshal(itemResult)
		}

		if err == nil {
			result.HttpStatusCode = 200
			result.Body = jsonArray
		} else {
			result.ErrorMessage = fmt.Sprintf("Error convert item. Here's why: %v\n", err)
			result.Error = err
			result.HttpStatusCode = 500
		}
	}

	return result
}

func (handler *DynamoDbHandler) getItem(ctx context.Context, wrenchContext *contexts.WrenchContext, bodyContext *contexts.BodyContext, keys map[string]types.AttributeValue) (*dynamodb.GetItemOutput, error) {

	response, err := handler.TableConnection.DynamoDbClient.GetItem(ctx, &dynamodb.GetItemInput{
		Key: keys, TableName: aws.String(handler.TableConnection.TableName),
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (handler *DynamoDbHandler) setError(wrenchContext *contexts.WrenchContext, bodyContext *contexts.BodyContext, span trace.Span, statusCode int, messageError string, err error) {
	wrenchContext.SetHasError(span, messageError, err)
	bodyContext.SetBodyAction(handler.ActionSettings, []byte(err.Error()))
	bodyContext.HttpStatusCode = statusCode
}

func (handler *DynamoDbHandler) getKeyFromItem(item map[string]types.AttributeValue) (map[string]types.AttributeValue, error) {

	partitionKeyValue, ok := item[handler.TableSettings.PartitionKeyName]
	if !ok {
		return nil, fmt.Errorf("the partition key %v not exist in item", handler.TableSettings.PartitionKeyName)
	}

	var sortKeyValue types.AttributeValue
	if len(handler.TableSettings.SortKeyName) > 0 {
		sortKeyValue, ok = item[handler.TableSettings.SortKeyName]
		if !ok {
			return nil, fmt.Errorf("the sort key %v not exist in item", handler.TableSettings.SortKeyName)
		}
	}

	if sortKeyValue != nil {
		var keyMap = map[string]types.AttributeValue{
			handler.TableSettings.PartitionKeyName: partitionKeyValue,
			handler.TableSettings.SortKeyName:      sortKeyValue}
		return keyMap, nil
	} else {
		var keyMap = map[string]types.AttributeValue{
			handler.TableSettings.PartitionKeyName: partitionKeyValue}
		return keyMap, nil
	}
}

func (handler *DynamoDbHandler) getKey(wrenchContext *contexts.WrenchContext, bodyContext *contexts.BodyContext) (map[string]types.AttributeValue, error) {
	partitionKeyValue := contexts.GetCalculatedValue(handler.ActionSettings.DynamoDb.Key.PartitionKeyValue, wrenchContext, bodyContext, handler.ActionSettings)
	var sortKeyValue interface{}

	if len(handler.TableSettings.SortKeyName) > 0 {
		sortKeyValue = contexts.GetCalculatedValue(handler.ActionSettings.DynamoDb.Key.SortKeyValue, wrenchContext, bodyContext, handler.ActionSettings)
	}
	marshalPartitionKeyValue, err := attributevalue.Marshal(partitionKeyValue)

	if err != nil {
		return nil, err
	}

	if sortKeyValue != nil {
		marshalSortKeyValue, err := attributevalue.Marshal(sortKeyValue)

		if err != nil {
			return nil, err
		}

		var keyMap = map[string]types.AttributeValue{
			handler.TableSettings.PartitionKeyName: marshalPartitionKeyValue,
			handler.TableSettings.SortKeyName:      marshalSortKeyValue}

		return keyMap, nil
	} else {
		var keyMap = map[string]types.AttributeValue{
			handler.TableSettings.PartitionKeyName: marshalPartitionKeyValue}

		return keyMap, nil
	}
}
