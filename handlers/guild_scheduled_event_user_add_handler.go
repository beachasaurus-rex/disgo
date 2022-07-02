package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

type gatewayHandlerGuildScheduledEventUserAdd struct{}

func (h *gatewayHandlerGuildScheduledEventUserAdd) EventType() gateway.EventType {
	return gateway.EventTypeGuildScheduledEventUserAdd
}

func (h *gatewayHandlerGuildScheduledEventUserAdd) New() any {
	return &gateway.EventGuildScheduledEventUser{}
}

func (h *gatewayHandlerGuildScheduledEventUserAdd) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	payload := *v.(*gateway.EventGuildScheduledEventUser)

	client.EventManager().DispatchEvent(&events.GuildScheduledEventUserAdd{
		GenericGuildScheduledEventUser: &events.GenericGuildScheduledEventUser{
			GenericEvent:          events.NewGenericEvent(client, sequenceNumber, shardID),
			GuildScheduledEventID: payload.GuildScheduledEventID,
			UserID:                payload.UserID,
			GuildID:               payload.GuildID,
		},
	})
}
