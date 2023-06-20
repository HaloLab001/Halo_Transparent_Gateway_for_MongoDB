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
	_ = x[ErrAuthenticationFailed-18]
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
	_ = x[ErrInvalidPipelineOperator-168]
	_ = x[ErrNotImplemented-238]
	_ = x[ErrDuplicateKey-11000]
	_ = x[ErrSetBadExpression-40272]
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
	_ = x[ErrOperatorWrongLenOfArgs-16020]
	_ = x[ErrFieldPathInvalidName-16410]
	_ = x[ErrGroupInvalidFieldPath-16872]
	_ = x[ErrGroupUndefinedVariable-17276]
	_ = x[ErrInvalidArg-28667]
	_ = x[ErrSliceFirstArg-28724]
	_ = x[ErrStageUnsetNoPath-31119]
	_ = x[ErrStageUnsetArrElementInvalidType-31120]
	_ = x[ErrStageUnsetInvalidType-31002]
	_ = x[ErrStageUnwindNoPath-28812]
	_ = x[ErrStageUnwindNoPrefix-28818]
	_ = x[ErrProjectionInEx-31253]
	_ = x[ErrProjectionExIn-31254]
	_ = x[ErrAggregatePositionalProject-31324]
	_ = x[ErrAggregateInvalidExpression-31325]
	_ = x[ErrWrongPositionalOperatorLocation-31394]
	_ = x[ErrExclusionPositionalProjection-31395]
	_ = x[ErrStageCountNonString-40156]
	_ = x[ErrStageCountNonEmptyString-40157]
	_ = x[ErrStageCountBadPrefix-40158]
	_ = x[ErrStageCountBadValue-40160]
	_ = x[ErrAddFieldsExpressionWrongAmountOfArgs-40181]
	_ = x[ErrStageGroupUnaryOperator-40237]
	_ = x[ErrStageGroupMultipleAccumulator-40238]
	_ = x[ErrStageGroupInvalidAccumulator-40234]
	_ = x[ErrStageInvalid-40323]
	_ = x[ErrEmptyFieldPath-40352]
	_ = x[ErrInvalidFieldPath-40353]
	_ = x[ErrMissingField-40414]
	_ = x[ErrFailedToParseInput-40415]
	_ = x[ErrCollStatsIsNotFirstStage-40415]
	_ = x[ErrFreeMonitoringDisabled-50840]
	_ = x[ErrValueNegative-51024]
	_ = x[ErrRegexOptions-51075]
	_ = x[ErrRegexMissingParen-51091]
	_ = x[ErrBadRegexOption-51108]
	_ = x[ErrBadPositionalProjection-51246]
	_ = x[ErrElementMismatchPositionalProjection-51247]
	_ = x[ErrEmptySubProject-51270]
	_ = x[ErrEmptyProject-51272]
	_ = x[ErrDuplicateField-4822819]
	_ = x[ErrStageSkipBadValue-5107200]
	_ = x[ErrStageLimitInvalidArg-5107201]
	_ = x[ErrStageCollStatsInvalidArg-5447000]
}

const _ErrorCode_name = "UnsetInternalErrorBadValueFailedToParseTypeMismatchAuthenticationFailedIllegalOperationNamespaceNotFoundIndexNotFoundPathNotViableConflictingUpdateOperatorsCursorNotFoundNamespaceExistsDollarPrefixedFieldNameInvalidIDEmptyFieldNameCommandNotFoundImmutableFieldCannotCreateIndexInvalidOptionsInvalidNamespaceIndexOptionsConflictIndexKeySpecsConflictOperationFailedDocumentValidationFailureInvalidPipelineOperatorNotImplementedLocation11000Location15947Location15948Location15955Location15958Location15959Location15969Location15973Location15974Location15975Location15976Location15981Location15998Location16020Location16410Location16872Location17276Location28667Location28724Location28812Location28818Location31002Location31119Location31120Location31253Location31254Location31324Location31325Location31394Location31395Location40156Location40157Location40158Location40160Location40181Location40234Location40237Location40238Location40272Location40323Location40352Location40353Location40414Location40415Location50840Location51024Location51075Location51091Location51108Location51246Location51247Location51270Location51272Location4822819Location5107200Location5107201Location5447000"

