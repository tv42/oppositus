// Code generated by "stringer -type=Channel"; DO NOT EDIT.

package channels

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[stable-1]
	_ = x[beta-2]
	_ = x[alpha-3]
}

const _Channel_name = "stablebetaalpha"

var _Channel_index = [...]uint8{0, 6, 10, 15}

func (i Channel) String() string {
	i -= 1
	if i < 0 || i >= Channel(len(_Channel_index)-1) {
		return "Channel(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _Channel_name[_Channel_index[i]:_Channel_index[i+1]]
}
