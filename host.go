package icinga

type Host struct {
	Name  string
	State int
}

var ErrNoHost = errors.New("no such host")

func (c *Client) Hosts() ([]Host, error) {
	_, err := c.get("/objects/hosts")
	if err != nil {
		return nil, err
	}
	return nil, err
}

func (c *Client) LookupHost(name string) (Host, error) {
	resp, err := c.get("/objects/hosts/" + name)
	if err != nil {
		return Host{}, err
	}
	if resp.StatusCode == http.StatusNotFound {
		return Host{}, fmt.Errorf("lookup %s: %w", name, ErrNoHost)
	}
	return Host{}, err
}

func (c *Client) DeleteHost(name string) error {
	if err := c.delete("/objects/hosts/" + name); err != nil {
		return fmt.Errorf("delete host %s: %w", name, err)
	}
	return nil
}
