package renderer_test

import (
	"context"
	"testing"

	. "github.com/go-courier/gox/pkg/dom"
	. "github.com/go-courier/gox/pkg/gox"
	"github.com/go-courier/gox/pkg/gox/renderer"
	"github.com/onsi/gomega"
)

type AppWithEffectHook struct {
	Setup   *func() func()
	Refresh string
}

func (a AppWithEffectHook) Render(ctx context.Context, children ...interface{}) interface{} {
	UseEffect(ctx, *a.Setup, []interface{}{a.Refresh})
	return nil
}

func TestRenderWithEffectHook(t *testing.T) {
	ctx := context.Background()
	root := Document.CreateElement("body")
	r := renderer.CreateRoot(root)

	counts := struct {
		setup   int
		cleanup int
	}{}

	setup := func() func() {
		counts.setup++
		return func() {
			counts.cleanup++
		}
	}

	_ = r.Render(ctx, H(AppWithEffectHook{
		Setup: &setup,
	})())

	gomega.NewWithT(t).Expect(counts.setup).To(gomega.Equal(1))
	gomega.NewWithT(t).Expect(counts.cleanup).To(gomega.Equal(0))

	t.Run("should not re setup when deps not changed", func(t *testing.T) {
		_ = r.Render(ctx, H(AppWithEffectHook{
			Setup: &setup,
		})())

		gomega.NewWithT(t).Expect(counts.setup).To(gomega.Equal(1))
		gomega.NewWithT(t).Expect(counts.cleanup).To(gomega.Equal(0))

		t.Run("should re setup when deps changed", func(t *testing.T) {
			_ = r.Render(ctx, H(AppWithEffectHook{
				Setup:   &setup,
				Refresh: "1",
			})())

			gomega.NewWithT(t).Expect(counts.setup).To(gomega.Equal(2))
			gomega.NewWithT(t).Expect(counts.cleanup).To(gomega.Equal(1))
		})

	})

}
