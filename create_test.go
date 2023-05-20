package grstack_test

import (
	"github.com/goslogan/grstack"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Create", Label("json", "query", "ft.create"), func() {

	It("can build the simplest JSON index", func() {
		createCmd := client.FTCreate(ctx, "jsimple", grstack.NewIndexBuilder().Schema(grstack.TextAttribute{
			Name:  "$.foo",
			Alias: "bar",
		}).On("json").Options())
		Expect(createCmd.Err()).NotTo(HaveOccurred())
		Expect(createCmd.String()).To(Equal("ft.create jsimple on json score 1 schema $.foo as bar text: true"))
	})

	It("can build a json index with options", func() {
		createCmd := client.FTCreate(ctx, "jwithoptions", grstack.NewIndexBuilder().
			Prefix("jaccount:").
			On("json").
			MaxTextFields().
			Score(0.5).
			Language("spanish").
			Schema(grstack.TextAttribute{
				Name:  "$.foo",
				Alias: "bar",
			}).Options())
		Expect(createCmd.Err()).NotTo(HaveOccurred())
		Expect(createCmd.String()).To(Equal("ft.create jwithoptions on json prefix 1 jaccount: language spanish score 0.5 maxtextfields schema $.foo as bar text: true"))
	})

})
