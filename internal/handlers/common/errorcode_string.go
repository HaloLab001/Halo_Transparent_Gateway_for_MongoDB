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
	_ = x[ErrRegexOptions-51075]
	_ = x[ErrRegexMissingParen-51091]
	_ = x[ErrBadRegexOption-51108]
}

const _ErrorCode_name = "UnsetInternalErrorBadValueFailedToParseTypeMismatchNamespaceNotFoundUnsuitableValueTypeConflictingUpdateOperatorsNamespaceExistsInvalidIDEmptyNameCommandNotFoundInvalidNamespaceOperationFailedDocumentValidationFailureNotImplementedMechanismUnavailableLocation15974Location15975Location28667Location28724Location31253Location31254Location40352Location40414Location40415Location50840Location51075Location51091Location51108"

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
	53:    _ErrorCode_name[128:137],
	56:    _ErrorCode_name[137:146],
	59:    _ErrorCode_name[146:161],
	73:    _ErrorCode_name[161:177],
	96:    _ErrorCode_name[177:192],
	121:   _ErrorCode_name[192:217],
	238:   _ErrorCode_name[217:231],
	334:   _ErrorCode_name[231:251],
	15974: _ErrorCode_name[251:264],
	15975: _ErrorCode_name[264:277],
	28667: _ErrorCode_name[277:290],
	28724: _ErrorCode_name[290:303],
	31253: _ErrorCode_name[303:316],
	31254: _ErrorCode_name[316:329],
	40352: _ErrorCode_name[329:342],
	40414: _ErrorCode_name[342:355],
	40415: _ErrorCode_name[355:368],
	50840: _ErrorCode_name[368:381],
	51075: _ErrorCode_name[381:394],
	51091: _ErrorCode_name[394:407],
	51108: _ErrorCode_name[407:420],
}

func (i ErrorCode) String() string {
	if str, ok := _ErrorCode_map[i]; ok {
		return str
	}
	return "ErrorCode(" + strconv.FormatInt(int64(i), 10) + ")"
}
