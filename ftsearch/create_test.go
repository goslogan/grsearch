package ftsearch_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/RedisLabs-Solution-Architects/go-search/ftsearch"
)

var _ = Describe("Create", func() {

	It("can build the simplest index", func() {
		createCmd := client.FTCreateIndex(ctx, "simple", ftsearch.NewIndexOptions().AddSchemaAttribute(ftsearch.TextAttribute{
			Name:  "foo",
			Alias: "bar",
		}))
		Expect(createCmd.Err()).NotTo(HaveOccurred())
		Expect(createCmd.String()).To(Equal("ft.create simple on hash score 1 schema foo as bar text: true"))
	})

	It("can build a hash index with options", func() {
		createCmd := client.FTCreateIndex(ctx, "withoptions", ftsearch.NewIndexOptions().
			AddPrefix("account:").
			WithMaxTextFields().
			WithScore(0.5).
			WithLanguage("spanish").
			AddSchemaAttribute(ftsearch.TextAttribute{
				Name:  "foo",
				Alias: "bar",
			}))
		Expect(createCmd.Err()).NotTo(HaveOccurred())
		Expect(createCmd.String()).To(Equal("ft.create withoptions on hash prefix 1 account: language spanish score 0.5 maxtextfields schema foo as bar text: true"))
	})

	It("can build a hash index with multiple schema entries", func() {
		createCmd := client.FTCreateIndex(ctx, "multiattrib", ftsearch.NewIndexOptions().
			AddSchemaAttribute(ftsearch.TextAttribute{
				Name:  "texttest",
				Alias: "xxtext",
			}).
			AddSchemaAttribute(ftsearch.NumericAttribute{
				Name:     "numtest",
				Sortable: true,
			}))
		Expect(createCmd.Err()).NotTo(HaveOccurred())
		Expect(createCmd.String()).To(Equal("ft.create multiattrib on hash score 1 schema texttest as xxtext text numtest numeric sortable: true"))
	})

})
