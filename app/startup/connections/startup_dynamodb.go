package connections

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
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

	if dynamodbConn != nil {
		dynamoDbTableConnections = make(map[string]*DynamoDbTableConnection)
		for _, table := range dynamodbConn.Tables {
			var sdkConfig aws.Config
			var awsErr error

			if dynamodbConn.Local {
				loaders := []func(*config.LoadOptions) error{config.WithRegion(dynamodbConn.LocalAwsRegion)}
				loaders = append(loaders,
					config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
						dynamodbConn.LocalAwsAccessKeyId,
						dynamodbConn.LocalAwsSecretAccessKey,
						"",
					)),
				)
				sdkConfig, awsErr = config.LoadDefaultConfig(ctx, loaders...)
			} else {
				sdkConfig, awsErr = config.LoadDefaultConfig(ctx)
			}

			if awsErr != nil {
				return awsErr
			}

			dynamoDbClient := dynamodb.NewFromConfig(sdkConfig, func(o *dynamodb.Options) {
				if dynamodbConn.Local {
					o.BaseEndpoint = aws.String(dynamodbConn.LocalEndpoint)
				}
			})

			var connection = new(DynamoDbTableConnection)
			connection.TableName = table.Name
			connection.DynamoDbClient = dynamoDbClient

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
