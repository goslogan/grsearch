package grsearch_test

import (
	grsearch "github.com/goslogan/grsearch"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Aggregate", Label("ft.aggregate"), func() {

	It("can run the most basic aggregates", func() {
		opts := grsearch.NewAggregateOptions()
		opts.Load = []grsearch.AggregateLoad{{Name: "owner"}, {Name: "customer"}}

		cmd := client.FTAggregate(ctx, "hcustomers", `*`, opts)
		Expect(cmd.Err()).NotTo(HaveOccurred())
		Expect(cmd.TotalResults()).To(Equal(int64(25)))

		keys := []string{}
		for key := range cmd.Val()[0] {
			keys = append(keys, key)
		}
		Expect(keys).To(ConsistOf("owner", "customer"))

	})

	It("can run a single group by", func() {
		opts := grsearch.NewAggregateBuilder().
			GroupBy(grsearch.NewGroupByBuilder().
				Property("@owner").
				Reduce(grsearch.ReduceSum("@balance", "total_balance")).
				GroupBy())

		cmd := client.FTAggregate(ctx, "hcustomers", "*", opts.Options())
		Expect(cmd.Err()).ToNot(HaveOccurred())
		Expect(cmd.TotalResults()).To(Equal(int64(3)))

		keys := []string{}
		for key := range cmd.Val()[0] {
			keys = append(keys, key)
		}
		Expect(keys).To(ConsistOf("owner", "total_balance"))
	})

	It("can sort ", func() {
		opts := grsearch.NewAggregateOptions()
		opts.Load = []grsearch.AggregateLoad{{Name: "owner"}, {Name: "customer"}}
		opts.Steps = []grsearch.AggregateStep{&grsearch.AggregateSort{Keys: []grsearch.AggregateSortKey{{Name: "owner"}}}}
		cmd := client.FTAggregate(ctx, "hcustomers", `*`, opts)
		Expect(cmd.Err()).NotTo(HaveOccurred())
		// Expect(cmd.TotalResults()).To(Equal(int64(25)))
	})

	It("can execute an APPLY statement", func() {
		opts := grsearch.NewAggregateBuilder().
			GroupBy(grsearch.NewGroupByBuilder().
				Property("@country").Property("@owner").
				Reduce(grsearch.ReduceSum("@balance", "total_balance")).
				GroupBy()).
			Apply("upper(@owner)", "uowner")
		cmd := client.FTAggregate(ctx, "hcustomers", "*", opts.Options())
		Expect(cmd.Err()).ToNot(HaveOccurred())
		Expect(cmd.TotalResults()).To(Equal(int64(13)))

	})

})
