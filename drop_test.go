package grstack_test

import (
	"github.com/goslogan/grstack"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Drop", func() {

	// TODO - fix this - we don't want to mess up the tests

	It("can drop an index but keep the docs", func() {
		createCmd := client.FTCreateIndex(ctx, "drop_test", grstack.NewIndexBuilder().Schema(grstack.TextAttribute{
			Name:  "foo",
			Alias: "bar",
		}).Options())
		Expect(createCmd.Err()).NotTo(HaveOccurred())

		cmd := client.DBSize(ctx)
		Expect(cmd.Err()).NotTo(HaveOccurred())
		size := cmd.Val()
		Expect(client.FTDropIndex(ctx, "drop_test", false).Err()).NotTo(HaveOccurred())
		Expect(client.DBSize(ctx).Val()).To(Equal(size))
	})

	// TODO - Fix to add a new index with one or two docs for drop test
	/* It("can drop an index and remove the docs", func() {
		Expect(client.DBSize(ctx).Val()).To(Equal(int64(392)))
		Expect(client.FTDropIndex(ctx, "drop_test", true).Err()).NotTo(HaveOccurred())
		Expect(client.DBSize(ctx).Val()).To(Equal(int64(0)))
	})*/

})
