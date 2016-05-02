package core

import "fmt"

func ConvertConfigUpstreamToMap(us []*Upstream) (usMap map[string]*Upstream) {

	usMap = make(map[string]*Upstream)

	for _, v := range us {
		usMap[CreateUpstreamKey(v.Target, v.Port)] = v
	}

	return usMap
}

func ConvertConfigUpstreamFromMap(usMap map[string]*Upstream) (us []*Upstream) {

	us = []*Upstream{}

	for _, v := range usMap {
		us = append(us, v)
	}
	fmt.Print(us)
	return
}
