package cache

import (
	"testing"
	"time"
)

func TestNewVacuumedCache(t *testing.T) {
	t.Parallel()

	var c Cache = NewVacuumedStdMapCache(1)

	c.Set("key1", "value1")
	c.SetWithTTl("key2", "value2", 3)

	if c.Get("key1") != "value1" {
		t.Error(`didn't get key1`)
	}
	if c.Get("key2") != "value2" {
		t.Error(`didn't get key2`)
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
	if c.Size() != 1 {
		t.Error(`size is not 1`)
	}

	time.Sleep(time.Duration(2) * time.Second)
	if c.Get("key1") != nil {
		t.Error(`got key1 after 4 sec`)
	}
	if c.Get("key2") != nil {
		t.Error(`got key2 after 4 sec`)
	}
	if c.Size() != 0 {
		t.Error(`size is not 0`)
	}
}
