# keylock

[![Build Status](https://travis-ci.com/mountain-viewer/keylock.svg?branch=master)](https://travis-ci.com/mountain-viewer/keylock)

A Golang utility package which provides a way to lock a set of keys so that any subsequent call with a non-empty subset of already locked keys makes the goroutine waiting.

# Usage

```
l := keylock.New()

canceled, unlock := l.LockKeys([]string{"a", "b"}, nil)  // locks {"a", "b"}
defer unlock()

go func() {
    _, _ = l.LockKeys([]string{"b", "c"}, nil)  // can't lock {"b", "c"} 
}()                                             // since "b" has been locked already

time.Sleep(time.Millisecond * 10)
canceled, _ = l.LockKeys([]string{"d"}, nil)  // locks {"d"}
```