var _ErrorCode_map = map[ErrorCode]string{
	0:       _ErrorCode_name[0:5],
	1:       _ErrorCode_name[5:18],
	2:       _ErrorCode_name[18:26],
	9:       _ErrorCode_name[26:39],
	14:      _ErrorCode_name[39:51],
	18:      _ErrorCode_name[51:71],
	20:      _ErrorCode_name[71:87],
	26:      _ErrorCode_name[87:104],
	27:      _ErrorCode_name[104:117],
	28:      _ErrorCode_name[117:130],
	40:      _ErrorCode_name[130:156],
	43:      _ErrorCode_name[156:170],
	48:      _ErrorCode_name[170:185],
	52:      _ErrorCode_name[185:208],
	53:      _ErrorCode_name[208:217],
	56:      _ErrorCode_name[217:231],
	59:      _ErrorCode_name[231:246],
	66:      _ErrorCode_name[246:260],
	67:      _ErrorCode_name[260:277],
	72:      _ErrorCode_name[277:291],
	73:      _ErrorCode_name[291:307],
	85:      _ErrorCode_name[307:327],
	86:      _ErrorCode_name[327:348],
	96:      _ErrorCode_name[348:363],
	121:     _ErrorCode_name[363:388],
	168:     _ErrorCode_name[388:411],
	238:     _ErrorCode_name[411:425],
	11000:   _ErrorCode_name[425:438],
	15947:   _ErrorCode_name[438:451],
	15948:   _ErrorCode_name[451:464],
	15955:   _ErrorCode_name[464:477],
	15958:   _ErrorCode_name[477:490],
	15959:   _ErrorCode_name[490:503],
	15969:   _ErrorCode_name[503:516],
	15973:   _ErrorCode_name[516:529],
	15974:   _ErrorCode_name[529:542],
	15975:   _ErrorCode_name[542:555],
	15976:   _ErrorCode_name[555:568],
	15981:   _ErrorCode_name[568:581],
	15998:   _ErrorCode_name[581:594],
	16020:   _ErrorCode_name[594:607],
	16410:   _ErrorCode_name[607:620],
	16872:   _ErrorCode_name[620:633],
	17276:   _ErrorCode_name[633:646],
	28667:   _ErrorCode_name[646:659],
	28724:   _ErrorCode_name[659:672],
	28812:   _ErrorCode_name[672:685],
	28818:   _ErrorCode_name[685:698],
	31002:   _ErrorCode_name[698:711],
	31119:   _ErrorCode_name[711:724],
	31120:   _ErrorCode_name[724:737],
	31253:   _ErrorCode_name[737:750],
	31254:   _ErrorCode_name[750:763],
	31324:   _ErrorCode_name[763:776],
	31325:   _ErrorCode_name[776:789],
	31394:   _ErrorCode_name[789:802],
	31395:   _ErrorCode_name[802:815],
	40156:   _ErrorCode_name[815:828],
	40157:   _ErrorCode_name[828:841],
	40158:   _ErrorCode_name[841:854],
	40160:   _ErrorCode_name[854:867],
	40181:   _ErrorCode_name[867:880],
	40234:   _ErrorCode_name[880:893],
	40237:   _ErrorCode_name[893:906],
	40238:   _ErrorCode_name[906:919],
	40272:   _ErrorCode_name[919:932],
	40323:   _ErrorCode_name[932:945],
	40352:   _ErrorCode_name[945:958],
	40353:   _ErrorCode_name[958:971],
	40414:   _ErrorCode_name[971:984],
	40415:   _ErrorCode_name[984:997],
	50840:   _ErrorCode_name[997:1010],
	51024:   _ErrorCode_name[1010:1023],
	51075:   _ErrorCode_name[1023:1036],
	51091:   _ErrorCode_name[1036:1049],
	51108:   _ErrorCode_name[1049:1062],
	51246:   _ErrorCode_name[1062:1075],
	51247:   _ErrorCode_name[1075:1088],
	51270:   _ErrorCode_name[1088:1101],
	51272:   _ErrorCode_name[1101:1114],
	4822819: _ErrorCode_name[1114:1129],
	5107200: _ErrorCode_name[1129:1144],
	5107201: _ErrorCode_name[1144:1159],
	5447000: _ErrorCode_name[1159:1174],
}

func (i ErrorCode) String() string {
	if str, ok := _ErrorCode_map[i]; ok {
		return str
	}
	return "ErrorCode(" + strconv.FormatInt(int64(i), 10) + ")"
}
