package utils

import (
	"encoding/json"
	"fmt"
	"sync"
)

func MapGroupingBy[K comparable, T any](vs []T, f func(T) K) map[K][]T {
	group := make(map[K][]T)
	for _, v := range vs {
		k := f(v)
		group[k] = append(group[k], v)
	}
	return group
}

// SliceToKeyMap transfer []K to map[K]struct{}
func SliceToKeyMap[K comparable](s []K) map[K]struct{} {
	m := make(map[K]struct{})
	for _, k := range s {
		m[k] = struct{}{}
	}
	return m
}

// SubtractMap return Map witch A exist and B not exist
func SubtractMap[M1 ~map[K]V1, M2 ~map[K]V2, K comparable, V1, V2 any](aMap M1, bMap M2) M1 {
	m := make(M1)
	for aK, aV := range aMap {
		_, ok := bMap[aK]
		if !ok {
			m[aK] = aV
		}
	}
	return m
}

// UnionMap return Map witch base on A and plus B
func UnionMap[M ~map[K]struct{}, K comparable](aMap M, bMap M) M {
	m := make(M)
	for aK, aV := range aMap {
		m[aK] = aV
	}
	for bK, bV := range bMap {
		m[bK] = bV
	}
	return m
}

func ToMapStringAny(o any) map[string]any {
	m := make(map[string]any)
	b, _ := json.Marshal(o)
	json.Unmarshal(b, &m)
	return m
}

func MapStringAnyToMapStringString(mOrigin map[string]any) map[string]string {
	m := make(map[string]string)
	for k, v := range mOrigin {
		m[k] = v.(string)
	}
	return m
}

func MapToSyncMap[K comparable, T any](srcMap map[K]T) sync.Map {
	var res sync.Map

	for k, v := range srcMap {
		res.Store(k, v)
	}

	return res
}

func LenSyncMap(m *sync.Map) int {
	var count int
	m.Range(func(key, value any) bool {
		count++
		return true
	})

	return count
}

// PrintSynMap print key, value in json format
func PrintSynMap(m *sync.Map) string {
	mapString := "{"

	m.Range(func(key, value any) bool {
		keyStr := fmt.Sprintf("%v", key)
		valueStr := fmt.Sprintf("%v", value)
		mapString = mapString + "\"" + keyStr + "\"" + ":" + "\"" + valueStr + "\"" + ","

		return true
	})
	mapString = mapString[:len(mapString)-1] + "}"

	return mapString
}
