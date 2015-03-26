package stress

import (
	"sync"

	"github.com/onsi/gomega/types"
)

type FailInterceptor struct {
	ginkgoFail types.GomegaFailHandler
	didFail    bool
	lock       *sync.Mutex
}

func NewFailInterceptor(fail types.GomegaFailHandler) *FailInterceptor {
	return &FailInterceptor{
		ginkgoFail: fail,
		lock:       &sync.Mutex{},
	}
}

func (f *FailInterceptor) Fail(message string, callerSkip ...int) {
	f.lock.Lock()
	f.didFail = true
	f.lock.Unlock()
	if len(callerSkip) == 0 {
		f.ginkgoFail(message, 1)
	} else {
		f.ginkgoFail(message, callerSkip[0]+1)
	}
}

func (f *FailInterceptor) Reset() {
	f.lock.Lock()
	f.didFail = false
	f.lock.Unlock()
}

func (f *FailInterceptor) DidFail() bool {
	f.lock.Lock()
	defer f.lock.Unlock()

	return f.didFail
}
