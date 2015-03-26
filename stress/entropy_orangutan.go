package stress

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/onsi/composition-ginkgo-example/helpers"
	. "github.com/onsi/ginkgo"
)

type EntropyOrangutan struct {
	Client *helpers.KeyValueStoreClient
	r      *rand.Rand
}

func NewEntropyOrangutan(client *helpers.KeyValueStoreClient) *EntropyOrangutan {
	return &EntropyOrangutan{
		Client: client,
		r:      rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (e *EntropyOrangutan) MakeAMess() {
	go func() {
		t := time.NewTicker(100 * time.Millisecond)
		for {
			<-t.C
			e.doSomethingTerrible()
		}
	}()
}

func (e *EntropyOrangutan) doSomethingTerrible() {
	n := e.r.Intn(100)
	switch {
	case n < 40:
		e.heavyReadLoad()
	case n < 80:
		e.heavyWriteLoad()
	default:
		e.deleteEverything()
	}
}

func (e *EntropyOrangutan) heavyReadLoad() {
	fmt.Fprintln(GinkgoWriter, "EO: Performing Heavy Read Load")
	for i := 0; i < 10; i++ {
		e.Client.GetPrefix("")
	}
}

func (e *EntropyOrangutan) heavyWriteLoad() {
	fmt.Fprintln(GinkgoWriter, "EO: Performing Heavy Write Load")
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("%d", e.r.Int63())
		value := strings.Repeat("7", 100)
		e.Client.Set(key, value)
	}
}

func (e *EntropyOrangutan) deleteEverything() {
	fmt.Fprintln(GinkgoWriter, "EO: Deleting Everything")
	e.Client.DeletePrefix("")
}
