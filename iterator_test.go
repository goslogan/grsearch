package grsearch_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/goslogan/grsearch"
)

var _ = Describe("Iterator", Label("search", "hash", "ft.search", "iterator"), func() {

	It("can iterate over a search returning a single hash result", func() {
		cmd := client.FTSearch(ctx, "hdocs", "HGET", grsearch.NewQueryBuilder().SortBy("command").Verbatim().Options())
		Expect(cmd.Err()).NotTo(HaveOccurred())
		iterator := cmd.Iterator(ctx)
		Expect(iterator.Err()).NotTo(HaveOccurred())
		Expect(iterator.Next(ctx)).To(BeTrue())
		Expect(iterator.Val().Key).To(Equal("hcommand:HGET"))
		Expect(iterator.Next(ctx)).To(BeFalse())
	})

	It("can iterate over a search returning fewer than the limit results", func() {
		cmd := client.FTSearch(ctx, "hdocs", "GET", grsearch.NewQueryBuilder().SortBy("command").Verbatim().Options())
		Expect(cmd.Err()).NotTo(HaveOccurred())
		iterator := cmd.Iterator(ctx)
		Expect(iterator.Err()).NotTo(HaveOccurred())
		Expect(iterator.Next(ctx)).To(BeTrue())
		Expect(iterator.Val().Key).To(Equal("hcommand:CONFIG_GET"))
		Expect(iterator.Next(ctx)).To(BeTrue())
		Expect(iterator.Val().Key).To(Equal("hcommand:GET"))
		Expect(iterator.Next(ctx)).To(BeTrue())
		Expect(iterator.Val().Key).To(Equal("hcommand:SLOWLOG_GET"))
		Expect(iterator.Next(ctx)).To(BeFalse())

	})

	It("can iterate over a search return in multiple calls", func() {
		cmd := client.FTSearch(ctx, "hdocs", "GET",
			grsearch.NewQueryBuilder().
				SortBy("command").
				Limit(0, 2).
				Verbatim().
				Options())
		Expect(cmd.Err()).NotTo(HaveOccurred())
		iterator := cmd.Iterator(ctx)
		Expect(iterator.Err()).NotTo(HaveOccurred())
		Expect(iterator.Next(ctx)).To(BeTrue())
		Expect(iterator.Val().Key).To(Equal("hcommand:CONFIG_GET"))
		Expect(iterator.Next(ctx)).To(BeTrue())
		Expect(iterator.Val().Key).To(Equal("hcommand:GET"))
		Expect(iterator.Next(ctx)).To(BeTrue())
		Expect(iterator.Val().Key).To(Equal("hcommand:SLOWLOG_GET"))
		Expect(iterator.Next(ctx)).To(BeFalse())

	})

})
