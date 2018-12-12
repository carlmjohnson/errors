package errors_test

import (
	"fmt"

	"github.com/carlmjohnson/errors"
)

type closer struct{}

func (c closer) Close() error {
	return fmt.Errorf("<had problem closing!>")
}

func openThingie() (c closer, err error) { return }

// Calling Close() on an io.WriteCloser can return an important error
// encountered while flushing to disk. Don't risk missing them by
// using a plain defer w.Close(). Use errors.Defer to capture the return value.
func ExampleDefer() {
	// If you just defer a close call, you can miss an error
	return1 := func() error {
		thing, err := openThingie()
		if err != nil {
			return err
		}
		defer thing.Close() // oh no, this returned an error!
		// do stuff...
		return err
	}()
	fmt.Println("return1 ==", return1)

	// Use errors.Defer and a named return to capture the error
	return2 := func() (err error) {
		thing, err := openThingie()
		if err != nil {
			return err
		}
		defer errors.Defer(&err, thing.Close)
		// do stuff...
		return err
	}()
	fmt.Println("return2 ==", return2)

	// Output:
	// return1 == <nil>
	// return2 == <had problem closing!>
}
