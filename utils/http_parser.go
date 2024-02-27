package utils

import (
	"strings"

	"golang.org/x/net/html"
)

func AttrEquals(attr []html.Attribute, key string, value string) bool {
	val, found := AttrContains(attr, key)
	if !found {
		return false
	}

	return strings.EqualFold(val, value)
}

func AttrContains(attr []html.Attribute, key string) (string, bool) {
	for _, attrKey := range attr {
		if attrKey.Key == key {
			return attrKey.Val, true
		}
	}

	return "", false
}
