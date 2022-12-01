package main

type Cache struct {
	client *icinga.Client
	hosts map[string]*hostEntry
	services map[string]*serviceEntry
}

func NewCache(client *icinga.Client) *Cache {
	return &Cache{
		client: client,
		hosts: make(map[string]*hostEntry),
		services: make(map[string]*serviceEntry),
	}
}

type serviceEntry struct {
	accessTime time.Time
	watcher *ServiceWatcher
	services []Service
}

func (c *Cache) GetServices(expr string) []icinga.Service {
	entry, ok := c.services[expr]
	if !ok {
		return nil
	}
	entry.accessTime = time.Now()
	if entry.wacher == nil {
		entry := entry
		entry.watcher = WatchServices(c.client, expr)
		go func() {
			for msg := range entry.watcher.Updates {
				if msg.Err != nil {
					// TODO(otl) ??
				}
				entry.services = msg.Services
			}
		}()
	}
	return entry.services
}

type hostEntry struct {
	accessTime time.Time
	hosts []Host
}

func (c *Cache) GetHosts(expr string) *hostEntry {}

func (c *Cache) evictOlderThan(d time.Duration) {
	t := time.Now().Sub(d)
	for _, entry := range c.services {
		if entry.accessTime.Before(t) {
			entry.kill <- true
			delete(c.services[expr])
		}
	}
	for _, entry := range c.Hosts {
		if entry.accessTime.Before(t) {
			entry.kill <- true
			delete(c.hosts[expr])
		}
	}
}
