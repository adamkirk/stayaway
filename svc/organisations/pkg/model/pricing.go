package model

import "math/big"

type Price struct {
	preTaxValue *big.Rat
	taxRatePercentage float64
}

func (p Price) TaxRatePercentage() float64 {
	return p.taxRatePercentage
}

func (p Price) PreTax() *big.Rat {
	copy := *p.preTaxValue

	return &copy
}

func (p Price) withTax() *big.Rat {
	m := big.NewFloat(100 + p.taxRatePercentage)
	mRat, _ := m.Rat(big.NewRat(1, 100))
	
	mRat.Mul(mRat, big.NewRat(1, 100))

	copy := p.PreTax()

	return copy.Mul(copy, mRat)
}

func (p Price) WithTaxString() string {
	return p.withTax().FloatString(2)
}

func (p Price) PreTaxString() string {
	return p.preTaxValue.FloatString(2)
}

func (p Price) AddPreTax(other Price) Price {
	copy := p.PreTax()
	otherCopy := other.PreTax()

	return Price{
		preTaxValue: copy.Add(copy, otherCopy),
		taxRatePercentage: p.taxRatePercentage,
	}
}

func (p Price) AddPostTax(other Price) Price {
	copy := p.withTax()
	otherCopy := other.withTax()

	return Price{
		preTaxValue: copy.Add(copy, otherCopy),
		taxRatePercentage: p.taxRatePercentage,
	}
}

// NewPrice builds a new price using the lowest possible units of currency as an 
// integer
// TODO: deal with currencies that don't use 100 units
func NewPrice(val int64, tax float64) Price {
	return Price{
		preTaxValue: big.NewRat(val, 100),
		taxRatePercentage: tax,
	}
}