package keylock

import (
	"sort"
	"sync"
)

type KeyLock struct {
	mtx   sync.Mutex
	locks map[string]chan struct{}
}

func NewKeyLock() *KeyLock {
	return &KeyLock{
		locks: make(map[string]chan struct{}),
	}
}

func (l *KeyLock) LockKeys(keys []string, cancel <-chan struct{}) (canceled bool, unlock func()) {
	sort.Strings(keys)
	acquired := make([]chan struct{}, 0)

	unlock = func() {
		l.mtx.Lock()
		for key := range l.locks {
			delete(l.locks, key)
		}
		l.mtx.Unlock()

		for _, lock := range acquired {
			close(lock)
		}
	}

	for _, key := range keys {
		for {
			l.mtx.Lock()
			otherLock, alreadyLocked := l.locks[key]
			lock := make(chan struct{})

			if !alreadyLocked {
				l.locks[key] = lock
			}
			l.mtx.Unlock()

			if !alreadyLocked {
				acquired = append(acquired, lock)
				break
			}

			select {
			case <-cancel:
				unlock()
				canceled = true
				return
			case <-otherLock:
				continue
			}
		}
	}

	canceled = false
	return
}
