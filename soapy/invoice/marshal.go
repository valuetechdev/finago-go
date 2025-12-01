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
		Price *string `xml:"Price,omitempty"` // Override the original field
		// Fields that we're hiding from the original struct
		OriginalPrice float64 `xml:"-"`
	}{
		Alias:         Alias(i),
		OriginalPrice: i.Price, // Hide the original Price
	}

	// Only include Price if it's not zero
	if i.Price != 0 {
		formattedPrice := fmt.Sprintf("%.2f", i.Price)
		modified.Price = &formattedPrice
	}

	return e.EncodeElement(modified, start)
}
