package grstack_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("JSON Array Append", func() {
	It("can append a single value to a simple array", func() {
		cmd1 := client.JSONSet(ctx, "append1", "$", `[]`)
		Expect(cmd1.Err()).NotTo(HaveOccurred())
		Expect(cmd1.Val()).To(Equal("OK"))

		cmd2 := client.JSONArrAppend(ctx, "append1", "$", 10)
		Expect(cmd2.Err()).NotTo(HaveOccurred())
		Expect(cmd2.Val()).To(Equal([]int64{1}))
	})

	It("can append a value in multiple locations in an object", func() {
		cmd1 := client.JSONSet(ctx, "append2", "$", `{"a": [10], "b": {"a": [12, 13]}}`)
		Expect(cmd1.Err()).NotTo(HaveOccurred())
		Expect(cmd1.Val()).To(Equal("OK"))

		cmd2 := client.JSONArrAppend(ctx, "append2", "$..a", 10)
		Expect(cmd2.Err()).NotTo(HaveOccurred())
		Expect(cmd2.Val()).To(Equal([]int64{2, 3}))
	})

})

var _ = Describe("JSON Array Index", func() {
	It("can get the single index of a value", func() {
		cmd1 := client.JSONSet(ctx, "index1", "$", `[100, 200, 300, 200]`)
		Expect(cmd1.Err()).NotTo(HaveOccurred())
		Expect(cmd1.Val()).To(Equal("OK"))

		cmd2 := client.JSONArrIndex(ctx, "index1", "$", 200)
		Expect(cmd2.Err()).NotTo(HaveOccurred())
		Expect(cmd2.Val()).To(Equal([]int64{1}))
	})

	It("can find a value in multiple arrays", func() {
		cmd1 := client.JSONSet(ctx, "index2", "$", `{"a": [10], "b": {"a": [12, 10]}}`)
		Expect(cmd1.Err()).NotTo(HaveOccurred())
		Expect(cmd1.Val()).To(Equal("OK"))

		cmd2 := client.JSONArrIndex(ctx, "index2", "$..a", 10)
		Expect(cmd2.Err()).NotTo(HaveOccurred())
		Expect(cmd2.Val()).To(Equal([]int64{0, 1}))
	})

	It("can find a value a nested array", func() {
		cmd1 := client.JSONSet(ctx, "index3", "$", `{"a": [10], "b": {"a": [12, 10]}}`)
		Expect(cmd1.Err()).NotTo(HaveOccurred())
		Expect(cmd1.Val()).To(Equal("OK"))

		cmd2 := client.JSONArrIndex(ctx, "index3", "$.b.a", 10)
		Expect(cmd2.Err()).NotTo(HaveOccurred())
		Expect(cmd2.Val()).To(Equal([]int64{1}))
	})

	It("can find a value given a start/top range", func() {
		cmd1 := client.JSONSet(ctx, "index4", "$", `{"a": [10], "b": {"a": [12, 10, 20, 12, 90, 10]}}`)
		Expect(cmd1.Err()).NotTo(HaveOccurred())
		Expect(cmd1.Val()).To(Equal("OK"))

		cmd2 := client.JSONArrIndexStartStop(ctx, "index4", "$.b.a", 12, 1, 4)
		Expect(cmd2.Err()).NotTo(HaveOccurred())
		Expect(cmd2.Val()).To(Equal([]int64{3}))
	})

})

