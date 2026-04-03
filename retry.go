package brickognize

import "time"

func retry(attempts int, fn func() error) error {
	var err error
	backoff := 500 * time.Millisecond

	for i := 0; i < attempts; i++ {
		err = fn()
		if err == nil {
			return nil
		}
		time.Sleep(backoff)
		backoff *= 2
	}
	return err
}
