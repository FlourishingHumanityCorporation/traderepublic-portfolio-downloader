package transaction

import (
	"errors"
	"log/slog"

	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/bus"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/pkg/traderepublic"
)

type Handler struct {
	resolver      *TypeResolver
	mapperFactory *DataMapperFactory
	eventBus      *bus.EventBus
}

func NewHandler(resolver *TypeResolver, mapperFactory *DataMapperFactory, eventBus *bus.EventBus) *Handler {
	return &Handler{
		resolver:      resolver,
		mapperFactory: mapperFactory,
		eventBus:      eventBus,
	}
}

func (h *Handler) Handle(event bus.Event) {
	var details traderepublic.TimelineDetailsJson

	// Unmarshal the JSON data from the event into a TimelineDetailsJson struct
	err := details.UnmarshalJSON(event.Data.([]byte))
	if err != nil {
		slog.Error("failed to unmarshal timeline detail", "id", event.ID, "error", err)

		return
	}

	model := Model{}

	err = h.resolver.SetType(details, &model)
	if err != nil {
		if errors.Is(err, ErrIgnoredTransactionReceived) {
			slog.Warn("ignored transaction received", "id", event.ID)

			return
		}

		if errors.Is(err, ErrCancelledTransactionReceived) {
			slog.Warn("cancelled transaction received", "id", event.ID)

			return
		}

		slog.Error("failed to resolve type", "id", event.ID, "err", err)

		return
	}

	mapper := h.mapperFactory.Make(details, &model)

	err = mapper.Map()
	if err != nil {
		slog.Error("failed to map", "id", event.ID, "err", err)

		return
	}

	h.eventBus.Publish(bus.NewEvent(bus.TopicModelReady, model.ID, model))
}
