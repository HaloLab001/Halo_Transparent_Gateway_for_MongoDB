// Code generated by "stringer -linecomment -type ErrorCode"; DO NOT EDIT.

package commonerrors

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
	_ = x[ErrIllegalOperation-20]
	_ = x[ErrNamespaceNotFound-26]
	_ = x[ErrIndexNotFound-27]
	_ = x[ErrUnsuitableValueType-28]
	_ = x[ErrConflictingUpdateOperators-40]
	_ = x[ErrCursorNotFound-43]
	_ = x[ErrNamespaceExists-48]
	_ = x[ErrDollarPrefixedFieldName-52]
	_ = x[ErrInvalidID-53]
	_ = x[ErrEmptyName-56]
	_ = x[ErrCommandNotFound-59]
	_ = x[ErrImmutableField-66]
	_ = x[ErrCannotCreateIndex-67]
	_ = x[ErrInvalidOptions-72]
	_ = x[ErrInvalidNamespace-73]
	_ = x[ErrIndexOptionsConflict-85]
	_ = x[ErrIndexKeySpecsConflict-86]
	_ = x[ErrOperationFailed-96]
	_ = x[ErrDocumentValidationFailure-121]
	_ = x[ErrNotImplemented-238]
	_ = x[ErrMechanismUnavailable-334]
	_ = x[ErrDuplicateKey-11000]
	_ = x[ErrStageGroupInvalidFields-15947]
	_ = x[ErrStageGroupID-15948]
	_ = x[ErrStageGroupMissingID-15955]
	_ = x[ErrStageLimitZero-15958]
	_ = x[ErrMatchBadExpression-15959]
	_ = x[ErrProjectBadExpression-15969]
	_ = x[ErrSortBadExpression-15973]
	_ = x[ErrSortBadValue-15974]
	_ = x[ErrSortBadOrder-15975]
	_ = x[ErrSortMissingKey-15976]
	_ = x[ErrStageUnwindWrongType-15981]
	_ = x[ErrPathContainsEmptyElement-15998]
	_ = x[ErrFieldPathInvalidName-16410]
	_ = x[ErrGroupInvalidFieldPath-16872]
	_ = x[ErrGroupUndefinedVariable-17276]
	_ = x[ErrInvalidArg-28667]
	_ = x[ErrSliceFirstArg-28724]
	_ = x[ErrStageUnwindNoPath-28812]
	_ = x[ErrStageUnwindNoPrefix-28818]
	_ = x[ErrProjectionInEx-31253]
	_ = x[ErrProjectionExIn-31254]
	_ = x[ErrStageCountNonString-40156]
	_ = x[ErrStageCountNonEmptyString-40157]
	_ = x[ErrStageCountBadPrefix-40158]
	_ = x[ErrStageCountBadValue-40160]
	_ = x[ErrStageGroupUnaryOperator-40237]
	_ = x[ErrStageGroupMultipleAccumulator-40238]
	_ = x[ErrStageGroupInvalidAccumulator-40234]
	_ = x[ErrStageInvalid-40323]
	_ = x[ErrEmptyFieldPath-40352]
	_ = x[ErrMissingField-40414]
	_ = x[ErrFailedToParseInput-40415]
	_ = x[ErrCollStatsIsNotFirstStage-40415]
	_ = x[ErrFreeMonitoringDisabled-50840]
	_ = x[ErrValueNegative-51024]
	_ = x[ErrRegexOptions-51075]
	_ = x[ErrRegexMissingParen-51091]
	_ = x[ErrBadRegexOption-51108]
	_ = x[ErrDuplicateField-4822819]
	_ = x[ErrStageSkipBadValue-5107200]
	_ = x[ErrStageLimitInvalidArg-5107201]
	_ = x[ErrStageCollStatsInvalidArg-5447000]
}

