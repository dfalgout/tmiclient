package tmiclient

import (
	"log"
	"strconv"
	"strings"
)

func isin(value string, values []string) bool {
	value = strings.ToLower(value)
	for _, v := range values {
		if value == strings.ToLower(v) {
			return true
		}
	}

	return false
}

func appendListUnique(list []string, value string) []string {
	for _, v := range list {
		if value == v {
			return list
		}
	}

	list = append(list, value)
	return list
}

func handleWriteError(err error) {
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func getUserName(id string) string {
	channel := strings.Split(id, "!")[0]
	return strings.Trim(channel, ":")
}

func getFieldMap(msg string) map[string]string {
	fieldMap := make(map[string]string)
	messageFields := strings.Split(strings.TrimPrefix(msg, "@"), ";")

	for _, v := range messageFields {
		kv := strings.Split(v, "=")
		fieldMap[kv[0]] = kv[1]
	}

	return fieldMap
}

func getBoolFromString(str string) bool {
	if str == "0" {
		return false
	}
	return true
}

func getIntFromString(str string) int {
	i, err := strconv.Atoi(str)
	if err != nil {
		log.Println(err)
	}

	return i
}

func getIntListFromString(str string) []int {
	var intItems []int

	items := strings.Split(str, ",")
	for _, v := range items {
		i, err := strconv.Atoi(v)
		if err != nil {
			log.Println(err)
			return intItems
		}
		intItems = append(intItems, i)
	}

	return intItems
}
