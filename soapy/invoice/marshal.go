package invoice

import (
	"encoding/xml"
	"fmt"
)

// NOTE: Custom marshaling for removing scientific notation from float64 Price
func (i InvoiceRow) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	type Alias InvoiceRow

	modified := struct {
		Alias
		Price             *string `xml:"Price,omitempty"`
		Quantity          *string `xml:"Quantity,omitempty"`
		VatRate           *string `xml:"VatRate,omitempty"`
		QuantityDelivered *string `xml:"QuantityDelivered,omitempty"`
		QuantityOrdered   *string `xml:"QuantityOrdered,omitempty"`
		QuantityRest      *string `xml:"QuantityRest,omitempty"`
		Cost              *string `xml:"Cost,omitempty"`
		InPrice           *string `xml:"InPrice,omitempty"`
		// Hide originals
		OriginalPrice             float64 `xml:"-"`
		OriginalQuantity          float64 `xml:"-"`
		OriginalVatRate           float64 `xml:"-"`
		OriginalQuantityDelivered float64 `xml:"-"`
		OriginalQuantityOrdered   float64 `xml:"-"`
		OriginalQuantityRest      float64 `xml:"-"`
		OriginalCost              float64 `xml:"-"`
		OriginalInPrice           float64 `xml:"-"`
	}{
		Alias:                     Alias(i),
		OriginalPrice:             i.Price,
		OriginalQuantity:          i.Quantity,
		OriginalVatRate:           i.VatRate,
		OriginalQuantityDelivered: i.QuantityDelivered,
		OriginalQuantityOrdered:   i.QuantityOrdered,
		OriginalQuantityRest:      i.QuantityRest,
		OriginalCost:              i.Cost,
		OriginalInPrice:           i.InPrice,
	}

	if i.Price != 0 {
		formattedPrice := fmt.Sprintf("%.2f", i.Price)
		modified.Price = &formattedPrice
	}
	if i.Quantity != 0 {
		formattedQuantity := fmt.Sprintf("%.2f", i.Quantity)
		modified.Quantity = &formattedQuantity
	}
	if i.VatRate != 0 {
		formattedVatRate := fmt.Sprintf("%.2f", i.VatRate)
		modified.VatRate = &formattedVatRate
	}
	if i.QuantityDelivered != 0 {
		formattedQuantityDelivered := fmt.Sprintf("%.2f", i.QuantityDelivered)
		modified.QuantityDelivered = &formattedQuantityDelivered
	}
	if i.QuantityOrdered != 0 {
		formattedQuantityOrdered := fmt.Sprintf("%.2f", i.QuantityOrdered)
		modified.QuantityOrdered = &formattedQuantityOrdered
	}
	if i.QuantityRest != 0 {
		formattedQuantityRest := fmt.Sprintf("%.2f", i.QuantityRest)
		modified.QuantityRest = &formattedQuantityRest
	}
	if i.Cost != 0 {
		formattedCost := fmt.Sprintf("%.2f", i.Cost)
		modified.Cost = &formattedCost
	}
	if i.InPrice != 0 {
		formattedInPrice := fmt.Sprintf("%.2f", i.InPrice)
		modified.InPrice = &formattedInPrice
	}

	return e.EncodeElement(modified, start)
}
