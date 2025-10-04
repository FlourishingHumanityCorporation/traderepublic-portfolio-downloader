package transaction_test

import (
	"fmt"
	"testing"

	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/transaction"
	"github.com/stretchr/testify/assert"
)

func TestItExtractsInstrumentISINFromIcon(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		input    string
		expected string
	}{
		{input: "logos/FR0003500008/v2", expected: "FR0003500008"},
		{input: "logos/DE000A0F5UF5/v2", expected: "DE000A0F5UF5"},
		{input: "logos/XF000DOT0011/v2", expected: "XF000DOT0011"},
		{input: "logos/IE00BK1PV551/v2", expected: "IE00BK1PV551"},
		{input: "logos/US6701002056/v2", expected: "US6701002056"},
		{input: "logos/XF000AVAX016/v2", expected: "XF000AVAX016"},
	}

	for i, testCase := range testCases {
		actual, err := transaction.ExtractInstrumentISINFromIcon(testCase.input)

		assert.NoError(t, err, fmt.Sprintf("case %d", i))
		assert.Equal(t, testCase.expected, actual, fmt.Sprintf("case %d", i))
	}
}

func TestItReturnsErrorOnIconContainsNoISIN(t *testing.T) {
	t.Parallel()

	testCases := []string{
		"logos/timeline_document/v2",
		"logos/timeline_interest_new/v2",
		"logos/timeline_plus_circle/v2",
		"logos/timeline_minus_circle/v2",
	}

	for i, testCase := range testCases {
		_, err := transaction.ExtractInstrumentISINFromIcon(testCase)

		assert.Error(t, err, fmt.Sprintf("case %d", i))
	}
}

func TestParseFloatFromResponse(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		input    string
		expected float64
	}{
		{input: "€1,000.00", expected: 1000},
		{input: "€5,462.13", expected: 5462.13},
		{input: "€5,000,462.13", expected: 5000462.13},
		{input: "€4.12", expected: 4.12},
		{input: "€217.30", expected: 217.30},
		{input: "€94.39", expected: 94.39},
		{input: "139.003939", expected: 139.003939},
		{input: "139.003939 × ", expected: 139.003939},
		{input: "200 × ", expected: 200},
	}

	for i, testCase := range testCases {
		actual, err := transaction.ParseStringToFloat64(testCase.input)

		assert.NoError(t, err, fmt.Sprintf("case %d", i))
		assert.Equal(t, testCase.expected, actual, fmt.Sprintf("case %d", i))
	}
}
