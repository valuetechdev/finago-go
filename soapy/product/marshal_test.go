package product

import (
	"encoding/xml"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshaledProduct(t *testing.T) {
	assert := assert.New(t)
	row := Product{
		Price: float64(1122062),
	}

	x, err := xml.Marshal(row)
	assert.NoError(err)
	assert.Equal("<Product><Price>1122062.00</Price></Product>", string(x))

	row = Product{
		Price: float64(1337.37),
	}

	x, err = xml.Marshal(row)
	assert.NoError(err)
	assert.Equal("<Product><Price>1337.37</Price></Product>", string(x))
}
