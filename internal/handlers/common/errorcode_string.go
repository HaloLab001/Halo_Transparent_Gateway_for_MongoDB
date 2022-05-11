// Code generated by "stringer -linecomment -type ErrorCode"; DO NOT EDIT.

package common

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[errInternalError-1]
	_ = x[ErrBadValue-2]
	_ = x[ErrFailedToParse-9]
	_ = x[ErrTypeMismatch-14]
	_ = x[ErrNamespaceNotFound-26]
	_ = x[ErrNamespaceExists-48]
	_ = x[ErrCommandNotFound-59]
	_ = x[ErrNotImplemented-238]
	_ = x[ErrSortBadValue-15974]
	_ = x[ErrInvalidArg-28667]
	_ = x[ErrSliceFirstArg-28724]
	_ = x[ErrProjectionInEx-31253]
	_ = x[ErrProjectionExIn-31254]
	_ = x[ErrRegexOptions-51075]
	_ = x[ErrRegexMissingParen-51091]
}

const _ErrorCode_name = "InternalErrorBadValueFailedToParseTypeMismatchNamespaceNotFoundNamespaceExistsCommandNotFoundNotImplementedLocation15974Location28667Location28724Location31253Location31254Location51075Location51091"

var _ErrorCode_map = map[ErrorCode]string{
	1:     _ErrorCode_name[0:13],
	2:     _ErrorCode_name[13:21],
	9:     _ErrorCode_name[21:34],
	14:    _ErrorCode_name[34:46],
	26:    _ErrorCode_name[46:63],
	48:    _ErrorCode_name[63:78],
	59:    _ErrorCode_name[78:93],
	238:   _ErrorCode_name[93:107],
	15974: _ErrorCode_name[107:120],
	28667: _ErrorCode_name[120:133],
	28724: _ErrorCode_name[133:146],
	31253: _ErrorCode_name[146:159],
	31254: _ErrorCode_name[159:172],
	51075: _ErrorCode_name[172:185],
	51091: _ErrorCode_name[185:198],
}

func (i ErrorCode) String() string {
	if str, ok := _ErrorCode_map[i]; ok {
		return str
	}
	return "ErrorCode(" + strconv.FormatInt(int64(i), 10) + ")"
}
