package cache

import (
	"testing"
	"time"
)

func TestNewStdMapCache(t *testing.T) {
	var c Cache
	c = NewStdMapCache(1)

	testCache(c, t)
}

func TestNewSyncMapCache(t *testing.T) {
	var c Cache
	c = NewSyncMapCache(1)

	testCache(c, t)
}

func testCache(c Cache, t *testing.T) {
	t.Parallel()
	if c.Size() != 0 {
		t.Error("size not 0")
	}

	c.Set("key1", "value1")
	c.SetWithTTl("key2", "value2", 3)

	if c.Get("key1") != "value1" {
		t.Error("cannot get value1")
	}
	if c.Get("key2") != "value2" {
		t.Error("cannot get value1")
	}
	if c.Size() != 2 {
		t.Error(`size is not 2`)
	}

	time.Sleep(time.Duration(2) * time.Second)
	if c.Get("key1") != nil {
		t.Error(`got key1 after 2 sec`)
	}
	if c.Get("key2") != "value2" {
		t.Error(`didn't get key2`)
	}
	if c.Size() != 2 {
		t.Error(`size is not 2`)
	}

	time.Sleep(time.Duration(2) * time.Second)
	if c.Get("key1") != nil {
		t.Error(`got key1 after 4 sec`)
	}
	if c.Get("key2") != nil {
		t.Error(`got key2 after 4 sec`)
	}
	if c.Size() != 2 {
		t.Error(`size is not 2`)
	}

	c.Del("key1")
	c.Del("key2")
	if c.Size() != 0 {
		t.Error(`size is not 0`)
	}
}
