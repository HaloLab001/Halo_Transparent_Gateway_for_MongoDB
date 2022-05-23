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
	_ = x[ErrTypeMismatch-14]
	_ = x[ErrNamespaceNotFound-26]
	_ = x[ErrNamespaceExists-48]
	_ = x[ErrCommandNotFound-59]
	_ = x[ErrInvalidNamespace-73]
	_ = x[ErrNotImplemented-238]
	_ = x[ErrSortBadValue-15974]
	_ = x[ErrInvalidArg-28667]
	_ = x[ErrSliceFirstArg-28724]
	_ = x[ErrElemMatchInclusionInExclusion-31253]
	_ = x[ErrElemMatchExclusionInInclusion-31254]
	_ = x[ErrElemMatchPositional-31255]
	_ = x[ErrElemMatchObjectRequired-31274]
	_ = x[ErrElemMatchNestedField-31275]
	_ = x[ErrRegexOptions-51075]
	_ = x[ErrRegexMissingParen-51091]
}

const _ErrorCode_name = "UnsetInternalErrorBadValueFailedToParseTypeMismatchNamespaceNotFoundNamespaceExistsCommandNotFoundInvalidNamespaceNotImplementedLocation15974Location28667Location28724Location31253Location31254Location31255Location31274Location31275Location51075Location51091"

var _ErrorCode_map = map[ErrorCode]string{
	0:     _ErrorCode_name[0:5],
	1:     _ErrorCode_name[5:18],
	2:     _ErrorCode_name[18:26],
	9:     _ErrorCode_name[26:39],
	14:    _ErrorCode_name[39:51],
	26:    _ErrorCode_name[51:68],
	48:    _ErrorCode_name[68:83],
	59:    _ErrorCode_name[83:98],
	73:    _ErrorCode_name[98:114],
	238:   _ErrorCode_name[114:128],
	15974: _ErrorCode_name[128:141],
	28667: _ErrorCode_name[141:154],
	28724: _ErrorCode_name[154:167],
	31253: _ErrorCode_name[167:180],
	31254: _ErrorCode_name[180:193],
	31255: _ErrorCode_name[193:206],
	31274: _ErrorCode_name[206:219],
	31275: _ErrorCode_name[219:232],
	51075: _ErrorCode_name[232:245],
	51091: _ErrorCode_name[245:258],
}

func (i ErrorCode) String() string {
	if str, ok := _ErrorCode_map[i]; ok {
		return str
	}
	return "ErrorCode(" + strconv.FormatInt(int64(i), 10) + ")"
}
