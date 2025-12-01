package account

import (
	"encoding/xml"
	"fmt"
)

// NOTE: Custom marshaling for removing scientific notation from float64 Price
func (i Entry) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	type Alias Entry

	modified := struct {
		Alias
		Amount *string `xml:"Amount,omitempty"` // Override the original field
		// Fields that we're hiding from the original struct
		OriginalAmount float64 `xml:"-"`
	}{
		Alias:          Alias(i),
		OriginalAmount: i.Amount, // Hide the original Price
	}

	// Only include Price if it's not zero
	if i.Amount != 0 {
		formattedPrice := fmt.Sprintf("%.4f", i.Amount)
		modified.Amount = &formattedPrice
	}

	return e.EncodeElement(modified, start)
}
