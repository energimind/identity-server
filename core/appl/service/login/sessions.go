package login

import "sync"

// sessions is a session cache.
type sessions struct {
	entries map[string]*session
	mx      sync.RWMutex
}

// newSessions creates a new cache.
func newSessions() *sessions {
	return &sessions{
		entries: make(map[string]*session),
	}
}

// put puts a session into the cache.
func (c *sessions) put(id string, s *session) {
	c.mx.Lock()
	defer c.mx.Unlock()

	c.entries[id] = s
}

// get gets a session from the cache.
func (c *sessions) get(id string) (*session, bool) {
	c.mx.RLock()
	defer c.mx.RUnlock()

	s, ok := c.entries[id]

	return s, ok
}

// delete deletes a session from the cache.
func (c *sessions) delete(id string) {
	c.mx.Lock()
	defer c.mx.Unlock()

	delete(c.entries, id)
}
