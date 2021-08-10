package info

import "strings"

func DecodeCSV(msg string) LogInfo {
	s := strings.Split(msg, ",")
	return LogInfo{
		UnixDate:       s[0],
		UtcNanoseconds: 0,
		Index:          0,
		Kind:           s[1],
		Elements:       s[2:],
	}
}

type LogInfo struct {
	UnixDate       string
	UtcNanoseconds int64
	Index          uint32
	Kind           string
	Elements       []string
}

func (f *LogInfo) Valid() bool {
	if f.UnixDate == "" {
		return false
	}
	if f.Kind == "" {
		return false
	}
	if len(f.Elements) == 0 {
		return false
	}
	for _, element := range f.Elements {
		if element == "" {
			return false
		}
	}
	return true
}

func (format *LogInfo) Compare(info LogInfo) bool {
	if !format.Valid() {
		return false
	}
	if !info.Valid() {
		return false
	}
	if format.Kind != info.Kind {
		return false
	}
	for i := range format.Elements {
		if format.Elements[i] != info.Elements[i] {
			return false
		}
	}
	return true
}
