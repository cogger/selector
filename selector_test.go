package selector

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"reflect"

	"golang.org/x/net/context"
	"gopkg.in/cogger/cogger.v1"
	"gopkg.in/cogger/cogger.v1/cogs"
)

var _ = Describe("New", func() {
	coggerInterface := reflect.TypeOf((*cogger.Cog)(nil)).Elem()
	selectorInterface := reflect.TypeOf((*Selector)(nil)).Elem()

	It("should create a cog", func() {
		cog := New()

		Expect(reflect.TypeOf(cog).Implements(coggerInterface)).To(BeTrue())
		Expect(<-cog.Do(context.Background())).ToNot(HaveOccurred())
	})

	It("should return a selector when you add a case", func() {
		selectorInterface := reflect.TypeOf((*Selector)(nil)).Elem()
		selector := New()
		result := selector.Case(func(ctx context.Context) bool { return true }, cogs.NoOp())
		Expect(reflect.TypeOf(result).Implements(selectorInterface)).To(BeTrue())
	})

	It("should return a selector when you add a default case", func() {
		selector := New()
		result := selector.Default(cogs.NoOp())
		Expect(reflect.TypeOf(result).Implements(selectorInterface)).To(BeTrue())
	})

	Context("when a executing a job", func() {
		It("should return nil where there are no errors", func() {
			Expect(<-New().Do(context.Background())).ToNot(HaveOccurred())
		})

		It("should return an error when there is an error", func() {
			testErr := errors.New("test error")
			s := New()
			s = s.Default(cogs.ReturnErr(testErr))
			Expect(<-s.Do(context.Background())).To(Equal(testErr))
		})

		It("should execute the first case that passes", func() {
			for i := 1; i < 10; i++ {
				s := New()
				count := 0

				failed := func(ctx context.Context) bool {
					count++
					return false
				}

				passed := func(ctx context.Context) bool {
					count++
					return true
				}

				for j := 0; j < 10; j++ {
					if j < i {
						s = s.Case(failed, cogs.NoOp())
					} else {
						s = s.Case(passed, cogs.NoOp())
					}
				}
				ctx := context.Background()
				defaultFired := false

				s = s.Default(cogs.Simple(ctx, func() error {
					defaultFired = true
					return nil
				}))

				err := <-s.Do(ctx)
				Expect(err).ToNot(HaveOccurred())
				Expect(count).To(Equal(i + 1))
				Expect(defaultFired).To(BeFalse())
			}
		})

		It("should execute the default case if none pass", func() {
			s := New()

			failed := func(ctx context.Context) bool { return false }

			for i := 0; i < 10; i++ {
				s = s.Case(failed, cogs.NoOp())
			}

			ctx := context.Background()
			defaultFired := false

			s = s.Default(cogs.Simple(ctx, func() error {
				defaultFired = true
				return nil
			}))

			err := <-s.Do(ctx)
			Expect(err).ToNot(HaveOccurred())
			Expect(defaultFired).To(BeTrue())
		})

	})

	It("should implement SetLimit function", func() {
		s := New()

		limit := &mockLimit{}

		s.SetLimit(limit)

		ctx := context.Background()
		Expect(<-s.Do(ctx)).ToNot(HaveOccurred())
		Expect(limit.NextHits).To(Equal(1))
		Expect(limit.Completed).To(BeTrue())
	})
})

type mockLimit struct {
	Completed bool
	NextHits  int
}

func (limit *mockLimit) Next(ctx context.Context) chan struct{} {
	next := make(chan struct{})
	go func() {
		limit.NextHits++
		next <- struct{}{}
	}()
	return next
}

func (limit *mockLimit) Done(ctx context.Context) {
	limit.Completed = true
}
