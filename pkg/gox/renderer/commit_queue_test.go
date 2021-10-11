package renderer

import (
	"testing"

	"github.com/onsi/gomega"
)

func Test_actionQueueBufferEach(t *testing.T) {
	q := []func(){
		func() {},
		func() {},
		func() {},
		func() {},
		func() {},
	}

	bufLens := make([]int, 0)

	actionQueueBufferEach(q, 2, func(queue []func()) {
		bufLens = append(bufLens, len(queue))
	})

	gomega.NewWithT(t).Expect(bufLens).To(gomega.Equal([]int{
		2, 2, 1,
	}))
}
