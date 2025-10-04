package transaction

type Model struct {
	ID             string
	Status         string
	Timestamp      CSVDateTime
	Type           Type
	AssetType      string `csv:"Asset type"`
	AssetName      string `csv:"Asset name"`
	ISIN           string
	Shares         float64
	SharePrice     float64 `csv:"Share price"`
	Profit         float64
	Gain           float64
	Fee            float64
	Debit          float64
	Credit         float64
	TaxAmount      float64 `csv:"Tax amount"`
	InvestedAmount float64 `csv:"-"`
	Documents      []string
}
