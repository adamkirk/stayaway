package pricing_test

import (
	"testing"

	"github.com/adamkirk-stayaway/organisations/internal/domain/pricing"
	"github.com/stretchr/testify/assert"
)

func TestTaxCalculations(t *testing.T) {
	tests := []struct{
		name string
		in int64
		tax float64
		out string
		outPretax string
	}{
		{
			name: "100 with 20% tax",
			in: 10000,
			tax: float64(20),
			out: "120.00",
			outPretax: "100.00",
		},

		{
			name: "100 with 20.5% tax",
			in: 10000,
			tax: float64(20.5),
			out: "120.50",
			outPretax: "100.00",
		},

		{
			name: "100.10 with 20% tax",
			in: 10010,
			tax: float64(20),
			out: "120.12",
			outPretax: "100.10",
		},

		{
			name: "1.10 with 20% tax",
			in: 110,
			tax: float64(20),
			out: "1.32",
			outPretax: "1.10",
		},

		{
			name: "0.30 with 20% tax",
			in: 30,
			tax: float64(20),
			out: "0.36",
			outPretax: "0.30",
		},

		{
			name: "746.35 with 20% tax",
			in: 74635,
			tax: float64(20),
			out: "895.62",
			outPretax: "746.35",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			p := pricing.NewPrice(test.in, test.tax)

			assert.Equal(tt, test.out, p.WithTaxString())
			assert.Equal(tt, test.outPretax, p.PreTaxString())
		})
	}
}


func TestAddPreTax(t *testing.T) {
	tests := []struct{
		name string
		subject pricing.Price
		toAdd pricing.Price
		expect string
	}{
		{
			name: "add two prices without tax",
			subject: pricing.NewPrice(10000, 20),
			toAdd: pricing.NewPrice(10000, 20),
			expect: "200.00",
		},
		{
			name: "add two prices without tax",
			subject: pricing.NewPrice(10010, 20),
			toAdd: pricing.NewPrice(30, 20),
			expect: "100.40",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			new := test.subject.AddPreTax(test.toAdd)

			assert.Equal(tt, test.expect, new.PreTaxString())
			assert.Equal(tt, test.subject.TaxRatePercentage(), new.TaxRatePercentage())
			
		})
	}
}

func TestAddPostTax(t *testing.T) {
	tests := []struct{
		name string
		subject pricing.Price
		toAdd pricing.Price
		expect string
	}{
		{
			name: "add two prices with tax",
			subject: pricing.NewPrice(10000, 20),
			toAdd: pricing.NewPrice(10000, 20),
			expect: "240.00",
		},
		{
			name: "add two prices with tax",
			subject: pricing.NewPrice(10010, 20),
			toAdd: pricing.NewPrice(30, 20),
			expect: "120.48",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			new := test.subject.AddPostTax(test.toAdd)

			assert.Equal(tt, test.expect, new.PreTaxString())
			assert.Equal(tt, test.subject.TaxRatePercentage(), new.TaxRatePercentage())
			
		})
	}
}

