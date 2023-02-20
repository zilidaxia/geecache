package consistendhash

import (
	"strconv"
	"testing"
)

func TestHashing(t *testing.T) {
	//这里使用自定义的hash算法
	hash := New(3, func(key []byte) uint32 {
		i,_ := strconv.Atoi(string(key))
		return uint32(i)
	})
	hash.Add("6","4","2")
	testCase :=map[string]string{
		"2": "2",
		"11":"2",
		"23":"4",
	}
	for k,v := range testCase {
		if hash.Get(k) != v {
			t.Errorf("Asking for %s, should have yielded %s", k, v)
		}
	}
	hash.Add("8")
	testCase["27"] = "8"
	for k,v :=range testCase {
		if hash.Get(k) != v {
			t.Errorf("Asking for %s, should have yielded %s", k, v)
		}
	}
}