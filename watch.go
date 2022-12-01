package icinga

// A ServiceWatcher continuously queries Icinga about a set of Services.
// A ServiceWatcher uses the next check time of Services to determine when to query Icinga again.
// If it cannot determine the next check time, the Watcher waits 1 minute before trying again.
type ServiceWatcher struct {
	kill chan bool
	Updates chan ServiceMsg
}

type HostWatcher struct {
	Kill chan bool
	Updates chan HostMsg
}

type ServiceMsg struct {
	Services []Service
	Err error
}

type HostMsg struct {
	Hosts []Host
	Err error
}

// Kill stops the Watcher irrevocably.
// A new Watcher must be created with WatchServices.
func (w *ServiceWatcher) Kill() {
	w.kill <- true
}

// WatchServices returns a new Watcher which uses client to
// continuously query Icinga for Services matching the given filter
// expression.
func WatchServices(client *Client, expr string) *Watcher {
	timer := time.NewTimer(0)
	kill := make(chan bool)
	ch := make(chan ServiceMsg)
	go func() {
		select {
		case <-kill:
			close(ch)
			return
		case <-timer.C:
			services, err := client.Services(expr)
			ch <- ServiceMsg{services, err}
			if err != nil {
				timer.Reset(1*time.Minute)
			} else {
				timer.Reset(time.Until(nextCheck(services)))
			}
		}
	}
	return &Watcher{
		kill: kill
		Updates: ch
	}
}

func WatchHosts(expr string, ch chan HostMsg, stop chan bool) {
	timer := time.NewTimer(0)
	select {
	case <-stop:
		close(ch)
		return
	case <-timer.C:
		hosts, err := client.Hosts(expr)
		ch <- HostMsg{hosts, err}
		timer.Reset(time.Until(NextCheck(hosts)))
	}
}
