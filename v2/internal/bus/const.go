package bus

const (
	// TopicAuthenticated is published when user authentication is successful
	TopicAuthenticated = "authenticated"

	// TopicAuthSessionRefreshed is published when session token is refreshed
	TopicAuthSessionRefreshed = "auth_session_refreshed"

	// TopicAuthFailed is published when authentication fails
	TopicAuthFailed = "auth_failed"

	// TopicTimelineTransactionsReceived is published when timeline transactions data is received
	TopicTimelineTransactionsReceived = "timeline_transactions_received"

	// TopicTimelineDetailsV2Received is published when timeline detail v2 data is received
	TopicTimelineDetailsV2Received = "timeline_detail_v2_received"

	// TopicInstrumentFetch is published to request instrument data fetch
	TopicInstrumentFetch = "instrument_fetch"

	// TopicInstrumentReceived is published when instrument data is received
	TopicInstrumentReceived = "instrument_received"

	// TopicModelCreated is published when a data model is created
	TopicModelCreated = "model_created"

	//TopicWorkersDone is published when all workers have completed their tasks
	TopicWorkersDone = "workers_done"
)
