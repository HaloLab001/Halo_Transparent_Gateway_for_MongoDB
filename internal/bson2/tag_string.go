// Code generated by "stringer -linecomment -type tag"; DO NOT EDIT.

package bson2

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[tagFloat64-1]
	_ = x[tagString-2]
	_ = x[tagDocument-3]
	_ = x[tagArray-4]
	_ = x[tagBinary-5]
	_ = x[tagUndefined-6]
	_ = x[tagObjectID-7]
	_ = x[tagBool-8]
	_ = x[tagTime-9]
	_ = x[tagNull-10]
	_ = x[tagRegex-11]
	_ = x[tagDBPointer-12]
	_ = x[tagJavaScript-13]
	_ = x[tagSymbol-14]
	_ = x[tagJavaScriptScope-15]
	_ = x[tagInt32-16]
	_ = x[tagTimestamp-17]
	_ = x[tagInt64-18]
	_ = x[tagDecimal-19]
	_ = x[tagMinKey-255]
	_ = x[tagMaxKey-127]
}

const (
	_tag_name_0 = "Float64StringDocumentArrayBinaryUndefinedObjectIDBoolTimeNullRegexDBPointerJavaScriptSymbolJavaScriptScopeInt32TimestampInt64Decimal"
	_tag_name_1 = "MaxKey"
	_tag_name_2 = "MinKey"
)

var (
	_tag_index_0 = [...]uint8{0, 7, 13, 21, 26, 32, 41, 49, 53, 57, 61, 66, 75, 85, 91, 106, 111, 120, 125, 132}
)

func (i tag) String() string {
	switch {
	case 1 <= i && i <= 19:
		i -= 1
		return _tag_name_0[_tag_index_0[i]:_tag_index_0[i+1]]
	case i == 127:
		return _tag_name_1
	case i == 255:
		return _tag_name_2
	default:
		return "tag(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
