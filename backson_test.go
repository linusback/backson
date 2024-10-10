package backson

import (
	"encoding/json"
	"slices"
	"testing"
)

type valueArrayTest[T comparable] struct {
	json     []byte
	expected []T
	ch       chan T
}

func (at *valueArrayTest[T]) test(t *testing.T) {
	err := ParseArray(at.json, at.ch)
	if err != nil {
		t.Error(err)
	}
	checkSame(t, at.ch, at.expected, at.json)
}

func Test_ParseArray(t *testing.T) {
	intTest := makeTest(t, []int{1, 2, 3})
	intTest.test(t)
	stringTest := makeTest(t, []string{"hello", "my", "name", "is"})
	stringTest.test(t)
	byteTest := makeTest(t, []uint8{3, 2, 1})
	byteTest.json = []byte("[3,2,1]")
	byteTest.test(t)
}

func checkSame[T comparable](t *testing.T, ch chan T, expected []T, js []byte) {
	result := make([]T, 0, len(expected))
	for val := range ch {
		result = append(result, val)
	}
	if !slices.Equal(result, expected) {
		t.Errorf("for json %s, expected: %v got %v", js, expected, result)
	}
	return
}

func makeTest[T comparable](t *testing.T, result []T) valueArrayTest[T] {
	b, err := json.Marshal(result)
	if err != nil {
		t.Error(err)
	}
	return valueArrayTest[T]{
		json:     b,
		expected: result,
		ch:       make(chan T, len(result)),
	}
}
