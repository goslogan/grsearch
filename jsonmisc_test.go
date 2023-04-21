package grstack_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Misc JSON", Label("misc"), func() {

	It("can clear a simple value", Label("json.set"), func() {
		cmd1 := client.JSONSet(ctx, "clear1", "$", `[1]`)
		Expect(cmd1.Err()).NotTo(HaveOccurred())
		Expect(cmd1.Val()).To(Equal("OK"))

		cmd2 := client.JSONClear(ctx, "clear1", "$")
		Expect(cmd2.Err()).NotTo(HaveOccurred())
		Expect(cmd2.Val()).To(Equal(int64(1)))

		cmd3 := client.JSONGet(ctx, "clear1", "$")
		Expect(cmd3.Err()).NotTo(HaveOccurred())
		Expect(cmd3.Val()).To(HaveLen(1))
		Expect(cmd3.Val()[0]).To(BeAssignableToTypeOf([]interface{}{1}))
		Expect(cmd3.Val()[0]).To(Equal([]interface{}{}))

	})

})
