package tracker

import "github.com/morrocker/errors"

func getInt64(n interface{}) (int64, error) {
	op := "tracker.getInt64()"
	var out int64
	if x, ok := n.(int64); ok {
		out = x
		return out, nil
	}
	if x, ok := n.(int); ok {
		out = int64(x)
		return out, nil
	}
	if x, ok := n.(int8); ok {
		out = int64(x)
		return out, nil
	}
	if x, ok := n.(int16); ok {
		out = int64(x)
		return out, nil
	}
	if x, ok := n.(int32); ok {
		out = int64(x)
		return out, nil
	}
	if x, ok := n.(uint); ok {
		out = int64(x)
		return out, nil
	}
	if x, ok := n.(uint8); ok {
		out = int64(x)
		return out, nil
	}
	if x, ok := n.(uint16); ok {
		out = int64(x)
		return out, nil
	}
	if x, ok := n.(uint32); ok {
		out = int64(x)
		return out, nil
	}
	if x, ok := n.(uint64); ok {
		out = int64(x)
		return out, nil
	}
	return 0, errors.New(op, "Given is not a valid intX")
}
