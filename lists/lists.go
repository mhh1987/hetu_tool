package lists

import (
	"slices"

	"code.chenji.com/cj/scan_paper/utils/valid"
)

// DiffList 获取两个数组的差集，sourceList，targetList 差集， 即为sourceList中包含，但是targetList中不包含的。
func DiffList[T int8 | int16 | int | int32 | int64](sourceList, targetList []T) []T {

	if valid.IsEmptySlice(sourceList) {
		return sourceList
	}
	if valid.IsEmptySlice(targetList) {
		return sourceList
	}
	resultList := make([]T, 0, len(sourceList))
	targetMap := make(map[T]interface{}, len(targetList))
	for _, t := range targetList {
		targetMap[t] = nil
	}
	for _, s := range sourceList {
		_, ok := targetMap[s]
		if !ok {
			resultList = append(resultList, s)
		}
	}
	return resultList
}

// Map 集合元素的映射
func Map[T any, K any](dataList []T, fc func(item T) K) []K {
	if valid.IsEmptySlice(dataList) {
		return nil
	}
	resultList := make([]K, 0, len(dataList))
	for _, item := range dataList {
		resultList = append(resultList, fc(item))
	}
	return resultList
}

// ToValueList 将map的value用集合收集
func ToValueList[K string | int8 | int16 | int | int32 | int64, V any](dataMap map[K]V) []V {
	resultList := make([]V, 0, len(dataMap))
	for _, item := range dataMap {
		resultList = append(resultList, item)
	}
	return resultList
}

func ToValueListAndMerge[K string | int8 | int16 | int | int32 | int64, V any](dataMap map[K][]V) []V {

	resultList := make([]V, 0, len(dataMap))
	for _, item := range dataMap {
		resultList = append(resultList, item...)
	}
	return resultList
}

// ToKeyList 将map的key用集合收集
func ToKeyList[K string | int8 | int16 | int | int32 | int64](dataMap map[K]any) []K {

	resultList := make([]K, 0, len(dataMap))
	for key := range dataMap {
		resultList = append(resultList, key)
	}
	return resultList
}

// SliceDistinct 数组元素去重
func SliceDistinct[K string | int8 | int16 | int | int32 | int64](dataList []K) []K {

	dataMap := make(map[K]interface{}, len(dataList))
	for _, data := range dataList {
		dataMap[data] = nil
	}
	return ToKeyList(dataMap)
}

// Filter 按照条件过滤出符合条件的数据
func Filter[T any](dataList []T, fc func(item T) bool) []T {

	resultList := make([]T, 0, len(dataList))
	for _, data := range dataList {
		if fc(data) {
			resultList = append(resultList, data)
		}
	}
	return resultList
}

// FindFirst 找到符合条件的第一个
func FindFirst[T any](dataList []T, fc func(item T) bool) (T, bool) {

	var result T
	index := slices.IndexFunc(dataList, fc)
	if index == -1 {
		return result, false
	}
	return dataList[index], true
}

func IsContain[T string | int8 | int16 | int | int32 | int64](source []T, target T) bool {

	if valid.IsEmptySlice(source) {
		return false
	}
	for _, s := range source {
		if s == target {
			return true
		}
	}

	return false
}

func MergeMap[K string | int8 | int16 | int | int32 | int64, V any](sourceMap map[K]V, targetMap map[K]V) map[K]V {

	if sourceMap == nil || len(sourceMap) <= 0 {
		return targetMap
	}
	if targetMap == nil || len(targetMap) <= 0 {
		return sourceMap
	}

	for k, v := range targetMap {
		if _, ok := sourceMap[k]; !ok {
			sourceMap[k] = v
		}
	}
	return sourceMap
}

func Partition[T any](dataList []T, itemLen int) [][]T {

	if len(dataList) <= 0 {
		return nil
	}
	if itemLen <= 0 {
		return nil
	}
	resultList := make([][]T, 0, len(dataList)/itemLen+1)
	for i := 0; i < len(dataList); i += itemLen {
		end := i + itemLen
		if end > len(dataList) {
			end = len(dataList)
		}
		resultList = append(resultList, dataList[i:end])
	}
	return resultList
}

// SumValue 数组元素求和
func SumValue[T int8 | int16 | int | int32 | int64 | float32 | float64](dataList []T) T {
	var result T
	for _, data := range dataList {
		result += data
	}
	return result
}

// FindMax 找到数组中的最大值
func FindMax[T int8 | int16 | int | int32 | int64 | float32 | float64](dataList []T) T {

	var result T
	for _, data := range dataList {
		if data > result {
			result = data
		}
	}
	return result
}
