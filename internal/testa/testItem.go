package testa

import "mangosteena/pkg/info"

type Method struct {
	//상태 들
	resultType  string
	message     string
	name        string
	maxCount    uint32
	match_count uint32
	oders       []Method_order
	historys    []history
}
type history struct {
	resultType  string
	message     string
	maxCount    uint32
	match_count uint32
	info        info.LogInfo
	History     map[string]info.LogInfo
}
type Method_order struct {
	status_tests Conditional
	cond         Conditional
}

type Conditional struct {
	And   []Conditional
	Or    []Conditional
	Not   []Conditional
	Value info.LogInfo
}

func (method *Method) valid(last map[string]info.LogInfo, info info.LogInfo) bool {
	return compareAction(info, method.oders[0].cond) &&
		compareState(last, method.oders[0].status_tests)
}

func (method *Method) addHistory(last map[string]info.LogInfo, info info.LogInfo) {
	method.historys = append(method.historys, history{
		resultType:  method.resultType,
		message:     method.message,
		maxCount:    method.maxCount,
		match_count: method.match_count,
		info:        info,
		History:     last,
	})
}

func compareAction(info info.LogInfo, cond Conditional) bool {
	if len(cond.And) != 0 {
		for i := range cond.And {
			if !compareAction(info, cond.And[i]) {
				return false
			}
		}
		return true
	} else if len(cond.Or) != 0 {
		for i := range cond.Or {
			if compareAction(info, cond.Or[i]) {
				return true
			}
		}
		return false
	} else if len(cond.Not) != 0 {
		for i := range cond.Not {
			if compareAction(info, cond.Not[i]) {
				return false
			}
		}
		return true
	} else {
		if !cond.Value.Valid() {
			return false
		}

		return cond.Value.Compare(info)
	}
}

func compareState(last map[string]info.LogInfo, cond Conditional) bool {
	if len(cond.And) != 0 {
		for i := range cond.And {
			if !compareState(last, cond.And[i]) {
				return false
			}
		}
		return true
	} else if len(cond.Or) != 0 {
		for i := range cond.Or {
			if compareState(last, cond.Or[i]) {
				return true
			}
		}
		return false
	} else if len(cond.Not) != 0 {
		for i := range cond.Not {
			if compareState(last, cond.Not[i]) {
				return false
			}
		}
		return true
	} else {
		if !cond.Value.Valid() {
			return false
		}

		if kkk, ok := last[cond.Value.Kind]; ok {
			return cond.Value.Compare(kkk)
		}
		return false
	}
}
