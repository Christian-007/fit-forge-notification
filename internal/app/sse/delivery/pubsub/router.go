package pubsub

import (
	"encoding/json"
	"log/slog"

	"github.com/Christian-007/fit-forge-notification/internal/pkg/appdependency"
	"github.com/Christian-007/fit-forge-notification/internal/pkg/messagebroker"
	"github.com/Christian-007/fit-forge-notification/internal/pkg/topics"
	"github.com/ThreeDotsLabs/watermill-googlecloud/pkg/googlecloud"
	"github.com/ThreeDotsLabs/watermill/message"
)

func Routes(router *message.Router, subscriber *googlecloud.Subscriber, appDependencies appdependency.AppDependency) {
	// Add handler into router
	router.AddNoPublisherHandler(
		"fan_out_point_rewards",
		topics.PointsRewarded,
		subscriber,
		func(msg *message.Message) error {
			var payload messagebroker.Message
			err := json.Unmarshal(msg.Payload, &payload)
			if err != nil {
				appDependencies.Logger.Error("[fan_out_point_rewards] failed to unmarshall JSON", slog.String("error", err.Error()))
				return err
			}

			appDependencies.InMemoryMessageBroker.Publish(payload.UserID, payload)

			appDependencies.Logger.Info("[fan_out_point_rewards] successfully fan out point rewards", slog.Any("payload", payload))
			return nil
		},
	)
}
