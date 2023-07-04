package metrics

import "testing"
import "strconv"

func BenchmarkMetadata(b *testing.B) {
	m := NewMetadata()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.AddOrUpdate("key" + strconv.Itoa(i), "value")
	}
}

func TestMetadataClear(t *testing.T) {
	m := NewMetadata()
	m.AddOrUpdate("key", "value")
	m.Clear()
	if tags := m.GetAll(); len(tags) != 0 {
		t.Errorf("m.GetAll(): 0 != %v\n", len(tags))
	}
}

func TestMetadataAddOrUpdate(t *testing.T) {
	m := NewMetadata()
	m.AddOrUpdate("key", "value")
	m.AddOrUpdate("key1", "value")
	m.AddOrUpdate("key2", "value")
	m.AddOrUpdate("key3", "value")
	m.AddOrUpdate("key4", "value")
	if tags := m.GetAll(); len(tags) != 5 {
		t.Errorf("m.GetAll(): 5 != %v\n", len(tags))
	}
}

func TestMetadataRemove(t *testing.T) {
	m := NewMetadata()
	m.AddOrUpdate("key", "value")
	m.Remove("key")
	if tags := m.GetAll(); len(tags) != 0 {
		t.Errorf("m.GetAll(): 0 != %v\n", len(tags))
	}
}
