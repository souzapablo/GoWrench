package connections

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/sns/types"

	"wrench/app/manifest/connection_settings"
)

type DynamoDbTableConnection struct {
	DynamoDbClient *dynamodb.Client
	TableName      string
}

var dynamoDbTableConnections map[string]*DynamoDbTableConnection

func loadConnectionsDynamodb(ctx context.Context, dynamodbConn *connection_settings.DynamodbConnectionSettings) error {

	var err error

	for _, table := range dynamodbConn.Tables {
		var connection = new(DynamoDbTableConnection)
		connection.TableName = table.Name

		_, err := connection.DynamoDbClient.DescribeTable(
			ctx, &dynamodb.DescribeTableInput{TableName: aws.String(connection.TableName)},
		)
		if err != nil {
			var notFoundEx *types.ResourceNotFoundException
			if errors.As(err, &notFoundEx) {
				log.Printf("Table %v does not exist.\n", connection.TableName)
			} else {
				log.Printf("Couldn't determine existence of table %v. Here's why: %v\n", connection.TableName, err)
			}
		}

		if err != nil {
			return err
		}

		dynamoDbTableConnections[table.Id] = connection
	}

	return err
}

func GetDynamoDbTableConnection(tableId string) (*DynamoDbTableConnection, error) {
	if len(tableId) == 0 ||
		len(dynamoDbTableConnections) == 0 ||
		dynamoDbTableConnections[tableId] == nil {
		return nil, fmt.Errorf("DynamoDb without connection to tableId: %v", tableId)
	}

	return dynamoDbTableConnections[tableId], nil
}
