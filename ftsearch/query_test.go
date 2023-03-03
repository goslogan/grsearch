package ftsearch_test

import (
	"math"

	"github.com/RedisLabs-Solution-Architects/go-search/ftsearch"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// ft.search accounts "@tam:nic\.gibson"
var _ = Describe("Query basics", func() {

	It("doesn't raise an error on a valid query", func() {
		Expect(client.FTSearch(ctx, "customers", `@email:ejowers0\@unblog\.fr`, nil).
			Err()).To(Not(HaveOccurred()))
	})
	It("can generate a valid query", func() {
		Expect(client.FTSearch(ctx, "customers", `@email:ejowers0\@unblog\.fr`, nil).
			String()).To(ContainSubstring(`ft.search customers @email:ejowers0\@unblog\.fr`))
	})

	It("can search for a specific result by attribute", func() {
		Expect(client.FTSearch(ctx, "customers", `@email:ejowers0\@unblog\.fr`, nil).Err()).NotTo(HaveOccurred())
	})

	It("can return a single result", func() {
		Expect(client.FTSearch(ctx, "customers", `@email:ejowers0\@unblog\.fr`, nil).Len()).To(Equal(1))
	})

	It("can return a map result", func() {
		Expect(client.FTSearch(ctx, "customers", `@id:{1121175}`, nil).Val()).To(Equal(
			map[string]*ftsearch.QueryResult{
				"account:1121175": {
					Score: 0,
					Value: map[string]string{
						"customer":      "Kandace Korneichuk",
						"email":         `kkorneichukc\@cpanel\.net`,
						"ip":            `148\.140\.255\.235`,
						"account_id":    "1121175",
						"account_owner": "nic.gibson",
						"balance":       "927.00",
					}}}))
	})

	It("will fail quietly with no search defined", func() {
		cmd := client.FTSearch(ctx, "customers", "", nil)
		Expect(cmd.Err()).NotTo(HaveOccurred())
		Expect(cmd.Len()).To(Equal(0))
	})

	It("can return all the results for a given tag", func() {
		Expect(client.FTSearch(ctx, "customers", `@owner:{nic\.gibson}`, nil).Len()).To(Equal(10))
	})

})

var _ = Describe("Query options", func() {

	It("will return empty results - NOCONTENT", func() {
		Expect(client.FTSearch(ctx, "customers", `@id:{1121175}`, ftsearch.NewQueryOptions().WithoutContent()).Val()).To(Equal(
			map[string]*ftsearch.QueryResult{
				"account:1121175": {
					Score: 0,
					Value: nil}}))
	})

	It("will return scores - WITHSCORES", func() {
		cmd := client.FTSearch(ctx, "customers", `@id:{1121175}`, ftsearch.NewQueryOptions().WithScores())
		Expect(cmd.Err()).NotTo(HaveOccurred())
		Expect(cmd.Val()["account:1121175"].Score).Should(BeNumerically(">=", 1))
	})

	It("will return filtered results - FILTER (numeric)", func() {
		cmd := client.FTSearch(ctx, "customers", `@owner:{nic\.gibson}`, ftsearch.NewQueryOptions().
			WithoutContent().
			AddFilter(ftsearch.NewQueryFilter("balance").WithMinExclusive(0).WithMaxInclusive(math.Inf(1))))
		Expect(cmd.Err()).NotTo(HaveOccurred())
		Expect(len(cmd.Val())).To(Equal(2))
	})

})
