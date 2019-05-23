package dlock

type DLockInf interface {
	Lock()
	UnLock()
}
