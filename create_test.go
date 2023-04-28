package grstack_test

import (
	"github.com/goslogan/grstack"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Create", func() {

	It("can build the simplest index", func() {
		createCmd := client.FTCreateIndex(ctx, "simple", grstack.NewIndexOptions().AddSchemaAttribute(grstack.TextAttribute{
			Name:  "foo",
			Alias: "bar",
		}))
		Expect(createCmd.Err()).NotTo(HaveOccurred())
		Expect(createCmd.String()).To(Equal("ft.create simple on hash score 1 schema foo as bar text: true"))
	})

	It("can build a hash index with options", func() {
		createCmd := client.FTCreateIndex(ctx, "withoptions", grstack.NewIndexOptions().
			AddPrefix("account:").
			WithMaxTextFields().
			WithScore(0.5).
			WithLanguage("spanish").
			AddSchemaAttribute(grstack.TextAttribute{
				Name:  "foo",
				Alias: "bar",
			}))
		Expect(createCmd.Err()).NotTo(HaveOccurred())
		Expect(createCmd.String()).To(Equal("ft.create withoptions on hash prefix 1 account: language spanish score 0.5 maxtextfields schema foo as bar text: true"))
	})

	It("can build a hash index with multiple schema entries", func() {
		createCmd := client.FTCreateIndex(ctx, "multiattrib", grstack.NewIndexOptions().
			AddSchemaAttribute(grstack.TextAttribute{
				Name:  "texttest",
				Alias: "xxtext",
			}).
			AddSchemaAttribute(grstack.NumericAttribute{
				Name:     "numtest",
				Sortable: true,
			}))
		Expect(createCmd.Err()).NotTo(HaveOccurred())
		Expect(createCmd.String()).To(Equal("ft.create multiattrib on hash score 1 schema texttest as xxtext text numtest numeric sortable: true"))
	})

	It("can build a hash index with multiple schema entries and a different language", func() {
		createCmd := client.FTCreateIndex(ctx, "language", grstack.NewIndexOptions().
			AddSchemaAttribute(grstack.TextAttribute{
				Name:  "texttest",
				Alias: "xxtext",
			}).
			AddSchemaAttribute(grstack.NumericAttribute{
				Name:     "numtest",
				Sortable: true,
			}).WithLanguage("german"))
		Expect(createCmd.Err()).NotTo(HaveOccurred())
		Expect(createCmd.String()).To(Equal("ft.create language on hash language german score 1 schema texttest as xxtext text numtest numeric sortable: true"))
	})

	It("can build an index with a language field and score field", func() {
		createCmd := client.FTCreateIndex(ctx, "langscore", grstack.NewIndexOptions().AddSchemaAttribute(grstack.TextAttribute{
			Name:  "foo",
			Alias: "bar",
		}).WithLanguageField("lng").WithScoreField("scr"))
		Expect(createCmd.Err()).NotTo(HaveOccurred())
		Expect(createCmd.String()).To(Equal("ft.create langscore on hash language_field lng score 1 score_field scr schema foo as bar text: true"))
	})

	It("can build an index with NOFIELDS", func() {
		createCmd := client.FTCreateIndex(ctx, "nofields", grstack.NewIndexOptions().AddSchemaAttribute(grstack.TextAttribute{
			Name:  "foo",
			Alias: "bar",
		}).WithNoFields())
		Expect(createCmd.Err()).NotTo(HaveOccurred())
		Expect(createCmd.String()).To(Equal("ft.create nofields on hash score 1 nofields schema foo as bar text: true"))
	})

	It("can build an index with NOHL", func() {
		createCmd := client.FTCreateIndex(ctx, "nohl", grstack.NewIndexOptions().AddSchemaAttribute(grstack.TextAttribute{
			Name:  "foo",
			Alias: "bar",
		}).WithNoHighlight())
		Expect(createCmd.Err()).NotTo(HaveOccurred())
		Expect(createCmd.String()).To(Equal("ft.create nohl on hash score 1 nohl schema foo as bar text: true"))
	})

	It("can build an index with NOOFFSETS", func() {
		createCmd := client.FTCreateIndex(ctx, "nooff", grstack.NewIndexOptions().AddSchemaAttribute(grstack.TextAttribute{
			Name:  "foo",
			Alias: "bar",
		}).WithNoOffsets())
		Expect(createCmd.Err()).NotTo(HaveOccurred())
		Expect(createCmd.String()).To(Equal("ft.create nooff on hash score 1 nooffsets schema foo as bar text: true"))
	})

	It("can build an index with NOFREQS", func() {
		createCmd := client.FTCreateIndex(ctx, "nofr", grstack.NewIndexOptions().AddSchemaAttribute(grstack.TextAttribute{
			Name:  "foo",
			Alias: "bar",
		}).WithNoFreqs())
		Expect(createCmd.Err()).NotTo(HaveOccurred())
		Expect(createCmd.String()).To(Equal("ft.create nofr on hash score 1 nofreqs schema foo as bar text: true"))
	})

})
