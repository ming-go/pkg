package dmutex

type Locker interface {
	Lock(string) error
	Unlock(string) error
}
