package icinga

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// An Event represents an event from the Icinga Event Stream.
type Event struct {
	// Type indicates the type of the stream, such as CheckResult.
	Type string `json:"type"`
	// Host is the name of an Icinga Host object which this event relates to.
	Host string `json:"host"`
	// Service is the name of an Icinga Service object which this event relates to.
	// It is empty when a CheckResult event of a Host object is received.
	Service         string       `json:"service"`
	Acknowledgement bool         `json:"acknowledgement"`
	CheckResult     *CheckResult `json:"check_result"`
	Error           error
}

// Subscribe returns a channel through which events from the
// corresponding Icinga Event Stream named in typ are sent.
// Queue is a unique identifier Icinga uses to manage stream clients.
// Filter is a filter expression which modifies which events will be received;
// the empty string means all events are sent.
//
// Any errors on initialising the connection are returned immediately as a value.
// Subsequent errors reading the stream are set in the Error field of sent Events.
// Callers should handle both cases and resubscribe as required.
func (c *Client) Subscribe(typ, queue, filter string) (<-chan Event, error) {
	m := map[string]interface{}{
		"types":  []string{typ},
		"queue":  queue,
		"filter": filter,
	}
	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(m); err != nil {
		return nil, fmt.Errorf("encode stream parameters: %w", err)
	}
	resp, err := c.post("/events", buf)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		iresp, err := parseResponse(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("request events: parse error response: %w", err)
		}
		return nil, fmt.Errorf("request events: %w", iresp.Error)
	}
	sc := bufio.NewScanner(resp.Body)
	ch := make(chan Event)
	go func() {
		for sc.Scan() {
			var ev Event
			if err := json.Unmarshal(sc.Bytes(), &ev); err != nil {
				ch <- Event{Error: fmt.Errorf("decode event: %v", err)}
				continue
			}
			ch <- ev
		}
		if sc.Err() != nil {
			ch <- Event{Error: fmt.Errorf("scan response: %w", sc.Err())}
		}
		resp.Body.Close()
		close(ch)
	}()
	return ch, nil
}
