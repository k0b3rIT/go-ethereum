package metrics

// Counters hold an int64 value that can be incremented and decremented.
type Metadata interface {
	Clear()
	AddOrUpdate(string, string)
	Remove(string)
	GetAll() map[string]string
	Snapshot() Metadata
}

// GetOrRegisterCounter returns an existing Counter or constructs and registers
// a new StandardCounter.
func GetOrRegisterMetadata(name string, r Registry) Metadata {
	if nil == r {
		r = DefaultRegistry
	}
	return r.GetOrRegister(name, NewMetadata).(Metadata)
}

// GetOrRegisterCounterForced returns an existing Counter or constructs and registers a
// new Counter no matter the global switch is enabled or not.
// Be sure to unregister the counter from the registry once it is of no use to
// allow for garbage collection.
func GetOrRegisterMetadataForced(name string, r Registry) Metadata {
	if nil == r {
		r = DefaultRegistry
	}
	return r.GetOrRegister(name, NewMetadataForced).(Metadata)
}

// NewCounter constructs a new StandardCounter.
func NewMetadata() Metadata {
	if !Enabled {
		return NilMetadata{}
	}
	return &StandardMetadata{tags: make(map[string] string)}
}

// NewCounterForced constructs a new StandardCounter and returns it no matter if
// the global switch is enabled or not.
func NewMetadataForced() Metadata {
	return &StandardMetadata{tags: make(map[string] string)}
}

// NewRegisteredCounter constructs and registers a new StandardCounter.
func NewRegisteredMetadata(name string, r Registry) Metadata {
	c := NewMetadata()
	if nil == r {
		r = DefaultRegistry
	}
	r.Register(name, c)
	return c
}

// NewRegisteredCounterForced constructs and registers a new StandardCounter
// and launches a goroutine no matter the global switch is enabled or not.
// Be sure to unregister the counter from the registry once it is of no use to
// allow for garbage collection.
func NewRegisteredMetadataForced(name string, r Registry) Metadata {
	c := NewMetadataForced()
	if nil == r {
		r = DefaultRegistry
	}
	r.Register(name, c)
	return c
}




// CounterSnapshot is a read-only copy of another Counter.
type MetadataSnapshot struct {
	tags map[string]string
}

// Clear panics.
func (MetadataSnapshot) Clear() {
	panic("Clear called on a CounterSnapshot")
}

func (MetadataSnapshot) AddOrUpdate(key string, value string) {
	panic("AddOrUpdate called on a MetadataSnapshot")
}

func (MetadataSnapshot) Remove(key string)() {
	panic("Remove called on a MetadataSnapshot")
}

func (m MetadataSnapshot) GetAll() map[string]string { 
	result := make(map[string]string)
	for k,v := range m.tags {
		result[k] = v
	  }
	return result;
}

// Snapshot returns the snapshot.
func (m MetadataSnapshot) Snapshot() Metadata { return m }




// NilCounter is a no-op Counter.
type NilMetadata struct{}

// Clear is a no-op.
func (NilMetadata) Clear() {}

func (NilMetadata) AddOrUpdate(key string, value string) {}

func (NilMetadata) Remove(key string)() {}

func (NilMetadata) GetAll() map[string]string { return make(map[string]string)}

// Snapshot is a no-op.
func (NilMetadata) Snapshot() Metadata { return NilMetadata{} }




// StandardCounter is the standard implementation of a Counter and uses the
// sync/atomic package to manage a single int64 value.
type StandardMetadata struct {
	tags map[string]string
}

// Clear sets the counter to zero.
func (c *StandardMetadata) Clear() {
	c.tags = make(map[string]string)
}

func (c *StandardMetadata) AddOrUpdate(key string, value string) {
	c.tags[key] = value;
}

func (c *StandardMetadata) Remove(key string) {
	delete(c.tags, key)
}

func (c *StandardMetadata) GetAll() map[string]string {
	result := make(map[string]string)
	for k,v := range c.tags {
		result[k] = v
	  }
	return result;
}


// Snapshot returns a read-only copy of the counter.
func (c *StandardMetadata) Snapshot() Metadata {
	return &MetadataSnapshot{tags: c.GetAll()}
}
