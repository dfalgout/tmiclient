package tmiclient

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsin(t *testing.T) {
	values := []string{"A", "B", "C", "d-f"}
	shouldBeTrue := isin("B", values)
	assert.Equal(t, shouldBeTrue, true, "Should find value in list")

	shouldBeFalse := isin("D", values)
	assert.Equal(t, shouldBeFalse, false, "Should not find value in list")

	shouldBeFalse = isin("-", values)
	assert.Equal(t, shouldBeFalse, false, "Should not find value in list")

	shouldBeTrue = isin("D-f", values)
	assert.Equal(t, shouldBeTrue, true, "Should find value in list")
}

func TestAppendListUnique(t *testing.T) {
	list := []string{"A", "B"}
	listA := appendListUnique(list, "C")
	assert.Equal(t, listA, []string{"A", "B", "C"}, "Should append unique value to end of list")

	listB := appendListUnique(list, "A")
	assert.Equal(t, listB, []string{"A", "B"}, "Should not append a value that already exists in list")
}

func TestGetUserName(t *testing.T) {
	value := ":u_lost!u_lost@u_lost.tmi.twitch.tv JOIN #chatrooms:44322889:04e762ec-ce8f-4cbc-b6a3-ffc871ab53da"
	result := getUserName(value)
	assert.Equal(t, result, "u_lost", "Should match the the username")
}

func TestGetFieldMap(t *testing.T) {
	shouldBe := map[string]string{
		"color":        "#0D4200",
		"mod":          "0",
		"room-id":      "1337",
		"subscriber":   "0",
		"tmi-sent-ts":  "1507246572675",
		"user-type":    "global_mod",
		"badges":       "global_mod/1,turbo/1",
		"display-name": "dallas",
		"emotes":       "25:0-4,12-16/1902:6-10",
		"id":           "b34ccfc7-4977-403a-8a94-33c6bac34fb8",
		"turbo":        "1",
		"user-id":      "1337",
	}
	value := "@badges=global_mod/1,turbo/1;color=#0D4200;display-name=dallas;emotes=25:0-4,12-16/1902:6-10;id=b34ccfc7-4977-403a-8a94-33c6bac34fb8;mod=0;room-id=1337;subscriber=0;tmi-sent-ts=1507246572675;turbo=1;user-id=1337;user-type=global_mod"
	result := getFieldMap(value)
	assert.Equal(t, result, shouldBe)
}

func TestGetBoolFromString(t *testing.T) {
	value := "1"
	result := getBoolFromString(value)
	assert.Equal(t, result, true)

	value = "0"
	result = getBoolFromString(value)
	assert.Equal(t, result, false)
}

func TestGetIntFromString(t *testing.T) {
	value := "25"
	result := getIntFromString(value)
	assert.Equal(t, result, 25)
}

func TestGetIntListFromString(t *testing.T) {
	value := "1,2,3,4,"
	result := getIntListFromString(value)
	assert.Equal(t, result, []int{1, 2, 3, 4})
}
