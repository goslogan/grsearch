package grstack_test

import (
	"math"

	grstack "github.com/goslogan/grstack"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// ft.search accounts "@tam:nic\.gibson"
var _ = Describe("Query basics", Label("query", "ft.search"), func() {

	It("doesn't raise an error on a valid query", func() {
		Expect(client.FTSearch(ctx, "hcustomers", `@email:ejowers0\@unblog\.fr`, nil).
			Err()).To(Not(HaveOccurred()))
	})
	It("can generate a valid query", func() {
		Expect(client.FTSearch(ctx, "hcustomers", `@email:ejowers0@unblog.fr`, nil).
			String()).To(ContainSubstring(`ft.search hcustomers @email:ejowers0@unblog.fr`))
	})

	It("can search for a specific result by attribute", func() {
		cmd := client.FTSearch(ctx, "hcustomers", `@email:ejowers0\@unblog\.fr`, nil)
		Expect(cmd.Err()).NotTo(HaveOccurred())
	})

	It("can return a single result", func() {
		Expect(client.FTSearch(ctx, "hcustomers", `@email:ejowers0\@unblog\.fr`, nil).Len()).To(Equal(1))
	})

	It("can return a map result", func() {
		Expect(client.FTSearch(ctx, "hcustomers", `@id:{1121175}`, nil).Val()).To(Equal(
			map[string]*grstack.QueryResult{
				"haccount:1121175": {
					Score: 0,
					Value: map[string]string{
						"account_owner": "lara.croft",
						"balance":       "927.00",
						"customer":      "Kandace Korneichuk",
						"email":         `kkorneichukc\@cpanel\.net`,
						"ip":            `148\.140\.255\.235`,
						"account_id":    "1121175",
					}}}))
	})

	It("will fail quietly with no search defined", func() {
		cmd := client.FTSearch(ctx, "hcustomers", "", nil)
		Expect(cmd.Err()).NotTo(HaveOccurred())
		Expect(cmd.Len()).To(Equal(0))
	})

	It("can return all the results for a given tag", func() {
		Expect(client.FTSearch(ctx, "hcustomers", `@owner:{lara\.croft}`, nil).Len()).To(Equal(10))
	})

})

var _ = Describe("Query options", Label("query", "ft.search"), func() {

	It("will return empty results - NOCONTENT", func() {
		opts := grstack.NewQueryOptions()
		opts.NoContent = true
		cmd := client.FTSearch(ctx, "hcustomers", `@id:{1121175}`, opts)
		Expect(cmd.Val()).To(Equal(
			map[string]*grstack.QueryResult{
				"haccount:1121175": {
					Score: 0,
					Value: nil}}))
	})

	It("will return scores - WITHSCORES", func() {
		opts := grstack.NewQueryOptions()
		opts.WithScores = true
		cmd := client.FTSearch(ctx, "hcustomers", `@id:{1121175}`, opts)
		Expect(cmd.Err()).NotTo(HaveOccurred())
		Expect(cmd.Val()["haccount:1121175"].Score).Should(BeNumerically(">=", 1))
	})

	It("will return filtered results - FILTER (numeric)", func() {
		opts := grstack.NewQueryOptions()
		opts.NoContent = true
		opts.Filters = []grstack.QueryFilter{
			{
				Attribute: "balance",
				Min:       grstack.FilterValue(0, true),
				Max:       grstack.FilterValue(math.Inf(1), false),
			},
		}
		cmd := client.FTSearch(ctx, "hcustomers", `@owner:{lara\.croft}`, opts)
		Expect(cmd.Err()).NotTo(HaveOccurred())
		Expect(len(cmd.Val())).To(Equal(2))
	})

	It("can explain a score", func() {
		opts := grstack.NewQueryOptions()
		opts.WithScores = true
		opts.ExplainScore = true
		cmd := client.FTSearch(ctx, "hcustomers", `@owner:{lara\.croft}`, opts)
		Expect(cmd.Err()).NotTo(HaveOccurred())
		Expect(len(cmd.Val())).NotTo(BeZero())
		Expect(cmd.Val()["haccount:806396"].Explanation).NotTo(BeNil())
	})

	It("can sort results", func() {
		results := []string{}
		opts := grstack.NewQueryOptions()
		opts.NoContent = true
		opts.SortBy = "customer"
		cmd := client.FTSearch(ctx, "hcustomers", `@owner:{lara\.croft}`, opts)
		Expect(cmd.Err()).NotTo(HaveOccurred())
		for k := range cmd.Val() {
			results = append(results, k)
		}
		Expect(results).To(ConsistOf([]string{"haccount:1339089", "haccount:239155", "haccount:575072", "haccount:765279", "haccount:1826581", "haccount:1371128", "haccount:1121175", "haccount:886088", "haccount:806396", "haccount:507187"}))

	})

})
