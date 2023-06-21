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
	_ = x[ErrUnauthorized-13]
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
	_ = x[ErrInvalidIndexSpecificationOption-197]
	_ = x[ErrInvalidPipelineOperator-168]
	_ = x[ErrNotImplemented-238]
	_ = x[ErrDuplicateKeyInsert-11000]
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
	_ = x[ErrUnsetPathCollision-31249]
	_ = x[ErrUnsetPathOverwrite-31250]
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

const _ErrorCode_name = "UnsetInternalErrorBadValueFailedToParseUnauthorizedTypeMismatchAuthenticationFailedIllegalOperationNamespaceNotFoundIndexNotFoundPathNotViableConflictingUpdateOperatorsCursorNotFoundNamespaceExistsDollarPrefixedFieldNameInvalidIDEmptyFieldNameCommandNotFoundImmutableFieldCannotCreateIndexInvalidOptionsInvalidNamespaceIndexOptionsConflictIndexKeySpecsConflictOperationFailedDocumentValidationFailureInvalidPipelineOperatorInvalidIndexSpecificationOptionNotImplementedLocation11000Location15947Location15948Location15955Location15958Location15959Location15969Location15973Location15974Location15975Location15976Location15981Location15998Location16020Location16410Location16872Location17276Location28667Location28724Location28812Location28818Location31002Location31119Location31120Location31249Location31250Location31253Location31254Location31324Location31325Location31394Location31395Location40156Location40157Location40158Location40160Location40234Location40237Location40238Location40272Location40323Location40352Location40353Location40414Location40415Location50840Location51024Location51075Location51091Location51108Location51246Location51247Location51270Location51272Location4822819Location5107200Location5107201Location5447000"

var _ErrorCode_map = map[ErrorCode]string{
	0:       _ErrorCode_name[0:5],
	1:       _ErrorCode_name[5:18],
	2:       _ErrorCode_name[18:26],
	9:       _ErrorCode_name[26:39],
	13:      _ErrorCode_name[39:51],
	14:      _ErrorCode_name[51:63],
	18:      _ErrorCode_name[63:83],
	20:      _ErrorCode_name[83:99],
	26:      _ErrorCode_name[99:116],
	27:      _ErrorCode_name[116:129],
	28:      _ErrorCode_name[129:142],
	40:      _ErrorCode_name[142:168],
	43:      _ErrorCode_name[168:182],
	48:      _ErrorCode_name[182:197],
	52:      _ErrorCode_name[197:220],
	53:      _ErrorCode_name[220:229],
	56:      _ErrorCode_name[229:243],
	59:      _ErrorCode_name[243:258],
	66:      _ErrorCode_name[258:272],
	67:      _ErrorCode_name[272:289],
	72:      _ErrorCode_name[289:303],
	73:      _ErrorCode_name[303:319],
	85:      _ErrorCode_name[319:339],
	86:      _ErrorCode_name[339:360],
	96:      _ErrorCode_name[360:375],
	121:     _ErrorCode_name[375:400],
	168:     _ErrorCode_name[400:423],
	197:     _ErrorCode_name[423:454],
	238:     _ErrorCode_name[454:468],
	11000:   _ErrorCode_name[468:481],
	15947:   _ErrorCode_name[481:494],
	15948:   _ErrorCode_name[494:507],
	15955:   _ErrorCode_name[507:520],
	15958:   _ErrorCode_name[520:533],
	15959:   _ErrorCode_name[533:546],
	15969:   _ErrorCode_name[546:559],
	15973:   _ErrorCode_name[559:572],
	15974:   _ErrorCode_name[572:585],
	15975:   _ErrorCode_name[585:598],
	15976:   _ErrorCode_name[598:611],
	15981:   _ErrorCode_name[611:624],
	15998:   _ErrorCode_name[624:637],
	16020:   _ErrorCode_name[637:650],
	16410:   _ErrorCode_name[650:663],
	16872:   _ErrorCode_name[663:676],
	17276:   _ErrorCode_name[676:689],
	28667:   _ErrorCode_name[689:702],
	28724:   _ErrorCode_name[702:715],
	28812:   _ErrorCode_name[715:728],
	28818:   _ErrorCode_name[728:741],
	31002:   _ErrorCode_name[741:754],
	31119:   _ErrorCode_name[754:767],
	31120:   _ErrorCode_name[767:780],
	31249:   _ErrorCode_name[780:793],
	31250:   _ErrorCode_name[793:806],
	31253:   _ErrorCode_name[806:819],
	31254:   _ErrorCode_name[819:832],
	31324:   _ErrorCode_name[832:845],
	31325:   _ErrorCode_name[845:858],
	31394:   _ErrorCode_name[858:871],
	31395:   _ErrorCode_name[871:884],
	40156:   _ErrorCode_name[884:897],
	40157:   _ErrorCode_name[897:910],
	40158:   _ErrorCode_name[910:923],
	40160:   _ErrorCode_name[923:936],
	40234:   _ErrorCode_name[936:949],
	40237:   _ErrorCode_name[949:962],
	40238:   _ErrorCode_name[962:975],
	40272:   _ErrorCode_name[975:988],
	40323:   _ErrorCode_name[988:1001],
	40352:   _ErrorCode_name[1001:1014],
	40353:   _ErrorCode_name[1014:1027],
	40414:   _ErrorCode_name[1027:1040],
	40415:   _ErrorCode_name[1040:1053],
	50840:   _ErrorCode_name[1053:1066],
	51024:   _ErrorCode_name[1066:1079],
	51075:   _ErrorCode_name[1079:1092],
	51091:   _ErrorCode_name[1092:1105],
	51108:   _ErrorCode_name[1105:1118],
	51246:   _ErrorCode_name[1118:1131],
	51247:   _ErrorCode_name[1131:1144],
	51270:   _ErrorCode_name[1144:1157],
	51272:   _ErrorCode_name[1157:1170],
	4822819: _ErrorCode_name[1170:1185],
	5107200: _ErrorCode_name[1185:1200],
	5107201: _ErrorCode_name[1200:1215],
	5447000: _ErrorCode_name[1215:1230],
}

func (i ErrorCode) String() string {
	if str, ok := _ErrorCode_map[i]; ok {
		return str
	}
	return "ErrorCode(" + strconv.FormatInt(int64(i), 10) + ")"
}
