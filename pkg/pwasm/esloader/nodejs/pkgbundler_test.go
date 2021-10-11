package nodejs

import "testing"

func TestPkgBundler(t *testing.T) {
	pb := PkgBundler{}
	pb.SetDefaults()

	t.Run("Bundle", func(t *testing.T) {
		pb.Bundle("workbox-core")
	})
}
