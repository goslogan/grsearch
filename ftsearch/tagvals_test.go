package ftsearch_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tagvals", func() {

	It("can get tag values from the index", func() {
		cmd := client.FTTagVals(ctx, "customers", "owner")
		Expect(cmd.Err()).NotTo(HaveOccurred())

		vals := cmd.Val()
		Expect(vals).To(ContainElements([]string{"nic.gibson", "daniel.preiskel", "igor.jovanovic"}))
		Expect(len(vals)).To(Equal(3))
	})
})
