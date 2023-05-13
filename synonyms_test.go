package grstack_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Synonyms", Ordered, func() {
	It("can add one or more synonyms", func() {
		cmd := client.FTSynUpdate(ctx, "hdocs", "synomyns", "hash", "map")
		Expect(cmd.Err()).NotTo(HaveOccurred())
		Expect(cmd.Val()).To(BeTrue())
		cmd = client.FTSynUpdate(ctx, "hdocs", "words", "list", "array")
		Expect(cmd.Err()).NotTo(HaveOccurred())
	})

	It("cam dump synonyms", func() {
		cmd := client.FTSynDump(ctx, "hdocs")
		Expect(cmd.Err()).NotTo(HaveOccurred())
		Expect(cmd.Val()).To(BeEquivalentTo(map[string][]string{
			"map":   {"hash"},
			"array": {"list"},
			"":      {"hash", "list"},
		}))
	})
})
