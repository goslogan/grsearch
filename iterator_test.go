package grsearch_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/goslogan/grsearch"
)

var _ = Describe("Iterator", Label("search", "hash", "ft.search", "iterator"), func() {

	It("can iterate over a search returning a single hash result", func() {
		options := grsearch.NewQueryOptions()
		options.SortBy = "email"
		options.Verbatim = true
		cmd := client.FTSearchHash(ctx, "hcustomers", `@id:{232226}`, options)
		Expect(cmd.Err()).NotTo(HaveOccurred())
		iterator := cmd.Iterator(ctx)
		Expect(iterator.Err()).NotTo(HaveOccurred())
		Expect(iterator.Next(ctx)).To(BeTrue())
		Expect(iterator.Val().Key).To(Equal("haccount:232226"))
		Expect(iterator.Next(ctx)).To(BeFalse())
	})

	It("can fail effectively if the search returns no results", func() {
		options := grsearch.NewQueryOptions()
		options.SortBy = "email"
		options.Verbatim = true
		cmd := client.FTSearchHash(ctx, "hcustomers", `@id:{11111}`, options)
		Expect(cmd.Err()).NotTo(HaveOccurred())
		iterator := cmd.Iterator(ctx)
		Expect(iterator.Err()).NotTo(HaveOccurred())
		Expect(iterator.Next(ctx)).To(BeFalse())
	})

	It("can iterate over a search returning fewer than the limit results", func() {
		options := grsearch.NewQueryOptions()
		options.SortBy = "email"
		options.Verbatim = true
		cmd := client.FTSearchHash(ctx, "hcustomers", `@country:{UK}`, options)
		Expect(cmd.Err()).NotTo(HaveOccurred())
		iterator := cmd.Iterator(ctx)
		Expect(iterator.Err()).NotTo(HaveOccurred())
		Expect(iterator.Next(ctx)).To(BeTrue())
		Expect(iterator.Val().Key).To(Equal("haccount:1888382"))
		Expect(iterator.Next(ctx)).To(BeTrue())
		Expect(iterator.Val().Key).To(Equal("haccount:1952347"))
		Expect(iterator.Next(ctx)).To(BeTrue())
		Expect(iterator.Val().Key).To(Equal("haccount:1371128"))
		Expect(iterator.Next(ctx)).To(BeTrue())
		Expect(iterator.Val().Key).To(Equal("haccount:419113"))
		Expect(iterator.Next(ctx)).To(BeFalse())
	})

	It("can iterate over a search returned in multiple calls", func() {
		options := grsearch.NewQueryOptions()
		options.SortBy = "email"
		options.Limit = &grsearch.Limit{Offset: 0, Num: 2}
		cmd := client.FTSearchHash(ctx, "hcustomers", `@country:{UK}`, options)
		Expect(cmd.Err()).NotTo(HaveOccurred())
		iterator := cmd.Iterator(ctx)
		Expect(iterator.Err()).NotTo(HaveOccurred())
		Expect(iterator.Next(ctx)).To(BeTrue())
		Expect(iterator.Val().Key).To(Equal("haccount:1888382"))
		Expect(iterator.Next(ctx)).To(BeTrue())
		Expect(iterator.Val().Key).To(Equal("haccount:1952347"))
		Expect(iterator.Next(ctx)).To(BeTrue())
		Expect(iterator.Val().Key).To(Equal("haccount:1371128"))
		Expect(iterator.Next(ctx)).To(BeTrue())
		Expect(iterator.Val().Key).To(Equal("haccount:419113"))
		Expect(iterator.Next(ctx)).To(BeFalse())
	})

})
