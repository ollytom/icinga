package icinga provides a client to the Icinga2 HTTP API.

[![godocs.io](http://godocs.io/olowe.co/icinga?status.svg)](http://godocs.io/olowe.co/icinga)

Send patches, questions or a friendly "hello" to the mailing list: [~otl/icinga@lists.sr.ht](mailto:~otl/icinga@lists.sr.ht) 🙂
Or, read [the archives][list].

## Quick Start

See the [package overview godoc][godocs] for getting started examples.

[godocs]: https://godocs.io/olowe.co/icinga

## Development

Some code is automatically generated. Ensure it's up-to-date before starting work:

	go generate

Make some changes, then run the tests:

	go test

Please send any patches to the [mailing list](https://lists.sr.ht/~otl/icinga):

	git send-email --to="~otl/icinga@lists.sr.ht" HEAD^

For those unfamiliar with this workflow, see [git-send-email.io][sendemail].

[list]: https://lists.sr.ht/~otl/icinga
[sendemail]: https://git-send-email.io

### Tests

Some tests dial an instance of Icinga2 running on the loopback address
and the standard Icinga2 port 5665 ("127.0.0.1:5665"). If this fails,
those tests are skipped. To run these tests, create the following API
user:

	object ApiUser "root" {
		password = "icinga"
		permissions = [ "*" ]
	}

Getting data from 127.0.0.1:5665 to an Icinga server is left as an
exercise to the reader!

Personally, I run an Alpine Linux virtual machine using qemu. You
could also use the [official Icinga2 container image][image].

[image]: https://hub.docker.com/r/icinga/icinga2

### Code generation

Source code for the basic lookup, create and delete operations of some
Icinga2 object types, such as Host and Service, are generated
automatically.

To generate the code, ensure the following tools are available:

* POSIX shell (/bin/sh)
* awk
* gofmt

The shell script crud.sh writes Go source code by reading a template
file and doing some text substitution. It loops through object types,
piping the template file crud.skel into the awk script crud.awk for
each.

crud.sh writes code to the standard output by default:

	./crud.sh

If the flag `-o` is set, code will be written to the file
specified instead of to standard output:

	./crud.sh -o crud.go

Code generation is used because the functions are trivial and call the exact
same underlying methods on Client anyway. The only thing that differs is the type.
Perhaps when Go gets type parameters then this will go away?

## Why Another Package?

The [icinga2 terraform provider][tf] uses the package [github.com/lrsmith/go-icinga2-api/iapi][lrsmith].
As I read the source code I felt I wasn't reading idiomatic Go as detailed in documents like [Effective Go][effectivego].
Other properties of `iapi` felt unusual to me:

* The client to the API has the  confusing name `server`.
* Every HTTP request creates a new http.Client.
* Types have superfluous names like `HostStruct` instead of just `Host`.
* Every response body from the API is decoded from JSON into one data strucutre, marshalled into JSON again, then unmarshalled back into another.
* Every error returned from a function has a new name, rather than reusing the idiomatic name `err`.

If I was being paid, I'd create a fork and contribute patches upstream to carefully avoid breaking functionality of existing users of `iapi`.

But I'm not being paid ;)

[effectivego]: https://go.dev/doc/effective_go
[tf]: https://registry.terraform.io/providers/Icinga/icinga2/latest
[lrsmith]: https://godocs.io/github.com/lrsmith/go-icinga2-api/iapi
