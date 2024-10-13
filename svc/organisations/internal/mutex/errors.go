package mutex

import "fmt"

type ErrLockNotClaimed struct {
	Key string
	Err error
}

func (e ErrLockNotClaimed) Error() string {
	return fmt.Sprintf("could not claim lock '%s': %s", e.Key, e.Err.Error())
}
