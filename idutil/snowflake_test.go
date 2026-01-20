package idutil_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/Jsharkc/mygopkg/idutil"
)

func TestNextID(t *testing.T) {
	sf := idutil.New()
	for i := 0; i < 100; i++ {
		go func(i int) {
			fmt.Println(i, sf.NextID(), sf.String())
		}(i)
	}

	time.Sleep(2e9)
}
