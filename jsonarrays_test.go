package grstack_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("JSON Array Append", func() {
	It("can append a single value to a simple array", func() {
		cmd1 := client.JSONSet(ctx, "testdoca1", "$", `[]`)
		Expect(cmd1.Err()).NotTo(HaveOccurred())
		Expect(cmd1.Val()).To(Equal("OK"))

		cmd2 := client.JSONArrAppend(ctx, "testdoca1", "$", 10)
		Expect(cmd2.Err()).NotTo(HaveOccurred())
		Expect(cmd2.Val()).To(Equal([]int64{1}))
	})

	It("can append a value in multiple locations in an object", func() {
		cmd1 := client.JSONSet(ctx, "testdoca2", "$", `{"a": [10], "b": {"a": [12, 13]}}`)
		Expect(cmd1.Err()).NotTo(HaveOccurred())
		Expect(cmd1.Val()).To(Equal("OK"))

		cmd2 := client.JSONArrAppend(ctx, "testdoca2", "$..a", 10)
		Expect(cmd2.Err()).NotTo(HaveOccurred())
		Expect(cmd2.Val()).To(Equal([]int64{2, 3}))
	})

})

var _ = Describe("JSON Array Index", func() {
	It("can get the single index of a value", func() {
		cmd1 := client.JSONSet(ctx, "testdoci1", "$", `[100, 200, 300, 200]`)
		Expect(cmd1.Err()).NotTo(HaveOccurred())
		Expect(cmd1.Val()).To(Equal("OK"))

		cmd2 := client.JSONArrIndex(ctx, "testdoci1", "$", 200)
		Expect(cmd2.Err()).NotTo(HaveOccurred())
		Expect(cmd2.Val()).To(Equal([]int64{1}))
	})

	It("can find a value in multiple arrays", func() {
		cmd1 := client.JSONSet(ctx, "testdoci2", "$", `{"a": [10], "b": {"a": [12, 10]}}`)
		Expect(cmd1.Err()).NotTo(HaveOccurred())
		Expect(cmd1.Val()).To(Equal("OK"))

		cmd2 := client.JSONArrIndex(ctx, "testdoci2", "$..a", 10)
		Expect(cmd2.Err()).NotTo(HaveOccurred())
		Expect(cmd2.Val()).To(Equal([]int64{0, 1}))
	})

	It("can find a value a nested array", func() {
		cmd1 := client.JSONSet(ctx, "testdoci3", "$", `{"a": [10], "b": {"a": [12, 10]}}`)
		Expect(cmd1.Err()).NotTo(HaveOccurred())
		Expect(cmd1.Val()).To(Equal("OK"))

		cmd2 := client.JSONArrIndex(ctx, "testdoci3", "$.b.a", 10)
		Expect(cmd2.Err()).NotTo(HaveOccurred())
		Expect(cmd2.Val()).To(Equal([]int64{1}))
	})

	It("can find a value given a start/top range", func() {
		cmd1 := client.JSONSet(ctx, "testdoci4", "$", `{"a": [10], "b": {"a": [12, 10, 20, 12, 90, 10]}}`)
		Expect(cmd1.Err()).NotTo(HaveOccurred())
		Expect(cmd1.Val()).To(Equal("OK"))

		cmd2 := client.JSONArrIndexStartStop(ctx, "testdoci4", "$.b.a", 12, 1, 4)
		Expect(cmd2.Err()).NotTo(HaveOccurred())
		Expect(cmd2.Val()).To(Equal([]int64{3}))
	})

})