var _ = Describe("JSON Array Insert", func() {
	It("can insert a single value", func() {
		cmd1 := client.JSONSet(ctx, "insert1", "$", `[100, 200, 300, 200]`)
		Expect(cmd1.Err()).NotTo(HaveOccurred())
		Expect(cmd1.Val()).To(Equal("OK"))

		cmd2 := client.JSONArrInsert(ctx, "insert1", "$", 0, 200)
		Expect(cmd2.Err()).NotTo(HaveOccurred())
		Expect(cmd2.Val()).To(Equal([]int64{5}))

		cmd3 := client.JSONGet(ctx, "insert1")
		Expect(cmd3.Err()).NotTo(HaveOccurred())
		Expect(cmd3.Val()).To(Equal([]interface{}{float64(200), float64(100), float64(200), float64(300), float64(200)}))

	})

	It("can insert multiple values", func() {
		cmd1 := client.JSONSet(ctx, "insert2", "$", `[100, 200, 300, 200]`)
		Expect(cmd1.Err()).NotTo(HaveOccurred())
		Expect(cmd1.Val()).To(Equal("OK"))

		cmd2 := client.JSONArrInsert(ctx, "insert2", "$", -1, 1, 2)
		Expect(cmd2.Err()).NotTo(HaveOccurred())
		Expect(cmd2.Val()).To(Equal([]int64{6}))

		cmd3 := client.JSONGet(ctx, "insert2")
		Expect(cmd3.Err()).NotTo(HaveOccurred())
		Expect(cmd3.Val()).To(Equal([]interface{}{float64(100), float64(200), float64(300), float64(1), float64(2), float64(200)}))
	})

	It("can't insert values if the path doesn't match", func() {
		cmd1 := client.JSONSet(ctx, "insert3", "$", `[100, 200, 300, 200]`)
		Expect(cmd1.Err()).NotTo(HaveOccurred())
		Expect(cmd1.Val()).To(Equal("OK"))

		cmd2 := client.JSONArrInsert(ctx, "insert3", "$.a", 0, 1)
		Expect(cmd2.Err()).NotTo(HaveOccurred())
		Expect(len(cmd2.Val())).To(Equal(0))

		cmd3 := client.JSONGet(ctx, "insert3")
		Expect(cmd3.Err()).NotTo(HaveOccurred())
		Expect(cmd3.Val()).To(Equal([]interface{}{float64(100), float64(200), float64(300), float64(200)}))
	})

	It("can insert into multiple arrays", func() {
		cmd1 := client.JSONSet(ctx, "insert4", "$", `{"a": [10], "b": {"a": [12, 10, 20, 12, 90, 10]}}`)
		Expect(cmd1.Err()).NotTo(HaveOccurred())
		Expect(cmd1.Val()).To(Equal("OK"))

		cmd2 := client.JSONArrInsert(ctx, "insert4", "$..a", 0, 1)
		Expect(cmd2.Err()).NotTo(HaveOccurred())
		Expect(cmd2.Val()).To(Equal([]int64{2, 7}))

		cmd3 := client.JSONGet(ctx, "insert4", "$.a[0]")
		Expect(cmd3.Err()).NotTo(HaveOccurred())
		Expect(cmd3.Val()).To(Equal([]interface{}{float64(1)}))

		cmd4 := client.JSONGet(ctx, "insert4", "$.b.a[0]")
		Expect(cmd4.Err()).NotTo(HaveOccurred())
		Expect(cmd4.Val()).To(Equal([]interface{}{float64(1)}))
	})

})

var _ = Describe("JSON Array Lengths", func() {
	It("can get the length of a single array", func() {
		cmd1 := client.JSONSet(ctx, "length1", "$", `[100, 200, 300, 200]`)
		Expect(cmd1.Err()).NotTo(HaveOccurred())
		Expect(cmd1.Val()).To(Equal("OK"))

		cmd2 := client.JSONArrLen(ctx, "length1", "$")
		Expect(cmd2.Err()).NotTo(HaveOccurred())
		Expect(cmd2.Val()).To(Equal([]int64{4}))

	})

	It("can get the length of multiple arrays", func() {
		cmd1 := client.JSONSet(ctx, "length2", "$", `{"a": [10], "b": {"a": [12, 10, 20, 12, 90, 10]}}`)
		Expect(cmd1.Err()).NotTo(HaveOccurred())
		Expect(cmd1.Val()).To(Equal("OK"))

		cmd2 := client.JSONArrLen(ctx, "length2", "$..a")
		Expect(cmd2.Err()).NotTo(HaveOccurred())
		Expect(cmd2.Val()).To(Equal([]int64{1, 6}))
	})

	It("can't get the length if the path doesnt match", func() {
		cmd1 := client.JSONSet(ctx, "length3", "$", `[100, 200, 300, 200]`)
		Expect(cmd1.Err()).NotTo(HaveOccurred())
		Expect(cmd1.Val()).To(Equal("OK"))

		cmd2 := client.JSONArrLen(ctx, "length3", "$.a")
		Expect(cmd2.Err()).NotTo(HaveOccurred())
		Expect(len(cmd2.Val())).To(Equal(0))
	})

})

