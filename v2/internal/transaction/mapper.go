package transaction

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/pkg/traderepublic"
	gocache "github.com/patrickmn/go-cache"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	ErrTransactionWithoutTypeReceived = errors.New("transaction without type received")
	ErrUnsupportedTransactionReceived = errors.New("unsupported transaction received")
)

// DataMapperFactory creates DataMapper instances with shared cache
type DataMapperFactory struct {
	cache *gocache.Cache
}

func NewDataMapperFactory(cache *gocache.Cache) *DataMapperFactory {
	return &DataMapperFactory{
		cache: cache,
	}
}

// Make creates a new DataMapper for the given transaction details and model
func (f *DataMapperFactory) Make(details traderepublic.TimelineDetailsJson, model *Model) *DataMapper {
	return NewDataMapper(details, model, f.cache)
}

// DataMapper maps Trade Republic transaction details to internal transaction model
type DataMapper struct {
	details  traderepublic.TimelineDetailsJson
	model    *Model
	cache    *gocache.Cache
	header   traderepublic.HeaderSection
	overview traderepublic.TableSection
}

func NewDataMapper(details traderepublic.TimelineDetailsJson, model *Model, cache *gocache.Cache) *DataMapper {
	return &DataMapper{
		details: details,
		model:   model,
		cache:   cache,
	}
}

// Map transforms Trade Republic transaction data into our internal model
func (m *DataMapper) Map() error {
	// Validate transaction has a recognized type
	if m.model.Type == TypeUnknown {
		return fmt.Errorf("%w: %s", ErrTransactionWithoutTypeReceived, m.details.Id)
	}

	// Only process portfolio-related transactions
	if !slices.Contains(PortfolioTypes, m.model.Type) {
		return fmt.Errorf("%w: %s", ErrUnsupportedTransactionReceived, m.details.Id)
	}

	// Extract required sections from transaction details
	header, err := m.details.SectionHeader()
	if err != nil {
		return fmt.Errorf("failed to find header section: %w", err)
	}

	m.header = header

	overview, err := m.details.FindSection(traderepublic.SectionOverview)
	if err != nil {
		return fmt.Errorf("failed to find overview section: %w", err)
	}

	m.overview = overview

	// Map basic transaction fields
	m.MapID()
	m.MapStatus()

	err = m.mapISIN()
	if err != nil {
		return err
	}

	err = m.mapTimestamp()
	if err != nil {
		return err
	}

	err = m.mapShares()
	if err != nil {
		return err
	}

	err = m.mapSharePrice()
	if err != nil {
		return err
	}

	err = m.mapFee()
	if err != nil {
		return err
	}

	err = m.mapTotal()
	if err != nil {
		return err
	}

	err = m.mapAsset()
	if err != nil {
		return err
	}

	err = m.mapProfit()
	if err != nil {
		return err
	}

	err = m.mapGain()
	if err != nil {
		return err
	}

	return nil
}

// MapID extracts transaction ID from details
func (m *DataMapper) MapID() {
	m.model.ID = string(m.details.Id)
}

// MapStatus extracts transaction status from header
func (m *DataMapper) MapStatus() {
	m.model.Status = cases.Title(language.English).String(string(m.header.Data.Status))
}

// mapTimestamp parses and sets transaction timestamp
func (m *DataMapper) mapTimestamp() error {
	timestampStr := m.header.Data.Timestamp

	timestamp, err := ParseTimestamp(timestampStr)
	if err != nil {
		return fmt.Errorf("failed to parse timestamp: %w", err)
	}

	m.model.Timestamp = CSVDateTime{Time: timestamp}

	return nil
}

// mapISIN extracts ISIN from either action payload or icon data
func (m *DataMapper) mapISIN() error {
	// Interest payments don't have ISIN
	if m.model.Type == TypeInterestPayment {
		return nil
	}

	// Try to get ISIN from action payload first
	if m.header.Action == nil {
		// Fallback: extract ISIN from icon data
		isin, err := ExtractInstrumentISINFromIcon(m.header.Data.Icon)
		if err != nil {
			return fmt.Errorf("failed to extract ISIN from icon: %w", err)
		}

		m.model.ISIN = isin

		return nil
	}

	m.model.ISIN = m.header.Action.Payload

	return nil
}

