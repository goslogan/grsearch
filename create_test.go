package grstack_test

import (
	"github.com/goslogan/grstack"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Create", Label("json", "query", "ft.create"), func() {

	It("can build the simplest JSON index", func() {
		options := grstack.NewIndexOptions()
		options.On = "JSON"
		options.Schema = []grstack.SchemaAttribute{&grstack.TextAttribute{
			Name:  "$.foo",
			Alias: "bar",
		}}
		createCmd := client.FTCreate(ctx, "jsimple", options)
		Expect(createCmd.Err()).NotTo(HaveOccurred())
		Expect(createCmd.String()).To(Equal("FT.CREATE jsimple ON JSON SCORE 1 SCHEMA $.foo AS bar TEXT: true"))
	})

	It("can build a json index with options", func() {
		options := grstack.NewIndexOptions()
		options.On = "JSON"
		options.Prefix = []string{"jaccount:"}
		options.Schema = []grstack.SchemaAttribute{&grstack.TextAttribute{
			Name:  "$.foo",
			Alias: "bar",
		}}
		options.MaxTextFields = true
		options.Score = 0.5
		options.Language = "spanish"
		createCmd := client.FTCreate(ctx, "jwithoptions", options)
		Expect(createCmd.Err()).NotTo(HaveOccurred())
		Expect(createCmd.String()).To(Equal("FT.CREATE jwithoptions ON JSON PREFIX 1 jaccount: LANGUAGE spanish SCORE 0.5 MAXTEXTFIELDS SCHEMA $.foo AS bar TEXT: true"))
	})

})
