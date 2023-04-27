package grstack_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Clear JSON", Label("misc", "json.clear"), func() {

	It("can clear a simple value", func() {
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

	It("can clear multiple values", func() {
		cmd1 := client.JSONSet(ctx, "clear2", "$", `{"a": [100, 200, 300, 200], "b": {"a": [100, 200, 300, 200]}}`)
		Expect(cmd1.Err()).NotTo(HaveOccurred())
		Expect(cmd1.Val()).To(Equal("OK"))

		cmd2 := client.JSONClear(ctx, "clear2", "$..a")
		Expect(cmd2.Err()).NotTo(HaveOccurred())
		Expect(cmd2.Val()).To(Equal(int64(2)))

		cmd3 := client.JSONGet(ctx, "clear2", "$.a")
		Expect(cmd3.Err()).NotTo(HaveOccurred())
		Expect(cmd3.Val()).To(HaveLen(1))
		Expect(cmd3.Val()[0]).To(BeAssignableToTypeOf([]interface{}{1}))
		Expect(cmd3.Val()[0]).To(Equal([]interface{}{}))

		cmd3 = client.JSONGet(ctx, "clear2", "$.b.a")
		Expect(cmd3.Err()).NotTo(HaveOccurred())
		Expect(cmd3.Val()).To(HaveLen(1))
		Expect(cmd3.Val()[0]).To(BeAssignableToTypeOf([]interface{}{1}))
		Expect(cmd3.Val()[0]).To(Equal([]interface{}{}))

	})

})

var _ = Describe("Delete JSON", Label("misc"), func() {

	It("can delete the root of a key", Label("json.del"), func() {
		cmd1 := client.JSONSet(ctx, "del1", "$", `[1]`)
		Expect(cmd1.Err()).NotTo(HaveOccurred())
		Expect(cmd1.Val()).To(Equal("OK"))

		cmd2 := client.JSONDel(ctx, "del1", "$")
		Expect(cmd2.Err()).NotTo(HaveOccurred())
		Expect(cmd2.Val()).To(Equal(int64(1)))

		cmd3 := client.JSONGet(ctx, "del1", "$")
		Expect(cmd3.Err()).NotTo(HaveOccurred())
		Expect(cmd3.Val()).To(HaveLen(0))
	})

	It("can delete a single value in a key", Label("json.del"), func() {
		cmd1 := client.JSONSet(ctx, "del2", "$", `{"a": [1,2,3], "b": {"a": [1,2,3], "b": "annie"}}`)
		Expect(cmd1.Err()).NotTo(HaveOccurred())
		Expect(cmd1.Val()).To(Equal("OK"))

		cmd2 := client.JSONDel(ctx, "del2", "$.a")
		Expect(cmd2.Err()).NotTo(HaveOccurred())
		Expect(cmd2.Val()).To(Equal(int64(1)))

		cmd3 := client.JSONGet(ctx, "del2", "$")
		Expect(cmd3.Err()).NotTo(HaveOccurred())
		Expect(cmd3.Val()).To(HaveLen(1))
		Expect(cmd3.Val()[0]).To(BeAssignableToTypeOf(map[string]interface{}{"foo": "bar"}))
		Expect(cmd3.Val()[0]).To(Equal(map[string]interface{}{
			"b": map[string]interface{}{
				"a": []interface{}{float64(1), float64(2), float64(3)},
				"b": "annie",
			},
		}))

	})

	It("can delete multiple values in a key", Label("json.del"), func() {
		cmd1 := client.JSONSet(ctx, "del3", "$", `{"a": [1,2,3], "b": {"a": [1,2,3], "b": "annie"}}`)
		Expect(cmd1.Err()).NotTo(HaveOccurred())
		Expect(cmd1.Val()).To(Equal("OK"))

		cmd2 := client.JSONDel(ctx, "del3", "$..a")
		Expect(cmd2.Err()).NotTo(HaveOccurred())
		Expect(cmd2.Val()).To(Equal(int64(2)))

		cmd3 := client.JSONGet(ctx, "del3", "$")
		Expect(cmd3.Err()).NotTo(HaveOccurred())
		Expect(cmd3.Val()).To(HaveLen(1))
		Expect(cmd3.Val()[0]).To(BeAssignableToTypeOf(map[string]interface{}{"foo": "bar"}))
		Expect(cmd3.Val()[0]).To(Equal(map[string]interface{}{
			"b": map[string]interface{}{
				"b": "annie",
			},
		}))

	})

	It("can forget the root of a key", Label("json.forget"), func() {
		cmd1 := client.JSONSet(ctx, "forget1", "$", `[1]`)
		Expect(cmd1.Err()).NotTo(HaveOccurred())
		Expect(cmd1.Val()).To(Equal("OK"))

		cmd2 := client.JSONForget(ctx, "forget1", "$")
		Expect(cmd2.Err()).NotTo(HaveOccurred())
		Expect(cmd2.Val()).To(Equal(int64(1)))

		cmd3 := client.JSONGet(ctx, "forget1", "$")
		Expect(cmd3.Err()).NotTo(HaveOccurred())
		Expect(cmd3.Val()).To(HaveLen(0))
	})

	It("can forget a single value in a key", Label("json.forget"), func() {
		cmd1 := client.JSONSet(ctx, "forget2", "$", `{"a": [1,2,3], "b": {"a": [1,2,3], "b": "annie"}}`)
		Expect(cmd1.Err()).NotTo(HaveOccurred())
		Expect(cmd1.Val()).To(Equal("OK"))

		cmd2 := client.JSONForget(ctx, "forget2", "$.a")
		Expect(cmd2.Err()).NotTo(HaveOccurred())
		Expect(cmd2.Val()).To(Equal(int64(1)))

		cmd3 := client.JSONGet(ctx, "forget2", "$")
		Expect(cmd3.Err()).NotTo(HaveOccurred())
		Expect(cmd3.Val()).To(HaveLen(1))
		Expect(cmd3.Val()[0]).To(BeAssignableToTypeOf(map[string]interface{}{"foo": "bar"}))
		Expect(cmd3.Val()[0]).To(Equal(map[string]interface{}{
			"b": map[string]interface{}{
				"a": []interface{}{float64(1), float64(2), float64(3)},
				"b": "annie",
			},
		}))

	})

	It("can forget multiple values in a key", Label("json.forget"), func() {
		cmd1 := client.JSONSet(ctx, "forget3", "$", `{"a": [1,2,3], "b": {"a": [1,2,3], "b": "annie"}}`)
		Expect(cmd1.Err()).NotTo(HaveOccurred())
		Expect(cmd1.Val()).To(Equal("OK"))

		cmd2 := client.JSONForget(ctx, "forget3", "$..a")
		Expect(cmd2.Err()).NotTo(HaveOccurred())
		Expect(cmd2.Val()).To(Equal(int64(2)))

		cmd3 := client.JSONGet(ctx, "forget3", "$")
		Expect(cmd3.Err()).NotTo(HaveOccurred())
		Expect(cmd3.Val()).To(HaveLen(1))
		Expect(cmd3.Val()[0]).To(BeAssignableToTypeOf(map[string]interface{}{"foo": "bar"}))
		Expect(cmd3.Val()[0]).To(Equal(map[string]interface{}{
			"b": map[string]interface{}{
				"b": "annie",
			},
		}))

	})

})
