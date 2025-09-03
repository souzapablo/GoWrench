package connections

import (
	"context"
	"wrench/app"
	"wrench/app/manifest/application_settings"
)

var ErrorLoadConnections []error

func LoadConnections(ctx context.Context) {
	settings := application_settings.ApplicationSettingsStatic

	if settings.Connections == nil {
		return
	}

	addIfError(loadConnectionNats(settings.Connections.Nats))
	addIfError(loadJetStreams(settings.Actions))
	addIfError(loadConnectionsKafka(settings.Connections.Kafka))
	addIfError(loadConnectionsRedis(settings.Connections.Redis))
	addIfError(loadConnectionsDynamodb(ctx, settings.Connections.DynamoDb))
}

func addIfError(err error) {
	if err != nil {
		app.LogError2("Error connections: %v", err)
		ErrorLoadConnections = append(ErrorLoadConnections, err)
	}
}
