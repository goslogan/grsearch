package grsearch_test

import (
	grsearch "github.com/goslogan/grsearch"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Aliases", Label("search", "alias"), func() {
	It("can add an alias to the customers hash index", func() {
		cmd := client.FTAliasAdd(ctx, "alias1", "hcustomers")
		Expect(cmd.Err()).NotTo(HaveOccurred())
	})

	It("can delete an alias from the docs index", Label("FT.ALIASADD", "FT.ALIASDEL", "FT.SEARCH", "hash"), func() {
		cmd := client.FTAliasAdd(ctx, "aliasdel", "hcustomers")
		Expect(cmd.Err()).NotTo(HaveOccurred())
		cmd = client.FTAliasDel(ctx, "aliasdel")
		Expect(cmd.Err()).NotTo(HaveOccurred())
		Expect(client.FTSearchHash(ctx, "aliasdel", "haccount:232226", grsearch.NewQueryOptions()).Err()).To(HaveOccurred())
	})

	It("can move an alias to a different index", Label("FT.AlIASADD", "FT.ALIASUPDATE", "FT.SEARCH", "hash"), func() {
		cmd := client.FTAliasAdd(ctx, "customeralias", "jcustomers")
		Expect(cmd.Err()).NotTo((HaveOccurred()))
		search := client.FTSearchHash(ctx, "customeralias", `@id:{536299}`, grsearch.NewQueryOptions())
		Expect(search.Err()).NotTo(HaveOccurred())
		Expect(search.Len()).To(BeZero())
		cmd = client.FTAliasUpdate(ctx, "customeralias", "hcustomers")
		Expect(cmd.Err()).NotTo((HaveOccurred()))
		search = client.FTSearchHash(ctx, "customeralias", `@id:{536299}`, grsearch.NewQueryOptions())
		Expect(search.Err()).NotTo(HaveOccurred())
		Expect(search.Len()).ToNot(BeZero())

	})
})
