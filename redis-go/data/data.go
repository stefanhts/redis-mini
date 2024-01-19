package data

import (
	"fmt"
)

type Element struct {
	Key   string
	Value string
}

var Store []Element

func initStore() {
	Store = []Element{}
}

func Push(key, val string) {
	Store = append(Store, Element{
		Key:   key,
		Value: val,
	})
}

func Pop() (Element, error) {
	if len(Store) < 1 {
		return Element{}, fmt.Errorf("Pop error: no elements in list")
	}
	el := Store[len(Store)-1]
	Store = Store[:len(Store)-1]
	return el, nil
}

func LLen() int {
	return len(Store)
}

func Get(key string) (string, error) {
	for _, el := range Store {
		if el.Key == key {
			return el.Value, nil
		}
	}
	return "", fmt.Errorf("Get error: key not found %s", key)
}

func GetInd(ind int) (Element, error) {
	if ind >= len(Store) {
		return Element{}, fmt.Errorf("GetInd: index %d is out of bounds %d", ind, len(Store))
	}
	return Store[ind], nil
}
