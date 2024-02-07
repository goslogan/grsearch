package grsearch_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/goslogan/grsearch"
)

var _ = Describe("We can build query options", Label("aggregate", "search", "ft.aggregate"), func() {

	It("can execute a very simple aggregate", func() {
		opts := grsearch.NewAggregateBuilder().
			GroupBy(grsearch.NewGroupByBuilder().
				Property("@owner").
				Reduce(grsearch.ReduceSum("@balance", "total_balance")).
				GroupBy())

		cmd := client.FTAggregate(ctx, "hcustomers", "*", opts.Options())
		Expect(cmd.Err()).ToNot(HaveOccurred())
	})

})

var _ = Describe("Synonyms", Ordered, Label("synonyms", "search"), func() {
	It("can add one or more synonyms", Label("ft.synupdate"), func() {
		cmd := client.FTSynUpdate(ctx, "hdocs", "synomyns", "hash", "map")
		Expect(cmd.Err()).NotTo(HaveOccurred())
		Expect(cmd.Val()).To(BeTrue())
		cmd = client.FTSynUpdate(ctx, "hdocs", "words", "list", "array")
		Expect(cmd.Err()).NotTo(HaveOccurred())
	})

	It("cam dump synonyms", Label("ft.syndump"), func() {
		cmd := client.FTSynDump(ctx, "hdocs")
		Expect(cmd.Err()).NotTo(HaveOccurred())
		Expect(cmd.Val()).To(BeEquivalentTo(map[string][]string{
			"map":   {"hash"},
			"array": {"list"},
			"":      {"hash", "list"},
		}))
	})
})

var _ = Describe("Tagvals", Label("search", "ft.tagvals"), func() {

	It("can get tag values from the index", func() {
		cmd := client.FTTagVals(ctx, "hcustomers", "owner")
		Expect(cmd.Err()).NotTo(HaveOccurred())

		vals := cmd.Val()
		Expect(vals).To(ContainElements([]string{"lara.croft", "ellen.ripley", "sarah.oconnor"}))
		Expect(len(vals)).To(Equal(3))
	})
})

var _ = Describe("Info", Label("search", "ft.info"), func() {

	It("can get info from a simple index", func() {
		cmd := client.FTInfo(ctx, "hcustomers")
		Expect(cmd.Err()).NotTo(HaveOccurred())
	})

	It("can recreate an index", Label("rebuild"), func() {
		cmd1 := client.FTInfo(ctx, "hcustomers")
		Expect(cmd1.Err()).NotTo(HaveOccurred())
		cmd2 := client.FTCreate(ctx, "hcustomersdup", cmd1.Val().Index)
		Expect(cmd2.Err()).NotTo(HaveOccurred())
		time.Sleep(5 * time.Second)
		cmd3 := client.FTInfo(ctx, "hcustomersdup")
		Expect(cmd3.Err()).NotTo(HaveOccurred())
		Expect(cmd1.Val().Index).To(Equal(cmd3.Val().Index))
		Expect(cmd1.Val().NumDocs).To(Equal(cmd3.Val().NumDocs))
		Expect(cmd1.Val().MaxDocId).To(Equal(cmd3.Val().MaxDocId))
		Expect(cmd1.Val().SortableValuesSize).To(Equal(cmd3.Val().SortableValuesSize))
	})

})
