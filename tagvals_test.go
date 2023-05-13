package grstack_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tagvals", func() {

	It("can get tag values from the index", func() {
		cmd := client.FTTagVals(ctx, "hcustomers", "owner")
		Expect(cmd.Err()).NotTo(HaveOccurred())

		vals := cmd.Val()
		Expect(vals).To(ContainElements([]string{"lara.croft", "ellen.ripley", "sarah.oconnor"}))
		Expect(len(vals)).To(Equal(3))
	})
})
