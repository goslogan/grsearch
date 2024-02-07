package grsearch_test

import (
	grsearch "github.com/goslogan/grsearch"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Aliases", func() {
	It("can add an alias to the docs index", func() {
		cmd := client.FTAliasAdd(ctx, "alias1", "hdocs")
		Expect(cmd.Err()).NotTo(HaveOccurred())
	})

	It("can delete an alias from the docs index", func() {
		cmd := client.FTAliasAdd(ctx, "aliasdel", "hdocs")
		Expect(cmd.Err()).NotTo(HaveOccurred())
		cmd = client.FTAliasDel(ctx, "aliasdel")
		Expect(cmd.Err()).NotTo(HaveOccurred())
		Expect(client.FTSearch(ctx, "aliasdel", "@command:set", grsearch.NewQueryOptions()).Err()).To(HaveOccurred())
	})

	It("can move an alias to a different index", func() {
		cmd := client.FTAliasAdd(ctx, "docalias", "hcustomers")
		Expect(cmd.Err()).NotTo((HaveOccurred()))
		search := client.FTSearch(ctx, "docalias", "@command:SET", grsearch.NewQueryOptions())
		Expect(search.Err()).NotTo(HaveOccurred())
		Expect(search.Len()).To(BeZero())
		cmd = client.FTAliasUpdate(ctx, "docalias", "hdocs")
		Expect(cmd.Err()).NotTo((HaveOccurred()))
		search = client.FTSearch(ctx, "docalias", "@command:SET", grsearch.NewQueryOptions())
		Expect(search.Err()).NotTo(HaveOccurred())
		Expect(search.Len()).ToNot(BeZero())

	})

	It("can find the same result via the original index and the alias", func() {
		cmd := client.FTAliasAdd(ctx, "alias2", "hdocs")
		Expect(cmd.Err()).NotTo(HaveOccurred())
		searchAlias := client.FTSearch(ctx, "alias2", "@command:set", grsearch.NewQueryOptions())
		Expect(searchAlias.Err()).NotTo(HaveOccurred())
		searchIndex := client.FTSearch(ctx, "hdocs", "@command:set", grsearch.NewQueryOptions())
		Expect(searchIndex.Err()).NotTo(HaveOccurred())
		Expect(searchAlias.Val()).To(BeEquivalentTo(searchIndex.Val()))

	})

})