const _ErrorCode_name = "UnsetInternalErrorBadValueFailedToParseTypeMismatchIllegalOperationNamespaceNotFoundIndexNotFoundPathNotViableConflictingUpdateOperatorsCursorNotFoundNamespaceExistsDollarPrefixedFieldNameInvalidIDEmptyNameCommandNotFoundImmutableFieldCannotCreateIndexInvalidOptionsInvalidNamespaceIndexOptionsConflictIndexKeySpecsConflictOperationFailedDocumentValidationFailureNotImplementedMechanismUnavailableLocation11000Location15947Location15948Location15955Location15958Location15959Location15969Location15973Location15974Location15975Location15976Location15981Location15998Location16410Location16872Location17276Location28667Location28724Location28812Location28818Location31253Location31254Location40156Location40157Location40158Location40160Location40234Location40237Location40238Location40323Location40352Location40414Location40415Location50840Location51024Location51075Location51091Location51108Location4822819Location5107200Location5107201Location5447000"

var _ErrorCode_map = map[ErrorCode]string{
	0:       _ErrorCode_name[0:5],
	1:       _ErrorCode_name[5:18],
	2:       _ErrorCode_name[18:26],
	9:       _ErrorCode_name[26:39],
	14:      _ErrorCode_name[39:51],
	20:      _ErrorCode_name[51:67],
	26:      _ErrorCode_name[67:84],
	27:      _ErrorCode_name[84:97],
	28:      _ErrorCode_name[97:110],
	40:      _ErrorCode_name[110:136],
	43:      _ErrorCode_name[136:150],
	48:      _ErrorCode_name[150:165],
	52:      _ErrorCode_name[165:188],
	53:      _ErrorCode_name[188:197],
	56:      _ErrorCode_name[197:206],
	59:      _ErrorCode_name[206:221],
	66:      _ErrorCode_name[221:235],
	67:      _ErrorCode_name[235:252],
	72:      _ErrorCode_name[252:266],
	73:      _ErrorCode_name[266:282],
	85:      _ErrorCode_name[282:302],
	86:      _ErrorCode_name[302:323],
	96:      _ErrorCode_name[323:338],
	121:     _ErrorCode_name[338:363],
	238:     _ErrorCode_name[363:377],
	334:     _ErrorCode_name[377:397],
	11000:   _ErrorCode_name[397:410],
	15947:   _ErrorCode_name[410:423],
	15948:   _ErrorCode_name[423:436],
	15955:   _ErrorCode_name[436:449],
	15958:   _ErrorCode_name[449:462],
	15959:   _ErrorCode_name[462:475],
	15969:   _ErrorCode_name[475:488],
	15973:   _ErrorCode_name[488:501],
	15974:   _ErrorCode_name[501:514],
	15975:   _ErrorCode_name[514:527],
	15976:   _ErrorCode_name[527:540],
	15981:   _ErrorCode_name[540:553],
	15998:   _ErrorCode_name[553:566],
	16410:   _ErrorCode_name[566:579],
	16872:   _ErrorCode_name[579:592],
	17276:   _ErrorCode_name[592:605],
	28667:   _ErrorCode_name[605:618],
	28724:   _ErrorCode_name[618:631],
	28812:   _ErrorCode_name[631:644],
	28818:   _ErrorCode_name[644:657],
	31253:   _ErrorCode_name[657:670],
	31254:   _ErrorCode_name[670:683],
	40156:   _ErrorCode_name[683:696],
	40157:   _ErrorCode_name[696:709],
	40158:   _ErrorCode_name[709:722],
	40160:   _ErrorCode_name[722:735],
	40234:   _ErrorCode_name[735:748],
	40237:   _ErrorCode_name[748:761],
	40238:   _ErrorCode_name[761:774],
	40323:   _ErrorCode_name[774:787],
	40352:   _ErrorCode_name[787:800],
	40414:   _ErrorCode_name[800:813],
	40415:   _ErrorCode_name[813:826],
	50840:   _ErrorCode_name[826:839],
	51024:   _ErrorCode_name[839:852],
	51075:   _ErrorCode_name[852:865],
	51091:   _ErrorCode_name[865:878],
	51108:   _ErrorCode_name[878:891],
	4822819: _ErrorCode_name[891:906],
	5107200: _ErrorCode_name[906:921],
	5107201: _ErrorCode_name[921:936],
	5447000: _ErrorCode_name[936:951],
}

func (i ErrorCode) String() string {
	if str, ok := _ErrorCode_map[i]; ok {
		return str
	}
	return "ErrorCode(" + strconv.FormatInt(int64(i), 10) + ")"
}
