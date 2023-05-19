package grstack_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/goslogan/grstack"
)

var _ = Describe("We can build query options", Label("builders", "ft.search"), func() {

	It("can execute a very simple aggregate", func() {
		opts := grstack.NewAggregateOptionsBuilder().
			GroupBy(grstack.NewGroupByBuilder().
				Property("@owner").
				Reduce(grstack.ReduceSum("@balance", "total_balance")).
				GroupBy())

		cmd := client.FTAggregate(ctx, "hcustomers", "*", opts.Options())
		Expect(cmd.Err()).ToNot(HaveOccurred())
	})

})
