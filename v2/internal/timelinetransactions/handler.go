package timelinetransactions

import (
	"context"
	"log/slog"
	"strconv"
	"sync"

	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/bus"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/message"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/pkg/traderepublic"
)

type Handler struct {
	eventBus  *bus.EventBus
	msgClient message.ClientInterface
}

func NewHandler(eventBus *bus.EventBus, msgClient message.ClientInterface) *Handler {
	return &Handler{
		eventBus:  eventBus,
		msgClient: msgClient,
	}
}

func (h *Handler) HandleFetch(event bus.Event) {
	ctx := context.Background()

	ch, err := h.msgClient.SubscribeToTimelineTransactions(ctx, "")
	if err != nil {
		slog.Error("failed to subscribe to timeline transactions", "error", err)
	}

	data := <-ch
	counter := int64(1)

	h.eventBus.Publish(bus.NewEvent(
		bus.TopicTimelineTransactionsReceived,
		strconv.FormatInt(counter, 10),
		data,
	))

	var response traderepublic.TimelineTransactionsJson

	err = response.UnmarshalJSON(data)
	if err != nil {
		slog.Error("failed to unmarshal timeline transactions", "error", err)
	}

	var mu sync.Mutex

	go func() {
		for response.Cursors.After != nil {
			ch, err = h.msgClient.SubscribeToTimelineTransactions(ctx, *response.Cursors.After)
			if err != nil {
				slog.Error("error subscribing to timeline transactions", "error", err)

				return
			}

			data = <-ch

			mu.Lock()

			counter++
			c := counter

			mu.Unlock()

			h.eventBus.Publish(bus.NewEvent(
				bus.TopicTimelineTransactionsReceived,
				strconv.FormatInt(c, 10),
				data,
			))

			err = response.UnmarshalJSON(data)
			if err != nil {
				slog.Error("error subscribing to timeline transactions", "error", err)

				return
			}
		}
	}()
}

func (h *Handler) HandleReceived(event bus.Event) {
	var transactions traderepublic.TimelineTransactionsJson

	err := transactions.UnmarshalJSON(event.Data.([]byte))
	if err != nil {
		slog.Error("invalid event data type", "expected", "traderepublic.TimelineTransactionsSchemaJson", "got", event.Data)
	}

	for _, transaction := range transactions.Items {
		err := h.msgClient.SubscribeToTimelineDetailV2(context.Background(), transaction.Id)
		if err != nil {
			slog.Error("failed to subscribe to timeline detail", "error", err, "transaction_id", transaction.Id)

			continue
		}
	}
}
