package ftsearch_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Synonyms", Ordered, func() {
	It("can add one or more synonyms", func() {
		cmd := client.FTSynUpdate(ctx, "docs", "synomyns", "hash", "map")
		Expect(cmd.Err()).NotTo(HaveOccurred())
		Expect(cmd.Val()).To(BeTrue())
		cmd = client.FTSynUpdate(ctx, "docs", "words", "list", "array")
		Expect(cmd.Err()).NotTo(HaveOccurred())
	})

	It("cam dump synonyms", func() {
		cmd := client.FTSynDump(ctx, "docs")
		Expect(cmd.Err()).NotTo(HaveOccurred())
		Expect(cmd.Val()).To(BeEquivalentTo(map[string][]string{
			"map":   {"hash"},
			"array": {"list"},
			"":      {"hash", "list"},
		}))
	})
})
