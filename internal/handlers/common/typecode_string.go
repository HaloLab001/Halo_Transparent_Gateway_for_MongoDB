// Code generated by "stringer -linecomment -type typeCode"; DO NOT EDIT.

package common

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[typeCodeUnknown-0]
	_ = x[typeCodeDouble-1]
	_ = x[typeCodeString-2]
	_ = x[typeCodeObject-3]
	_ = x[typeCodeArray-4]
	_ = x[typeCodeBinData-5]
	_ = x[typeCodeObjectID-7]
	_ = x[typeCodeBool-8]
	_ = x[typeCodeDate-9]
	_ = x[typeCodeNull-10]
	_ = x[typeCodeRegex-11]
	_ = x[typeCodeInt-16]
	_ = x[typeCodeTimestamp-17]
	_ = x[typeCodeLong-18]
}

const (
	_typeCode_name_0 = "typeCodeUnknowndoublestringobjectarraybinData"
	_typeCode_name_1 = "objectIdbooldatenullregex"
	_typeCode_name_2 = "inttimestamplong"
)

var (
	_typeCode_index_0 = [...]uint8{0, 15, 21, 27, 33, 38, 45}
	_typeCode_index_1 = [...]uint8{0, 8, 12, 16, 20, 25}
	_typeCode_index_2 = [...]uint8{0, 3, 12, 16}
)

func (i typeCode) String() string {
	switch {
	case 0 <= i && i <= 5:
		return _typeCode_name_0[_typeCode_index_0[i]:_typeCode_index_0[i+1]]
	case 7 <= i && i <= 11:
		i -= 7
		return _typeCode_name_1[_typeCode_index_1[i]:_typeCode_index_1[i+1]]
	case 16 <= i && i <= 18:
		i -= 16
		return _typeCode_name_2[_typeCode_index_2[i]:_typeCode_index_2[i+1]]
	default:
		return "typeCode(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
