package appdependency

import (
	"log/slog"

	"github.com/Christian-007/fit-forge-notification/internal/pkg/cache"
	"github.com/Christian-007/fit-forge-notification/internal/pkg/messagebroker"
	"github.com/Christian-007/fit-forge-notification/internal/pkg/security"
)

type AppDependency struct {
	AppDependencyOptions
}

type AppDependencyOptions struct {
	Logger                *slog.Logger
	InMemoryMessageBroker *messagebroker.InMemoryMessageBroker
	RedisClient           *cache.RedisCache
	SecretManagerClient   security.SecretManageProvider
}

func NewAppDependency(options AppDependencyOptions) AppDependency {
	return AppDependency{
		options,
	}
}
