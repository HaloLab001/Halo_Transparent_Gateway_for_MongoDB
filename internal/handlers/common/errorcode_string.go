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
	_ = x[ErrUnsuitableValueType-28]
	_ = x[ErrConflictingUpdateOperators-40]
	_ = x[ErrNamespaceExists-48]
	_ = x[ErrCommandNotFound-59]
	_ = x[ErrInvalidNamespace-73]
	_ = x[ErrOperationFailed-96]
	_ = x[ErrDocumentValidationFailure-121]
	_ = x[ErrNotImplemented-238]
	_ = x[ErrFailedToParseInput-40415]
	_ = x[ErrSortBadValue-15974]
	_ = x[ErrSortBadOrder-15975]
	_ = x[ErrInvalidArg-28667]
	_ = x[ErrSliceFirstArg-28724]
	_ = x[ErrProjectionInEx-31253]
	_ = x[ErrProjectionExIn-31254]
	_ = x[ErrMissingField-40414]
	_ = x[ErrFreeMonitoringDisabled-50840]
	_ = x[ErrRegexOptions-51075]
	_ = x[ErrRegexMissingParen-51091]
}

const _ErrorCode_name = "UnsetInternalErrorBadValueFailedToParseTypeMismatchNamespaceNotFoundUnsuitableValueTypeConflictingUpdateOperatorsNamespaceExistsCommandNotFoundInvalidNamespaceOperationFailedDocumentValidationFailureNotImplementedLocation15974Location15975Location28667Location28724Location31253Location31254Location40414Location40415Location50840Location51075Location51091"

var _ErrorCode_map = map[ErrorCode]string{
	0:     _ErrorCode_name[0:5],
	1:     _ErrorCode_name[5:18],
	2:     _ErrorCode_name[18:26],
	9:     _ErrorCode_name[26:39],
	14:    _ErrorCode_name[39:51],
	26:    _ErrorCode_name[51:68],
	28:    _ErrorCode_name[68:87],
	40:    _ErrorCode_name[87:113],
	48:    _ErrorCode_name[113:128],
	59:    _ErrorCode_name[128:143],
	73:    _ErrorCode_name[143:159],
	96:    _ErrorCode_name[159:174],
	121:   _ErrorCode_name[174:199],
	238:   _ErrorCode_name[199:213],
	15974: _ErrorCode_name[213:226],
	15975: _ErrorCode_name[226:239],
	28667: _ErrorCode_name[239:252],
	28724: _ErrorCode_name[252:265],
	31253: _ErrorCode_name[265:278],
	31254: _ErrorCode_name[278:291],
	40414: _ErrorCode_name[291:304],
	40415: _ErrorCode_name[304:317],
	50840: _ErrorCode_name[317:330],
	51075: _ErrorCode_name[330:343],
	51091: _ErrorCode_name[343:356],
}

func (i ErrorCode) String() string {
	if str, ok := _ErrorCode_map[i]; ok {
		return str
	}
	return "ErrorCode(" + strconv.FormatInt(int64(i), 10) + ")"
}
