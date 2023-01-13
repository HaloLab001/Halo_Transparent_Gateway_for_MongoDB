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
	_ = x[ErrCursorNotFound-43]
	_ = x[ErrNamespaceExists-48]
	_ = x[ErrInvalidID-53]
	_ = x[ErrEmptyName-56]
	_ = x[ErrCommandNotFound-59]
	_ = x[ErrInvalidNamespace-73]
	_ = x[ErrOperationFailed-96]
	_ = x[ErrDocumentValidationFailure-121]
	_ = x[ErrNotImplemented-238]
	_ = x[ErrMechanismUnavailable-334]
	_ = x[ErrSortBadValue-15974]
	_ = x[ErrSortBadOrder-15975]
	_ = x[ErrInvalidArg-28667]
	_ = x[ErrSliceFirstArg-28724]
	_ = x[ErrProjectionInEx-31253]
	_ = x[ErrProjectionExIn-31254]
	_ = x[ErrEmptyFieldPath-40352]
	_ = x[ErrMissingField-40414]
	_ = x[ErrFailedToParseInput-40415]
	_ = x[ErrFreeMonitoringDisabled-50840]
	_ = x[ErrBatchSizeNegative-51024]
	_ = x[ErrRegexOptions-51075]
	_ = x[ErrRegexMissingParen-51091]
	_ = x[ErrBadRegexOption-51108]
}

const _ErrorCode_name = "UnsetInternalErrorBadValueFailedToParseTypeMismatchNamespaceNotFoundUnsuitableValueTypeConflictingUpdateOperatorsCursorNotFoundNamespaceExistsInvalidIDEmptyNameCommandNotFoundInvalidNamespaceOperationFailedDocumentValidationFailureNotImplementedMechanismUnavailableLocation15974Location15975Location28667Location28724Location31253Location31254Location40352Location40414Location40415Location50840Location51024Location51075Location51091Location51108"

var _ErrorCode_map = map[ErrorCode]string{
	0:     _ErrorCode_name[0:5],
	1:     _ErrorCode_name[5:18],
	2:     _ErrorCode_name[18:26],
	9:     _ErrorCode_name[26:39],
	14:    _ErrorCode_name[39:51],
	26:    _ErrorCode_name[51:68],
	28:    _ErrorCode_name[68:87],
	40:    _ErrorCode_name[87:113],
	43:    _ErrorCode_name[113:127],
	48:    _ErrorCode_name[127:142],
	53:    _ErrorCode_name[142:151],
	56:    _ErrorCode_name[151:160],
	59:    _ErrorCode_name[160:175],
	73:    _ErrorCode_name[175:191],
	96:    _ErrorCode_name[191:206],
	121:   _ErrorCode_name[206:231],
	238:   _ErrorCode_name[231:245],
	334:   _ErrorCode_name[245:265],
	15974: _ErrorCode_name[265:278],
	15975: _ErrorCode_name[278:291],
	28667: _ErrorCode_name[291:304],
	28724: _ErrorCode_name[304:317],
	31253: _ErrorCode_name[317:330],
	31254: _ErrorCode_name[330:343],
	40352: _ErrorCode_name[343:356],
	40414: _ErrorCode_name[356:369],
	40415: _ErrorCode_name[369:382],
	50840: _ErrorCode_name[382:395],
	51024: _ErrorCode_name[395:408],
	51075: _ErrorCode_name[408:421],
	51091: _ErrorCode_name[421:434],
	51108: _ErrorCode_name[434:447],
}

func (i ErrorCode) String() string {
	if str, ok := _ErrorCode_map[i]; ok {
		return str
	}
	return "ErrorCode(" + strconv.FormatInt(int64(i), 10) + ")"
}
