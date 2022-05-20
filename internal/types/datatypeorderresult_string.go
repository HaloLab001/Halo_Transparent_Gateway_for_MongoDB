// Code generated by "stringer -linecomment -type dataTypeOrderResult"; DO NOT EDIT.

package types

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[nullDataType-1]
	_ = x[nanDataType-2]
	_ = x[numbersDataType-3]
	_ = x[stringDataType-4]
	_ = x[objectDataType-5]
	_ = x[arrayDataType-6]
	_ = x[binDataType-7]
	_ = x[objectIdDataType-8]
	_ = x[booleanDataType-9]
	_ = x[dateDataType-10]
	_ = x[timestampDataType-11]
	_ = x[regexDataType-12]
}

const _dataTypeOrderResult_name = "nullDataTypenanDataTypenumbersDataTypestringDataTypeobjectDataTypearrayDataTypebinDataTypeobjectIdDataTypebooleanDataTypedateDataTypetimestampDataTyperegexDataType"

var _dataTypeOrderResult_index = [...]uint8{0, 12, 23, 38, 52, 66, 79, 90, 106, 121, 133, 150, 163}

func (i dataTypeOrderResult) String() string {
	i -= 1
	if i >= dataTypeOrderResult(len(_dataTypeOrderResult_index)-1) {
		return "dataTypeOrderResult(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _dataTypeOrderResult_name[_dataTypeOrderResult_index[i]:_dataTypeOrderResult_index[i+1]]
}
