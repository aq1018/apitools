package rangeheader

type rangeDescriptor struct {
	field          string
	startExclusive bool
	start          string
	end            string
	endExclusive   bool
}

func parseRangeDescriptor(rangeStr string) (*rangeDescriptor, error) {
	matches := rangeRegex.FindStringSubmatch(rangeStr)
	if matches == nil {
		// no match, it is not a valid range string
		return nil, ErrInvalidRange
	}

	var v rangeDescriptor

	// give the matched strings a name
	v.field = matches[1]
	v.startExclusive = matches[2] != ""
	v.start = matches[3]
	v.end = matches[4]
	v.endExclusive = matches[5] != ""

	return &v, nil
}

func (rt *rangeDescriptor) validate() error {
	// ensure the fields are valid
	if rt.start != "" && rt.end != "" {
		// name 1..100 is not allowed
		return ErrInvalidRange
	}
	if rt.startExclusive && rt.start == "" {
		// name ]..100 is not allowed
		return ErrInvalidRange
	}
	if rt.endExclusive && rt.end == "" {
		// name 1..[ is not allowed
		return ErrInvalidRange
	}

	return nil
}

func (rt *rangeDescriptor) value() (value string) {
	if rt.start != "" {
		return rt.start
	} else {
		return rt.end
	}
}

func (rt *rangeDescriptor) operator() RangeOperator {
	if rt.start != "" {
		if rt.startExclusive {
			return RangeOperatorGreaterThan
		} else {
			return RangeOperatorGreaterOrEqualTo
		}
	} else {
		if rt.endExclusive {
			return RangeOperatorLessThan
		} else {
			return RangeOperatorLessOrEqualTo
		}
	}
}
