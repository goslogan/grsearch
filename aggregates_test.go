package grstack_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/goslogan/grstack"
)

var _ = Describe("We can build query options", Label("builders", "ft.search"), func() {

	It("can execute a very simple aggregate", func() {
		opts := grstack.NewAggregateBuilder().
			GroupBy(grstack.NewGroupByBuilder().
				Property("@owner").
				Reduce(grstack.ReduceSum("@balance", "total_balance")).
				GroupBy())

		cmd := client.FTAggregate(ctx, "customers", "*", opts.Options())
		Expect(cmd.Err()).ToNot(HaveOccurred())
	})

})
