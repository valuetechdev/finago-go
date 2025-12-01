package product

import (
	"encoding/xml"
	"fmt"
)

// NOTE: Custom marshaling for removing scientific notation from float64 Price
func (p Product) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	type Alias Product

	modified := struct {
		Alias

		CashPriceIncTax         *string `xml:"CashPriceIncTax,omitempty"`
		OriginalCashPriceIncTax float64 `xml:"-"`

		WebPrice         *string `xml:"WebPrice,omitempty"`
		OriginalWebPrice float64 `xml:"-"`

		TaxRate         *string `xml:"TaxRate,omitempty"`
		OriginalTaxRate float64 `xml:"-"`

		Stock         *string `xml:"Stock,omitempty"`
		OriginalStock float64 `xml:"-"`

		InPrice         *string `xml:"InPrice,omitempty"`
		OriginalInPrice float64 `xml:"-"`

		Cost         *string `xml:"Cost,omitempty"`
		OriginalCost float64 `xml:"-"`

		Price         *string `xml:"Price,omitempty"`
		OriginalPrice float64 `xml:"-"`

		Weight         *string `xml:"Weight,omitempty"`
		OriginalWeight float64 `xml:"-"`

		MinimumStock         *string `xml:"MinimumStock,omitempty"`
		OriginalMinimumStock float64 `xml:"-"`

		OrderProposal         *string `xml:"OrderProposal,omitempty"`
		OriginalOrderProposal float64 `xml:"-"`
	}{
		Alias:                   Alias(p),
		OriginalCashPriceIncTax: p.CashPriceIncTax,
		OriginalWebPrice:        p.WebPrice,
		OriginalTaxRate:         p.TaxRate,
		OriginalStock:           p.Stock,
		OriginalInPrice:         p.InPrice,
		OriginalCost:            p.Cost,
		OriginalPrice:           p.Price,
		OriginalWeight:          p.Weight,
		OriginalMinimumStock:    p.MinimumStock,
		OriginalOrderProposal:   p.OrderProposal,
	}

	// Only include CashPriceIncTax if it's not zero
	if p.CashPriceIncTax != 0 {
		formattedPrice := fmt.Sprintf("%.2f", p.CashPriceIncTax)
		modified.CashPriceIncTax = &formattedPrice
	}

	// Only include WebPrice if it's not zero
	if p.WebPrice != 0 {
		formattedPrice := fmt.Sprintf("%.2f", p.WebPrice)
		modified.WebPrice = &formattedPrice
	}

	// Only include TaxRate if it's not zero
	if p.TaxRate != 0 {
		formattedPrice := fmt.Sprintf("%.2f", p.TaxRate)
		modified.TaxRate = &formattedPrice
	}

	// Only include Stock if it's not zero
	if p.Stock != 0 {
		formattedPrice := fmt.Sprintf("%.2f", p.Stock)
		modified.Stock = &formattedPrice
	}

	// Only include InPrice if it's not zero
	if p.InPrice != 0 {
		formattedPrice := fmt.Sprintf("%.2f", p.InPrice)
		modified.InPrice = &formattedPrice
	}

	// Only include Cost if it's not zero
	if p.Cost != 0 {
		formattedPrice := fmt.Sprintf("%.2f", p.Cost)
		modified.Cost = &formattedPrice
	}

	// Only include Price if it's not zero
	if p.Price != 0 {
		formattedPrice := fmt.Sprintf("%.2f", p.Price)
		modified.Price = &formattedPrice
	}

	// Only include Weight if it's not zero
	if p.Weight != 0 {
		formattedPrice := fmt.Sprintf("%.2f", p.Weight)
		modified.Weight = &formattedPrice
	}

	// Only include MinimumStock if it's not zero
	if p.MinimumStock != 0 {
		formattedPrice := fmt.Sprintf("%.2f", p.MinimumStock)
		modified.MinimumStock = &formattedPrice
	}

	// Only include OrderProposal if it's not zero
	if p.OrderProposal != 0 {
		formattedPrice := fmt.Sprintf("%.2f", p.OrderProposal)
		modified.OrderProposal = &formattedPrice
	}

	return e.EncodeElement(modified, start)
}