// mapAsset fetches and sets asset name and type using ISIN
func (m *DataMapper) mapAsset() error {
	if m.model.ISIN == "" {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	instr, err := m.getInstrument(ctx, m.model.ISIN)
	if err != nil {
		return fmt.Errorf("failed to get instrument: %w", err)
	}

	m.model.AssetName = *instr.ShortName
	m.model.AssetType = cases.Title(language.English).String(string(instr.TypeId))

	return nil
}

// mapShares extracts share quantity, handling different transaction formats
func (m *DataMapper) mapShares() error {
	var str string

	// Handle different transaction formats based on type
	switch m.model.Type {
	case TypeSavingsPlanPre202502, TypeBuyOrderPre202502, TypeSellOrderPre202502, TypeLimitSellPre202502, TypeDividendsIncome:
		// Pre-2025 format: shares in transaction section
		trnSection, err := m.details.FindSection(traderepublic.SectionTransaction)
		if err != nil {
			return fmt.Errorf("failed to find transaction section: %w", err)
		}

		shares, err := trnSection.FindData(traderepublic.DataShares)
		if err != nil {
			return fmt.Errorf("failed to find shares data: %w", err)
		}

		if shares.Detail.DisplayValue != nil {
			str = shares.Detail.DisplayValue.Text

			break
		}

		str = shares.Detail.Text
	case TypeSavingsPlan, TypeLimitSell, TypeSellOrder, TypeBuyOrder:
		// Current format: shares in overview section
		trn, err := m.overview.FindData(traderepublic.DataTransaction)
		if err != nil {
			return fmt.Errorf("failed to find transaction data: %w", err)
		}

		if trn.Detail.DisplayValue.Prefix == nil {
			return fmt.Errorf("failed to read trn.Detail.DisplayValue.Prefix: nil '%s'", m.model.ID)
		}

		str = *trn.Detail.DisplayValue.Prefix
	case TypeInterestPayment:
		// Interest payments don't have shares
		return nil
	}

	shares, err := ParseStringToFloat64(str)
	if err != nil {
		return fmt.Errorf("failed to parse float from shares: %w", err)
	}

	// Make shares negative for sell transactions
	if slices.Contains(GainTypes, m.model.Type) {
		m.model.Shares = -shares

		return nil
	}

	m.model.Shares = shares

	return nil
}

// mapSharePrice extracts price per share, handling different transaction formats
func (m *DataMapper) mapSharePrice() error {
	var str string

	switch m.model.Type {
	case TypeDividendsIncome:
		// For dividends, use dividend per share
		trnSection, err := m.details.FindSection(traderepublic.SectionTransaction)
		if err != nil {
			return fmt.Errorf("failed to find transaction section: %w", err)
		}

		dividendsPerShare, err := trnSection.FindData(traderepublic.DataDividendPerShare)
		if err != nil {
			return fmt.Errorf("failed to find shares data: %w", err)
		}

		if dividendsPerShare.Detail.DisplayValue != nil {
			str = dividendsPerShare.Detail.DisplayValue.Text

			break
		}

		str = dividendsPerShare.Detail.Text
	case TypeSavingsPlanPre202502, TypeBuyOrderPre202502, TypeSellOrderPre202502, TypeLimitSellPre202502:
		// Pre-2025 format: share price in transaction section
		trnSection, err := m.details.FindSection(traderepublic.SectionTransaction)
		if err != nil {
			return fmt.Errorf("failed to find transaction section: %w", err)
		}

		sharePrice, err := trnSection.FindData(traderepublic.DataSharePrice)
		if err != nil {
			return fmt.Errorf("failed to find shares data: %w", err)
		}

		if sharePrice.Detail.DisplayValue != nil {
			str = sharePrice.Detail.DisplayValue.Text

			break
		}

		str = sharePrice.Detail.Text
	case TypeSavingsPlan, TypeLimitSell, TypeSellOrder, TypeBuyOrder:
		// Current format: share price in overview section
		trn, err := m.overview.FindData(traderepublic.DataTransaction)
		if err != nil {
			return fmt.Errorf("failed to find transaction data: %w", err)
		}

		str = trn.Detail.DisplayValue.Text
	case TypeInterestPayment:
		// Interest payments don't have share price
		return nil
	}

	sharePrice, err := ParseStringToFloat64(str)
	if err != nil {
		return fmt.Errorf("failed to parse float from share price: %w", err)
	}

	m.model.SharePrice = sharePrice

	return nil
}

// mapFee extracts transaction fees (dividends are fee-free)
func (m *DataMapper) mapFee() error {
	var str string

	switch m.model.Type {
	case TypeDividendsIncome:
		// Dividends don't have fees
		return nil
	case TypeSavingsPlanPre202502, TypeBuyOrderPre202502, TypeSellOrderPre202502, TypeLimitSellPre202502:
		// Pre-2025 format: fee in transaction section
		trnSection, err := m.details.FindSection(traderepublic.SectionTransaction)
		if err != nil {
			return fmt.Errorf("failed to find transaction section: %w", err)
		}

		fee, err := trnSection.FindData(traderepublic.DataFee)
		if err != nil {
			return fmt.Errorf("failed to find fee data: %w", err)
		}

		if fee.Detail.DisplayValue != nil {
			str = fee.Detail.DisplayValue.Text

			break
		}

		str = fee.Detail.Text
	case TypeSavingsPlan, TypeLimitSell, TypeSellOrder, TypeBuyOrder:
		// Current format: fee in overview section
		fee, err := m.overview.FindData(traderepublic.DataFee)
		if err != nil {
			return fmt.Errorf("failed to find fee data: %w", err)
		}

		str = fee.Detail.DisplayValue.Text
	case TypeInterestPayment:
		// Interest payments don't have fees
		return nil
	}

	// Handle free transactions
	if str != "Free" {
		fee, err := ParseStringToFloat64(str)
		if err != nil {
			return fmt.Errorf("failed to parse float from fee: %w", err)
		}

		m.model.Fee = fee
	}

	return nil
}

// mapTotal extracts total transaction amount and sets debit/credit accordingly
func (m *DataMapper) mapTotal() error {
	var str string

	switch m.model.Type {
	case TypeSavingsPlanPre202502, TypeBuyOrderPre202502, TypeSellOrderPre202502, TypeLimitSellPre202502, TypeDividendsIncome, TypeInterestPayment:
		// Pre-2025 format: total in transaction section
		trnSection, err := m.details.FindSection(traderepublic.SectionTransaction)
		if err != nil {
			return fmt.Errorf("failed to find transaction section: %w", err)
		}

		total, err := trnSection.FindData(traderepublic.DataTotal)
		if err != nil {
			return fmt.Errorf("failed to find total data: %w", err)
		}

		if total.Detail.DisplayValue != nil {
			str = total.Detail.DisplayValue.Text

			break
		}

		str = total.Detail.Text
	case TypeSavingsPlan, TypeLimitSell, TypeSellOrder, TypeBuyOrder:
		// Current format: total in overview section
		total, err := m.overview.FindData(traderepublic.DataTotal)
		if err != nil {
			return fmt.Errorf("failed to find total data: %w", err)
		}

		str = total.Detail.DisplayValue.Text
	}

	total, err := ParseStringToFloat64(str)
	if err != nil {
		return fmt.Errorf("failed to parse float from total: %w", err)
	}

	// For credit transactions (money coming in), set credit field
	if slices.Contains(CreditTypes, m.model.Type) {
		m.model.Credit = &total

		return nil
	}

	// Default to debit (money going out)
	m.model.Debit = &total

	return nil
}

// mapProfit extracts realized profit for sell transactions
func (m *DataMapper) mapProfit() error {
	// Interest payments and dividends don't have profit
	if m.model.Type == TypeInterestPayment {
		return nil
	}

	// Skip non-gain transactions and dividends
	if m.model.Type == TypeDividendsIncome || !slices.Contains(GainTypes, m.model.Type) {
		return nil
	}

	performance, err := m.details.FindSection(traderepublic.SectionPerformance)
	if err != nil {
		return fmt.Errorf("failed to find performance section: %w", err)
	}

	profitData, err := performance.FindData(traderepublic.DataProfit)
	if err != nil {
		return fmt.Errorf("failed to find profit data: %w", err)
	}

	str := profitData.Detail.Text

	profit, err := ParseStringToFloat64(str)
	if err != nil {
		return fmt.Errorf("failed to parse float from profit: %w", err)
	}

	m.model.Profit = &profit

	return nil
}

// mapGain extracts gain/loss percentage for sell transactions or dividend amount
func (m *DataMapper) mapGain() error {
	if m.model.Type == TypeDividendsIncome || m.model.Type == TypeInterestPayment {
		// For dividends, gain equals credit amount
		m.model.Gain = m.model.Credit

		return nil
	}

	// Skip non-gain transactions
	if !slices.Contains(GainTypes, m.model.Type) {
		return nil
	}

	performance, err := m.details.FindSection(traderepublic.SectionPerformance)
	if err != nil {
		return fmt.Errorf("failed to find performance section: %w", err)
	}

	gainData, err := performance.FindData(traderepublic.DataGain)
	if err != nil {
		return fmt.Errorf("failed to find gain data: %w", err)
	}

	str := gainData.Detail.Text

	gain, err := ParseStringToFloat64(str)
	if err != nil {
		return fmt.Errorf("failed to parse float from gain: %w", err)
	}

	if *gainData.Detail.Trend == "negative" {
		negativeGain := -gain
		m.model.Gain = &negativeGain

		return nil
	}

	m.model.Gain = &gain

	return nil
}

func (m *DataMapper) getInstrument(ctx context.Context, isin string) (traderepublic.InstrumentJson, error) {
	for {
		select {
		case <-ctx.Done():
			return traderepublic.InstrumentJson{}, fmt.Errorf("context timeout: %w", ctx.Err())

		default:
			entry, found := m.cache.Get(isin)
			if found {
				instr, ok := entry.(traderepublic.InstrumentJson)
				if !ok {
					continue
				}

				return instr, nil
			}
		}
	}
}
