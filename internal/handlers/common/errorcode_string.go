// Code generated by "stringer -linecomment -type ErrorCode"; DO NOT EDIT.

package common

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[errUnset-0]
	_ = x[errInternalError-1]
	_ = x[ErrBadValue-2]
	_ = x[ErrFailedToParse-9]
	_ = x[ErrNamespaceNotFound-26]
	_ = x[ErrNamespaceExists-48]
	_ = x[ErrCommandNotFound-59]
	_ = x[ErrNotImplemented-238]
	_ = x[ErrSortBadValue-15974]
	_ = x[ErrProjectionInEx-31253]
	_ = x[ErrProjectionExIn-31254]
	_ = x[ErrRegexOptions-51075]
}

const (
	_ErrorCode_name_0 = "UnsetInternalErrorBadValue"
	_ErrorCode_name_1 = "FailedToParse"
	_ErrorCode_name_2 = "NamespaceNotFound"
	_ErrorCode_name_3 = "NamespaceExists"
	_ErrorCode_name_4 = "CommandNotFound"
	_ErrorCode_name_5 = "NotImplemented"
	_ErrorCode_name_6 = "Location15974"
	_ErrorCode_name_7 = "Location31253Location31254"
	_ErrorCode_name_8 = "Location51075"
)

var (
	_ErrorCode_index_0 = [...]uint8{0, 5, 18, 26}
	_ErrorCode_index_7 = [...]uint8{0, 13, 26}
)

func (i ErrorCode) String() string {
	switch {
	case 0 <= i && i <= 2:
		return _ErrorCode_name_0[_ErrorCode_index_0[i]:_ErrorCode_index_0[i+1]]
	case i == 9:
		return _ErrorCode_name_1
	case i == 26:
		return _ErrorCode_name_2
	case i == 48:
		return _ErrorCode_name_3
	case i == 59:
		return _ErrorCode_name_4
	case i == 238:
		return _ErrorCode_name_5
	case i == 15974:
		return _ErrorCode_name_6
	case 31253 <= i && i <= 31254:
		i -= 31253
		return _ErrorCode_name_7[_ErrorCode_index_7[i]:_ErrorCode_index_7[i+1]]
	case i == 51075:
		return _ErrorCode_name_8
	default:
		return "ErrorCode(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
