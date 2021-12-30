Icinga2 servers in the wild often serve self-signed certificates which
fail verification by Go's tls client. To ignore the errors, Dial the server
with a modified http.Client:

	t := http.DefaultTransport.(*http.Transport)
	t.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	c := http.DefaultClient
	c.Transport = t
	client, err := icinga.Dial(host, user, pass, c)
	if err != nil {
		// handle error
	}
	...

## Why?

The terraform provider...
