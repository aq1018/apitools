package rangeheader

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type StringValue string

func (sv *StringValue) Set(v string) error {
	*sv = StringValue(v)
	return nil
}

func TestParse(t *testing.T) {
	r := Range{
		Max:   1000,
		Sort:  RangeSortAscending,
		Value: new(StringValue),
	}

	err := Parse("id ]foo..;sort=desc;max=100", &r)
	assert.NoError(t, err)
	assert.Equal(t, "id", r.Field)
	assert.Equal(t, uint(100), r.Max)
	assert.Equal(t, RangeOperatorGreaterThan, r.Operator)
	assert.Equal(t, RangeSortDescending, r.Sort)
	v, _ := (r.Value).(*StringValue)
	assert.Equal(t, "foo", string(*v))
}
