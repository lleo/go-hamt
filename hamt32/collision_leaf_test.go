package hamt32

import (
	"log"
	"testing"

	"github.com/lleo/go-hamt/key"
	"github.com/lleo/go-hamt/stringkey"
)

func TestCollisonLeafDelImmutable(t *testing.T) {
	var a = stringkey.New("aaa")
	var b = stringkey.New("bbb")
	var c = stringkey.New("ccc")
	var kvs = []key.KeyVal{{Key: a, Val: 1}, {Key: b, Val: 2}, {Key: c, Val: 3}}
	var leaf = newCollisionLeaf(kvs)
	var newLeaf, val, deleted = leaf.del(b)
	if !deleted {
		t.Fatalf("failed to delete key b=%s", b)
	}
	if val != 2 {
		t.Fatalf("val,%d != 2 from leaf.del(%s)", val, b)
	}
	if len(newLeaf.keyVals()) != 2 {
		for i, kv := range newLeaf.keyVals() {
			log.Printf("TestCollisonLeafDelImmutable: %d: %s", i, kv)
		}
		log.Printf("newLeaf => %s", newLeaf)
		t.Fatalf("len(newLeaf.keyVals()),%d != 2", len(newLeaf.keyVals()))
	}
	if len(kvs) != 3 {
		t.Fatalf("len(kvs),%d != 3", len(kvs))
	}
}
