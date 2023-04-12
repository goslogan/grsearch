package grstack_test

import (
	"github.com/RedisLabs-Solution-Architects/grstack"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Aliases", func() {
	It("can add an alias to the docs index", func() {
		cmd := client.FTAliasAdd(ctx, "alias1", "docs")
		Expect(cmd.Err()).NotTo(HaveOccurred())
	})

	It("can delete an alias from the docs index", func() {
		cmd := client.FTAliasAdd(ctx, "aliasdel", "docs")
		Expect(cmd.Err()).NotTo(HaveOccurred())
		cmd = client.FTAliasDel(ctx, "aliasdel")
		Expect(cmd.Err()).NotTo(HaveOccurred())
		Expect(client.FTSearch(ctx, "aliasdel", "@command:set", grstack.NewQueryOptions()).Err()).To(HaveOccurred())
	})

	It("can move an alias to a different index", func() {
		cmd := client.FTAliasAdd(ctx, "docalias", "customers")
		Expect(cmd.Err()).NotTo((HaveOccurred()))
		search := client.FTSearch(ctx, "docalias", "@command:set", grstack.NewQueryOptions())
		Expect(search.Err()).NotTo(HaveOccurred())
		Expect(len(search.Val())).To(BeZero())
		cmd = client.FTAliasUpdate(ctx, "docalias", "docs")
		Expect(cmd.Err()).NotTo((HaveOccurred()))
		search = client.FTSearch(ctx, "docalias", "@command:set", grstack.NewQueryOptions())
		Expect(search.Err()).NotTo(HaveOccurred())
		Expect(len(search.Val())).ToNot(BeZero())

	})

	It("can find the same result via the original index and the alias", func() {
		cmd := client.FTAliasAdd(ctx, "alias2", "docs")
		Expect(cmd.Err()).NotTo(HaveOccurred())
		searchAlias := client.FTSearch(ctx, "alias2", "@command:set", grstack.NewQueryOptions())
		Expect(searchAlias.Err()).NotTo(HaveOccurred())
		searchIndex := client.FTSearch(ctx, "docs", "@command:set", grstack.NewQueryOptions())
		Expect(searchIndex.Err()).NotTo(HaveOccurred())
		Expect(searchAlias.Val()).To(BeEquivalentTo(searchIndex.Val()))

	})

})
