package configs

import (
	"context"

	"github.com/omniful/go_commons/config"
)

func Environment(ctx context.Context) string {
	return config.GetString(ctx, "env")
}

func GetDestinationBucketName(ctx context.Context) string {
	return config.GetString(ctx, "s3.bucket")
}

func ImportCitiesPusherEventName(ctx context.Context) string {
	return config.GetString(ctx, "pusher.events.city.import")
}

func ExportCitiesPusherEventName(ctx context.Context) string {
	return config.GetString(ctx, "pusher.events.city.export")
}

func GetPusherChannelName(ctx context.Context) string {
	return config.GetString(ctx, "pusher.channel")
}

func GetCityQueueName(ctx context.Context) string {
	return config.GetString(ctx, "worker.city.name")
}

func GetCityWorkerCount(ctx context.Context) int64 {
	return config.GetInt64(ctx, "worker.city.workerCount")
}

func GetCityRegion(ctx context.Context) string {
	return config.GetString(ctx, "worker.city.region")
}

func GetCityAccount(ctx context.Context) string {
	return config.GetString(ctx, "worker.city.account")
}

func GetCityEndpoint(ctx context.Context) string {
	return config.GetString(ctx, "worker.city.endpoint")
}

func GetCitiesExportBatchSize(ctx context.Context) int {
	return config.GetInt(ctx, "constants.cities.export.batchSize")
}

func GetLocationMappingFromGoogleEnabled(ctx context.Context) bool {
	return config.GetBool(ctx, "featureFlag.locationMappingFromGoogleEnabled")
}

func IsGeocodeAddressStringActive(ctx context.Context) bool {
	return config.GetBool(ctx, "featureFlag.geocodeAddressString")
}
