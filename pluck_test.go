package gfl

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type case1 struct {
	A string `json:"a"`
	B int
}

type case2 struct {
	A []string `json:"a"`
	B case1    `json:"b"`
}

type case3 struct {
	A bool     `json:"a"`
	B bool     `json:"b,omitempty"`
	C *case1   `json:"c,omitempty"`
	D []string `json:"d,omitempty"`
}

type case4 struct {
	case1
}

type case5 struct {
	case1
	A int `json:"a"`
}

func TestFieldSet(t *testing.T) {
	fs1 := fieldSet{"a": nil, "b": nil}
	assert.Equal(t, newFieldSet("a", "b"), fs1)
	fs2 := fieldSet{"a": nil, "b": fieldSet{"b1": nil}}
	assert.Equal(t, newFieldSet("a", "b.b1"), fs2)
}

func TestPluckFromStruct(t *testing.T) {
	item := case1{"a", 0}
	assert.Equal(t, Pluck(item), item)
	assert.Equal(t, Pluck(item, "a"), map[string]interface{}{"a": "a"})
	assert.Equal(t, Pluck(item, "b"), map[string]interface{}{})
	assert.Equal(t, Pluck(item, "B"), map[string]interface{}{"B": 0})
	assert.Equal(t, Pluck(&item, "a"), map[string]interface{}{"a": "a"})
}

func TestPluckFromNestedStruct(t *testing.T) {
	c1 := case1{"a", 0}
	item := case2{[]string{"a"}, c1}
	assert.Equal(t, Pluck(item, "b"), map[string]interface{}{"b": c1})
	assert.Equal(t, Pluck(item, "a", "b.a"), map[string]interface{}{
		"a": []string{"a"},
		"b": map[string]interface{}{"a": "a"},
	})
}

func TestPluckFromSlice(t *testing.T) {
	items := []case1{case1{"a1", 0}, case1{"a2", 1}}
	assert.Equal(t, Pluck(items), items)
	assert.Equal(t, Pluck(items, "a"), []interface{}{
		map[string]interface{}{"a": "a1"},
		map[string]interface{}{"a": "a2"},
	})
	assert.Equal(t, Pluck(&items, "a"), []interface{}{
		map[string]interface{}{"a": "a1"},
		map[string]interface{}{"a": "a2"},
	})
}

func TestPluckWithTags(t *testing.T) {
	item := case3{}
	assert.Equal(t, Pluck(item, "a", "b"), map[string]interface{}{"a": false})
}

func TestPluckFromAnonymous(t *testing.T) {
	c1 := case4{case1{"a", 1}}
	assert.Equal(t, Pluck(c1, "a", "b", "B"), map[string]interface{}{"a": "a", "B": 1})
	c2 := case5{case1{"a", 1}, 0}
	assert.Equal(t, Pluck(c2, "a", "B"), map[string]interface{}{"a": 0, "B": 1})
}
