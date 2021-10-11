package hash

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestMurmurHash2(t *testing.T) {
	NewWithT(t).Expect(MurmurHash2String([]byte("something"), 0)).To(Equal("crsxd7"))
}
