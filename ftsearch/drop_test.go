package ftsearch_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/RedisLabs-Solution-Architects/go-search/ftsearch"
)

var _ = Describe("Drop", func() {

	BeforeEach(func() {
		// Requires redis on localhost:6379 with search module!
		Expect(client.FTCreateIndex(ctx, "drop_test", ftsearch.NewIndexOptions().AddSchemaAttribute(ftsearch.TextAttribute{
			Name:  "foo",
			Alias: "bar",
		})).Err()).NotTo(HaveOccurred())

		for _, row := range testData {
			Expect(client.HSet(ctx, fmt.Sprintf("account:%s", row[4]),
				"customer", row[0]+" "+row[1],
				"email", row[2],
				"ip", row[3],
				"account_id", row[4],
				"account_owner", row[5],
			).Err()).NotTo(HaveOccurred())
		}
	})

	It("can drop an index but keep the docs", func() {
		Expect(client.DBSize(ctx).Val()).To(Equal(int64(392)))
		Expect(client.FTDropIndex(ctx, "drop_test", false).Err()).NotTo(HaveOccurred())
		Expect(client.DBSize(ctx).Val()).To(Equal(int64(392)))
	})

	/* It("can drop an index and remove the docs", func() {
		Expect(client.DBSize(ctx).Val()).To(Equal(int64(392)))
		Expect(client.FTDropIndex(ctx, "drop_test", true).Err()).NotTo(HaveOccurred())
		Expect(client.DBSize(ctx).Val()).To(Equal(int64(0)))
	})*/

})
