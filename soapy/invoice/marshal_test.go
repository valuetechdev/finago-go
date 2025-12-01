package invoice

import (
	"encoding/xml"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshaledInvoiceRow(t *testing.T) {
	assert := assert.New(t)
	row := InvoiceRow{
		Price: float64(1122062),
	}

	x, err := xml.Marshal(row)
	assert.NoError(err)
	assert.Equal("<InvoiceRow><Price>1122062.00</Price></InvoiceRow>", string(x))

	row = InvoiceRow{
		Price: float64(1337.37),
	}

	x, err = xml.Marshal(row)
	assert.NoError(err)
	assert.Equal("<InvoiceRow><Price>1337.37</Price></InvoiceRow>", string(x))
}
