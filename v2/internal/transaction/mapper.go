package transaction

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/pkg/traderepublic"
	gocache "github.com/patrickmn/go-cache"
)

var (
	ErrTransactionWithoutTypeReceived = errors.New("transaction without type received")
	ErrUnsupportedTransactionReceived = errors.New("unsupported transaction received")
)

type DataMapperFactory struct {
	cache *gocache.Cache
}

func NewDataMapperFactory(cache *gocache.Cache) *DataMapperFactory {
	return &DataMapperFactory{
		cache: cache,
	}
}

func (f *DataMapperFactory) Make(details traderepublic.TimelineDetailsJson, model *Model) *DataMapper {
	return NewDataMapper(details, model, f.cache)
}

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

func (m *DataMapper) Map() error {
	if m.model.Type == TypeUnknown {
		return fmt.Errorf("%w: %s", ErrTransactionWithoutTypeReceived, m.details.Id)
	}

	if !slices.Contains(PortfolioTypes, m.model.Type) {
		return fmt.Errorf("%w: %s", ErrUnsupportedTransactionReceived, m.details.Id)
	}

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

func (m *DataMapper) MapID() {
	m.model.ID = string(m.details.Id)
}

func (m *DataMapper) MapStatus() {
	m.model.Status = string(m.header.Data.Status)
}

func (m *DataMapper) mapTimestamp() error {
	timestampStr := m.header.Data.Timestamp

	timestamp, err := ParseTimestamp(timestampStr)
	if err != nil {
		return fmt.Errorf("failed to parse timestamp: %w", err)
	}

	m.model.Timestamp = CSVDateTime{Time: timestamp}

	return nil
}

func (m *DataMapper) mapISIN() error {
	if m.header.Action == nil {
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
	m.model.AssetType = string(instr.TypeId)

	return nil
}

func (m *DataMapper) mapShares() error {
	var str string

	switch m.model.Type {
	case TypeSavingsPlanPre202502, TypeBuyOrderPre202502, TypeSellOrderPre202502, TypeLimitSellPre202502, TypeDividendsIncome:
		trnSection, err := m.details.FindSection(traderepublic.SectionTransaction)
		if err != nil {
			return fmt.Errorf("failed to find transaction section: %w", err)
		}

		shares, err := trnSection.FindData(traderepublic.DataShares)
		if err != nil {
			return fmt.Errorf("failed to find shares data: %w", err)
		}

		str = shares.Detail.Text
	case TypeSavingsPlan, TypeLimitSell, TypeSellOrder, TypeBuyOrder:
		trn, err := m.overview.FindData(traderepublic.DataTransaction)
		if err != nil {
			return fmt.Errorf("failed to find transaction data: %w", err)
		}

		str = *trn.Detail.DisplayValue.Prefix
	}

	shares, err := ParseFloatFromResponse(str)
	if err != nil {
		return fmt.Errorf("failed to parse float from shares: %w", err)
	}

	m.model.Shares = shares

	if slices.Contains(GainTypes, m.model.Type) {
		m.model.Shares = -shares
	}

	return nil
}

func (m *DataMapper) mapSharePrice() error {
	var str string

	switch m.model.Type {
	case TypeDividendsIncome:
		trnSection, err := m.details.FindSection(traderepublic.SectionTransaction)
		if err != nil {
			return fmt.Errorf("failed to find transaction section: %w", err)
		}

		dividendsPerShare, err := trnSection.FindData(traderepublic.DataDividendPerShare)
		if err != nil {
			return fmt.Errorf("failed to find shares data: %w", err)
		}

		str = dividendsPerShare.Detail.Text
	case TypeSavingsPlanPre202502, TypeBuyOrderPre202502, TypeSellOrderPre202502, TypeLimitSellPre202502:
		trnSection, err := m.details.FindSection(traderepublic.SectionTransaction)
		if err != nil {
			return fmt.Errorf("failed to find transaction section: %w", err)
		}

		sharePrice, err := trnSection.FindData(traderepublic.DataSharePrice)
		if err != nil {
			return fmt.Errorf("failed to find shares data: %w", err)
		}

		str = sharePrice.Detail.Text
	case TypeSavingsPlan, TypeLimitSell, TypeSellOrder, TypeBuyOrder:
		trn, err := m.overview.FindData(traderepublic.DataTransaction)
		if err != nil {
			return fmt.Errorf("failed to find transaction data: %w", err)
		}

		str = trn.Detail.DisplayValue.Text
	}

	sharePrice, err := ParseFloatFromResponse(str)
	if err != nil {
		return fmt.Errorf("failed to parse float from share price: %w", err)
	}

	m.model.SharePrice = sharePrice

	return nil
}

func (m *DataMapper) mapFee() error {
	var str string

	switch m.model.Type {
	case TypeDividendsIncome:
		return nil
	case TypeSavingsPlanPre202502, TypeBuyOrderPre202502, TypeSellOrderPre202502, TypeLimitSellPre202502:
		trnSection, err := m.details.FindSection(traderepublic.SectionTransaction)
		if err != nil {
			return fmt.Errorf("failed to find transaction section: %w", err)
		}

		fee, err := trnSection.FindData(traderepublic.DataFee)
		if err != nil {
			return fmt.Errorf("failed to find fee data: %w", err)
		}

		str = fee.Detail.Text
	case TypeSavingsPlan, TypeLimitSell, TypeSellOrder, TypeBuyOrder:
		fee, err := m.overview.FindData(traderepublic.DataFee)
		if err != nil {
			return fmt.Errorf("failed to find fee data: %w", err)
		}

		str = fee.Detail.Text
	}

	if str != "Free" {
		fee, err := ParseFloatFromResponse(str)
		if err != nil {
			return fmt.Errorf("failed to parse float from fee: %w", err)
		}

		m.model.Fee = fee
	}

	return nil
}

func (m *DataMapper) mapTotal() error {
	var str string

	switch m.model.Type {
	case TypeSavingsPlanPre202502, TypeBuyOrderPre202502, TypeSellOrderPre202502, TypeLimitSellPre202502, TypeDividendsIncome:
		trnSection, err := m.details.FindSection(traderepublic.SectionTransaction)
		if err != nil {
			return fmt.Errorf("failed to find transaction section: %w", err)
		}

		total, err := trnSection.FindData(traderepublic.DataTotal)
		if err != nil {
			return fmt.Errorf("failed to find total data: %w", err)
		}

		str = total.Detail.Text
	case TypeSavingsPlan, TypeLimitSell, TypeSellOrder, TypeBuyOrder:
		total, err := m.overview.FindData(traderepublic.DataTotal)
		if err != nil {
			return fmt.Errorf("failed to find total data: %w", err)
		}

		str = total.Detail.Text
	}

	total, err := ParseFloatFromResponse(str)
	if err != nil {
		return fmt.Errorf("failed to parse float from total: %w", err)
	}

	m.model.Debit = total

	if slices.Contains(CreditTypes, m.model.Type) {
		m.model.Credit = total
	}

	return nil
}

func (m *DataMapper) mapProfit() error {
	if m.model.Type == TypeDividendsIncome {
		return nil
	}

	if !slices.Contains(GainTypes, m.model.Type) {
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

	profit, err := ParseFloatFromResponse(str)
	if err != nil {
		return fmt.Errorf("failed to parse float from profit: %w", err)
	}

	m.model.Profit = profit

	return nil
}

func (m *DataMapper) mapGain() error {
	if m.model.Type == TypeDividendsIncome {
		m.model.Gain = m.model.Credit

		return nil
	}

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

	gain, err := ParseFloatFromResponse(str)
	if err != nil {
		return fmt.Errorf("failed to parse float from gain: %w", err)
	}

	m.model.Gain = gain

	if *gainData.Detail.Trend == "negative" {
		m.model.Gain = -gain
	}

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
