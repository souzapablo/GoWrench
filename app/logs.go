package app

import (
	"log"
	"time"

	otelLog "go.opentelemetry.io/otel/log"
)

type WrenchErrorLog struct {
	Message string
	Error   error
}

func LogInfo(msg string) {
	log.Print(msg)

	if Logger != nil {
		Logger.Emit(GetContext(), getRecord(msg, otelLog.SeverityInfo))
	}
}

func LogWarning(msg string) {
	log.Print(msg)

	if Logger != nil {
		Logger.Emit(GetContext(), getRecord(msg, otelLog.SeverityWarn))
	}
}

func LogError(err WrenchErrorLog) {
	log.Print(err)

	if Logger != nil {
		Logger.Emit(GetContext(), getRecord(err.Message, otelLog.SeverityError))
	}
}

func LogError2(msg string, err error) {
	LogError(WrenchErrorLog{Message: msg, Error: err})
}

func getRecord(msg string, severity otelLog.Severity) otelLog.Record {
	var record otelLog.Record
	record.SetSeverity(severity)
	record.AddAttributes(otelLog.KeyValue{Key: "instance", Value: otelLog.StringValue(GetInstanceID())})
	record.SetBody(otelLog.StringValue(msg))
	record.SetTimestamp(time.Now())
	return record
}
