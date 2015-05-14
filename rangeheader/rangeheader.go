package rangeheader

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func Parse(contentRange string, r *Range) error {
	tokens := tokenize(contentRange, ";")

	fmt.Println(tokens)

	desc, err := parseRangeDescriptor(tokens[0])
	fmt.Println(desc)
	if err != nil {
		return err
	}

	if err := desc.validate(); err != nil {
		return err
	}

	options, err := parseRangeOptions(tokens[1:])
	if err != nil {
		return err
	}

	max, err := options.max(r.Max)
	if err != nil {
		return err
	}

	sort, err := options.sort(r.Sort)
	if err != nil {
		return err
	}

	if err := r.Value.Set(desc.value()); err != nil {
		return err
	}

	r.Field = desc.field
	r.Operator = desc.operator()
	r.Max = max
	r.Sort = sort

	return nil
}

func tokenize(str, tokenizer string) []string {
	var result []string
	for _, token := range strings.Split(str, tokenizer) {
		v := strings.Trim(token, " \t")
		if v != "" {
			result = append(result, v)
		}
	}
	return result
}

func parseRangeOptions(tokens []string) (rangeOptions, error) {
	result := make(rangeOptions)

	for _, token := range tokens {
		tuple := tokenize(token, "=")
		if len(tuple) != 2 || tuple[0] == "" {
			return nil, ErrInvalidRange
		}
		result[tuple[0]] = tuple[1]
	}

	return result, nil
}

func (ro rangeOptions) max(defaultMax uint) (uint, error) {
	maxStr, ok := ro["max"]

	if !ok || maxStr == "" {
		return defaultMax, nil
	}

	v, err := strconv.ParseUint(maxStr, 10, 32)
	if err != nil {
		// max is not integer
		return 0, ErrInvalidRange
	}
	return uint(v), nil
}

func (ro rangeOptions) sort(defaultSort RangeSort) (RangeSort, error) {
	orderStr, ok := ro["sort"]

	if !ok || orderStr == "" {
		return defaultSort, nil
	}

	switch orderStr {
	case "asc":
		return RangeSortAscending, nil
	case "desc":
		return RangeSortDescending, nil
	default:
		return 0, ErrInvalidRange
	}
}

type (
	RangeOperator int

	RangeSort int

	RangeValue interface {
		Set(string) error
	}

	Range struct {
		Field    string
		Value    RangeValue
		Max      uint
		Operator RangeOperator
		Sort     RangeSort
	}

	rangeOptions map[string]string
)

var (
	ErrInvalidRange = errors.New("invalid range")
	rangeRegex      = regexp.MustCompile(`\A(\w+)\s+(\]?)([\w()-]*)\.\.([\w()-]*)(\[?)\z`)
	sortRegex       = regexp.MustCompile(`\Asort=(asc|desc)\z`)
)

const (
	RangeSortAscending RangeSort = iota
	RangeSortDescending
)

const (
	RangeOperatorGreaterOrEqualTo RangeOperator = iota
	RangeOperatorGreaterThan
	RangeOperatorLessThan
	RangeOperatorLessOrEqualTo
)