var _ = Describe("JSON Array Popping", func() {
	It("remove one item from the end of an array", func() {
		cmd1 := client.JSONSet(ctx, "pop1", "$", `[100, 200, 300, 200]`)
		Expect(cmd1.Err()).NotTo(HaveOccurred())
		Expect(cmd1.Val()).To(Equal("OK"))

		cmd2 := client.JSONArrPop(ctx, "pop1", "$", -1)
		Expect(cmd2.Err()).NotTo(HaveOccurred())
		Expect(cmd2.Val()).To(Equal([]string{"200"}))

		// Remember that JSON.Get returns a slice. Our result is the first
		// item in that slice.
		cmd3 := client.JSONGet(ctx, "pop1", "$")
		Expect(cmd3.Err()).NotTo(HaveOccurred())
		Expect(cmd3.Val()).To(HaveLen(1))
		Expect(cmd3.Val()[0]).To(BeAssignableToTypeOf([]interface{}{1}))
		Expect(cmd3.Val()[0]).To(Equal(([]interface{}{float64(100), float64(200), float64(300)})))
		Expect(cmd3.Val()[0]).To(HaveLen(3))

	})

	It("can remove one item from the start of an array", func() {
		cmd1 := client.JSONSet(ctx, "pop2", "$", `[100, 200, 300, 200]`)
		Expect(cmd1.Err()).NotTo(HaveOccurred())
		Expect(cmd1.Val()).To(Equal("OK"))

		cmd2 := client.JSONArrPop(ctx, "pop2", "$", 0)
		Expect(cmd2.Err()).NotTo(HaveOccurred())
		Expect(cmd2.Val()).To(Equal([]string{"100"}))

		// Remember that JSON.Get returns a slice. Our result is the first
		// item in that slice.
		cmd3 := client.JSONGet(ctx, "pop2", "$")
		Expect(cmd3.Err()).NotTo(HaveOccurred())
		Expect(cmd3.Val()).To(HaveLen(1))
		Expect(cmd3.Val()[0]).To(BeAssignableToTypeOf([]interface{}{1}))

		Expect(cmd3.Val()[0]).To(Equal([]interface{}{float64(200), float64(300), float64(200)}))
		Expect(cmd3.Val()[0]).To(HaveLen(3))
	})

	It("can't get remove an item if the path doesnt match", func() {
		cmd1 := client.JSONSet(ctx, "pop3", "$", `[100, 200, 300, 200]`)
		Expect(cmd1.Err()).NotTo(HaveOccurred())
		Expect(cmd1.Val()).To(Equal("OK"))

		cmd2 := client.JSONArrPop(ctx, "pop3", "$.a", 1)
		Expect(cmd2.Err()).NotTo(HaveOccurred())
		Expect(len(cmd2.Val())).To(Equal(0))
	})

	It("can remove items from the middle of structures", func() {
		cmd1 := client.JSONSet(ctx, "pop4", "$", `[100, 200, 300, 200]`)
		Expect(cmd1.Err()).NotTo(HaveOccurred())
		Expect(cmd1.Val()).To(Equal("OK"))

		cmd2 := client.JSONArrPop(ctx, "pop4", "$", 2)
		Expect(cmd2.Err()).NotTo(HaveOccurred())
		Expect(cmd2.Val()).To(Equal([]string{"300"}))

		cmd3 := client.JSONGet(ctx, "pop4", "$")
		Expect(cmd3.Err()).NotTo(HaveOccurred())
		Expect(cmd3.Val()).To(HaveLen(1))
		Expect(cmd3.Val()[0]).To(BeAssignableToTypeOf([]interface{}{1}))
		Expect(cmd3.Val()[0]).To(Equal([]interface{}{float64(100), float64(200), float64(200)}))
	})

	It("can remove items from complex structures", func() {
		cmd1 := client.JSONSet(ctx, "pop5", "$", `{"a": [100, 200, 300, 200], "b": {"a": [100, 200, 300, 200]}}`)
		Expect(cmd1.Err()).NotTo(HaveOccurred())
		Expect(cmd1.Val()).To(Equal("OK"))

		cmd2 := client.JSONArrPop(ctx, "pop5", "$..a", 1)
		Expect(cmd2.Err()).NotTo(HaveOccurred())
		Expect(cmd2.Val()).To(Equal([]string{"200", "200"}))

		cmd3 := client.JSONGet(ctx, "pop5", "$.a")
		Expect(cmd3.Err()).NotTo(HaveOccurred())
		Expect(cmd3.Val()).To(HaveLen(1))
		Expect(cmd3.Val()[0]).To(BeAssignableToTypeOf([]interface{}{1}))
		Expect(cmd3.Val()[0]).To(Equal([]interface{}{float64(100), float64(300), float64(200)}))

		cmd3 = client.JSONGet(ctx, "pop5", "$.b.a")
		Expect(cmd3.Err()).NotTo(HaveOccurred())
		Expect(cmd3.Val()).To(HaveLen(1))
		Expect(cmd3.Val()[0]).To(BeAssignableToTypeOf([]interface{}{1}))
		Expect(cmd3.Val()[0]).To(Equal([]interface{}{float64(100), float64(300), float64(200)}))

	})

})

