package transaction

import (
	"errors"
	"fmt"

	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/pkg/traderepublic"
)

const (
	TypeUnknown              Type = "unknown"                   // Unknown transaction type
	TypeIgnored              Type = "ingored"                   // Ignored transaction type
	TypeSavingsPlan          Type = "Savings plan"              // Savings plan transaction
	TypeSavingsPlanPre202502 Type = "Savings plan pre Feb 2025" // Savings plan transaction
	TypeCardPayment          Type = "Card payment"              // Card payment transaction
	TypeCardRefund           Type = "Card refund"               // Card refund transaction
	TypeBuyOrder             Type = "Buy order"                 // Buy order transaction
	TypeBuyOrderPre202502    Type = "Buy order pre Feb 2025"    // Buy order transaction
	TypeSellOrder            Type = "Sell order"                // Sell order transaction
	TypeSellOrderPre202502   Type = "Sell order pre Feb 2025"   // Sell order transaction
	TypeLimitSell            Type = "Limit sell"
	TypeLimitSellPre202502   Type = "Limit sell pre Feb 2025"
	TypeDividendsIncome      Type = "Dividends income" // Dividends income transaction
	TypeRoundUp              Type = "Round up"         // Round up transaction
	TypeSaveback             Type = "Saveback"         // Saveback transaction
	TypeDeposit              Type = "Deposit"          // Deposit transaction
	TypeWithdrawal           Type = "Withdrawal"       // Withdrawal transaction
	TypeInterestPayment      Type = "Interest payment" // Interest payment transaction
)

var (
	ErrCancelledTransactionReceived = errors.New("canceled transaction type received")
	ErrIgnoredTransactionReceived   = errors.New("ignored transaction type received")
	ErrUnknownTransactionReceived   = errors.New("unknown transaction type received")

	PortfolioTypes = []Type{
		TypeSavingsPlan,
		TypeSavingsPlanPre202502,
		TypeBuyOrder,
		TypeBuyOrderPre202502,
		TypeSellOrder,
		TypeSellOrderPre202502,
		TypeLimitSell,
		TypeLimitSellPre202502,
		TypeDividendsIncome,
		TypeRoundUp,
		TypeSaveback,
		TypeDeposit,
		TypeWithdrawal,
		TypeInterestPayment,
	}

	GainTypes = []Type{
		TypeSellOrder,
		TypeSellOrderPre202502,
		TypeLimitSell,
		TypeLimitSellPre202502,
		TypeDividendsIncome,
		TypeInterestPayment,
	}

	CreditTypes = []Type{
		TypeSellOrder,
		TypeSellOrderPre202502,
		TypeLimitSell,
		TypeLimitSellPre202502,
		TypeDividendsIncome,
		TypeInterestPayment,
		TypeDeposit,
	}
)

// Type represents the type of a transaction.
type Type string

// TypeResolver resolves the type of a transaction based on its details.
type TypeResolver struct {
}

// NewTypeResolver creates a new instance of TypeResolver.
func NewTypeResolver() *TypeResolver {
	return &TypeResolver{}
}

// Resolve determines the type of a transaction from its details.
func (r *TypeResolver) SetType(details traderepublic.TimelineDetailsJson, model *Model) error {
	header, err := details.SectionHeader()
	if err == nil {
		if header.Data.Status == "canceled" {
			return ErrCancelledTransactionReceived
		}
	}

	overview, err := details.FindSection(traderepublic.SectionOverview)
	if err != nil {
		return fmt.Errorf("failed to find overview section: %w", err)
	}

	// Check for ignored transactions
	_, err = overview.FindData(traderepublic.DataCardVerification)
	if err == nil {
		return fmt.Errorf("%w: %s", ErrIgnoredTransactionReceived, details.Id)
	}

	// Check for card payment transactions
	_, err = overview.FindData(traderepublic.DataCardPayment)
	if err == nil {
		model.Type = TypeCardPayment

		return nil
	}

	// Check for card refund transactions
	_, err = overview.FindData(traderepublic.DataCardRefund)
	if err == nil {
		model.Type = TypeCardRefund

		return nil
	}

	// // Check for deposit transactions
	// _, err = details.FindSection(traderepublic.SectionSender)
	// if err == nil {
	// 	model.Type = TypeDeposit

	// 	return
	// }
	// _, err = overview.FindData(traderepublic.DataFrom)
	// if err == nil {
	// 	model.Type = TypeDeposit

	// 	return
	// }

	// // Check for withdrawal transactions
	// _, err = overview.FindData(traderepublic.DataTo)
	// if err == nil {
	// 	model.Type = TypeWithdrawal

	// 	return
	// }

	// payment, err := overview.FindData(traderepublic.DataPayment)
	// if err == nil {
	// 	if payment.Detail.Text == "Direct Debit" {
	// 		model.Type = TypeDeposit

	// 		return
	// 	}
	// }

	// Check for savings plan transactions
	_, err = overview.FindData(traderepublic.DataSavingsPlan)
	if err == nil {
		model.Type = TypeSavingsPlan

		return nil
	}

	event, err := overview.FindData(traderepublic.DataEvent)
	if err == nil {
		switch event.Detail.Text {
		case "Income", "Cash dividend":
			model.Type = TypeDividendsIncome

			return nil
		case "Tax Settlement":
			return fmt.Errorf("%w: %s", ErrIgnoredTransactionReceived, details.Id)
		}
	}

	// Check for buy order transactions
	orderType, err := overview.FindData(traderepublic.DataOrderType)
	if err == nil {
		switch orderType.Detail.Text {
		case "Savings plan":
			model.Type = TypeSavingsPlanPre202502

			return nil
		case "Buy":
			model.Type = TypeBuyOrderPre202502

			return nil
		case "Sell":
			model.Type = TypeSellOrderPre202502

			return nil
		case "Limit Sell":
			model.Type = TypeLimitSellPre202502

			return nil
		}
	}

	// Check for interest payment transactions
	_, err = overview.FindData(traderepublic.DataAverageBalance)
	if err == nil {
		model.Type = TypeInterestPayment

		return nil
	}
	steps, err := details.SectionSteps()
	if err == nil {
		_, err = steps.FindStep(traderepublic.StepInterestPayment)
		if err == nil {
			model.Type = TypeInterestPayment

			return nil
		}
	}

	// // Check for saveback transactions
	// _, err = overview.FindData(traderepublic.DataSaveback)
	// if err == nil {
	// 	model.Type = TypeSaveback

	// 	return
	// }

	// // Check for round up transactions
	// _, err = overview.FindData(traderepublic.DataRoundUp)
	// if err == nil {
	// 	model.Type = TypeRoundUp

	// 	return
	// }

	// Check for limit sell transactions
	_, err = overview.FindData(traderepublic.DataLimitSell)
	if err == nil {
		model.Type = TypeLimitSell

		return nil
	}

	// Check for sell transactions
	_, err = overview.FindData(traderepublic.DataSell)
	if err == nil {
		model.Type = TypeSellOrder

		return nil
	}

	// Check for buy transactions
	_, err = overview.FindData(traderepublic.DataBuy)
	if err == nil {
		model.Type = TypeBuyOrder

		return nil
	}

	return fmt.Errorf("%w: %s", ErrUnknownTransactionReceived, details.Id)
}
