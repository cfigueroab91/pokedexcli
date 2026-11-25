package pokecache

import (
	"bytes"
	"testing"
	"time"
)

func TestAddGet(t *testing.T) {
	const interval = 5 * time.Second

	cases := []struct {
		key string
		val []byte
	}{
		{key: "https://example.com", val: []byte("testdata")},
		{key: "https://example.com/path", val: []byte("moretestdata")},
		{key: "empty", val: []byte{}},
	}

	cache := NewCache(interval)

	for _, c := range cases {
		cache.Add(c.key, c.val)
	}

	// Todas las claves deben seguir presentes: ninguna sobrescribe a otra.
	for _, c := range cases {
		val, ok := cache.Get(c.key)
		if !ok {
			t.Errorf("expected to find key %q", c.key)
			continue
		}
		if !bytes.Equal(val, c.val) {
			t.Errorf("key %q: expected %q, got %q", c.key, c.val, val)
		}
	}
}

func TestGetMissingKey(t *testing.T) {
	cache := NewCache(5 * time.Second)
	cache.Add("exists", []byte("value"))

	val, ok := cache.Get("missing")
	if ok {
		t.Errorf("expected ok to be false for a missing key")
	}
	if val != nil {
		t.Errorf("expected nil value for a missing key, got %q", val)
	}
}

func TestAddReplacesValue(t *testing.T) {
	cache := NewCache(5 * time.Second)

	cache.Add("key", []byte("first"))
	cache.Add("other", []byte("untouched"))
	cache.Add("key", []byte("second"))

	val, ok := cache.Get("key")
	if !ok {
		t.Fatalf("expected to find key")
	}
	if string(val) != "second" {
		t.Errorf("expected %q, got %q", "second", val)
	}

	other, ok := cache.Get("other")
	if !ok || string(other) != "untouched" {
		t.Errorf("replacing one key changed another: got %q, %v", other, ok)
	}
}

func TestFreshEntryStaysAvailable(t *testing.T) {
	cache := NewCache(time.Hour)
	cache.Add("key", []byte("value"))

	// Con un intervalo largo, el ticker no llega a disparar y la entrada sigue.
	time.Sleep(10 * time.Millisecond)

	val, ok := cache.Get("key")
	if !ok {
		t.Fatalf("expected fresh entry to still be in the cache")
	}
	if string(val) != "value" {
		t.Errorf("expected %q, got %q", "value", val)
	}
}

func TestReapRemovesOldEntries(t *testing.T) {
	const interval = 10 * time.Millisecond

	cache := NewCache(interval)
	cache.Add("key", []byte("value"))

	if _, ok := cache.Get("key"); !ok {
		t.Fatalf("expected entry to exist right after adding it")
	}

	// La limpieza corre en otra goroutine, así que damos varios ticks de margen
	// en vez de comprobar en un instante exacto.
	deadline := time.Now().Add(20 * interval)
	for time.Now().Before(deadline) {
		if _, ok := cache.Get("key"); !ok {
			return
		}
		time.Sleep(interval / 2)
	}

	t.Errorf("expected entry to be reaped after %v", 20*interval)
}

func TestReapKeepsNewerEntries(t *testing.T) {
	cache := NewCache(time.Hour)

	cache.Add("old", []byte("old"))
	cache.Add("new", []byte("new"))

	// Llamamos a reap directamente con un "ahora" simulado, sin depender del reloj.
	cache.mu.Lock()
	entry := cache.entries["old"]
	entry.createdAt = time.Now().Add(-2 * time.Hour)
	cache.entries["old"] = entry
	cache.mu.Unlock()

	cache.reap(time.Now(), time.Hour)

	if _, ok := cache.Get("old"); ok {
		t.Errorf("expected old entry to be removed")
	}
	if _, ok := cache.Get("new"); !ok {
		t.Errorf("expected new entry to be kept")
	}
}