var _ = Describe("JSON Array Trimming", func() {

	It("can trim a simple array by one element", func() {
		cmd1 := client.JSONSet(ctx, "trim1", "$", `[100, 200, 300, 200]`)
		Expect(cmd1.Err()).NotTo(HaveOccurred())
		Expect(cmd1.Val()).To(Equal("OK"))

		cmd2 := client.JSONArrTrim(ctx, "trim1", "$", 0, 2)
		Expect(cmd2.Err()).NotTo(HaveOccurred())
		Expect(cmd2.Val()).To(Equal([]int64{3}))

		cmd3 := client.JSONGet(ctx, "trim1", "$")
		Expect(cmd3.Err()).NotTo(HaveOccurred())
		Expect(cmd3.Val()).To(HaveLen(1))
		Expect(cmd3.Val()[0]).To(BeAssignableToTypeOf([]interface{}{1}))
		Expect(cmd3.Val()[0]).To(Equal([]interface{}{float64(100), float64(200), float64(300)}))
	})

	It("can trim multiple values from the start of an array", func() {
		cmd1 := client.JSONSet(ctx, "trim2", "$", `[100, 200, 300, 200]`)
		Expect(cmd1.Err()).NotTo(HaveOccurred())
		Expect(cmd1.Val()).To(Equal("OK"))

		cmd2 := client.JSONArrTrim(ctx, "trim2", "$", 2, 3)
		Expect(cmd2.Err()).NotTo(HaveOccurred())
		Expect(cmd2.Val()).To(Equal([]int64{2}))

		cmd3 := client.JSONGet(ctx, "trim2", "$")
		Expect(cmd3.Err()).NotTo(HaveOccurred())
		Expect(cmd3.Val()).To(HaveLen(1))
		Expect(cmd3.Val()[0]).To(BeAssignableToTypeOf([]interface{}{1}))
		Expect(cmd3.Val()[0]).To(Equal([]interface{}{float64(300), float64(200)}))
	})

	It("can trim multiple values from the end of an array", func() {
		cmd1 := client.JSONSet(ctx, "trim3", "$", `[100, 200, 300, 200]`)
		Expect(cmd1.Err()).NotTo(HaveOccurred())
		Expect(cmd1.Val()).To(Equal("OK"))

		cmd2 := client.JSONArrTrim(ctx, "trim3", "$", 2, 3)
		Expect(cmd2.Err()).NotTo(HaveOccurred())
		Expect(cmd2.Val()).To(Equal([]int64{2}))

		cmd3 := client.JSONGet(ctx, "trim3", "$")
		Expect(cmd3.Err()).NotTo(HaveOccurred())
		Expect(cmd3.Val()).To(HaveLen(1))
		Expect(cmd3.Val()[0]).To(BeAssignableToTypeOf([]interface{}{1}))
		Expect(cmd3.Val()[0]).To(Equal([]interface{}{float64(300), float64(200)}))
	})

	It("can trim in the middle of an array", func() {
		cmd1 := client.JSONSet(ctx, "trim4", "$", `[100, 200, 300, 200]`)
		Expect(cmd1.Err()).NotTo(HaveOccurred())
		Expect(cmd1.Val()).To(Equal("OK"))

		cmd2 := client.JSONArrTrim(ctx, "trim4", "$", 1, 2)
		Expect(cmd2.Err()).NotTo(HaveOccurred())
		Expect(cmd2.Val()).To(Equal([]int64{2}))

		cmd3 := client.JSONGet(ctx, "trim4", "$")
		Expect(cmd3.Err()).NotTo(HaveOccurred())
		Expect(cmd3.Val()).To(HaveLen(1))
		Expect(cmd3.Val()[0]).To(BeAssignableToTypeOf([]interface{}{1}))
		Expect(cmd3.Val()[0]).To(Equal([]interface{}{float64(200), float64(300)}))
	})

	It("can trim when multiple paths match", func() {
		cmd1 := client.JSONSet(ctx, "trim5", "$", `{"a": [100, 200, 300, 200], "b": {"a": [100, 200, 300, 200]}}`)
		Expect(cmd1.Err()).NotTo(HaveOccurred())
		Expect(cmd1.Val()).To(Equal("OK"))

		cmd2 := client.JSONArrTrim(ctx, "trim5", "$..a", 1, 2)
		Expect(cmd2.Err()).NotTo(HaveOccurred())
		Expect(cmd2.Val()).To(Equal([]int64{2, 2}))

		cmd3 := client.JSONGet(ctx, "trim5", "$.a")
		Expect(cmd3.Err()).NotTo(HaveOccurred())
		Expect(cmd3.Val()).To(HaveLen(1))
		Expect(cmd3.Val()[0]).To(BeAssignableToTypeOf([]interface{}{1}))
		Expect(cmd3.Val()[0]).To(Equal([]interface{}{float64(200), float64(300)}))

		cmd3 = client.JSONGet(ctx, "trim5", "$.b.a")
		Expect(cmd3.Err()).NotTo(HaveOccurred())
		Expect(cmd3.Val()).To(HaveLen(1))
		Expect(cmd3.Val()[0]).To(BeAssignableToTypeOf([]interface{}{1}))
		Expect(cmd3.Val()[0]).To(Equal([]interface{}{float64(200), float64(300)}))
	})

})
