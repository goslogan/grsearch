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

var _ = Describe("Increment Numbers", Label("misc", "json.numincrby"), func() {

	It("can increment a single value by one", func() {
		cmd1 := client.JSONSet(ctx, "incr1", "$", "[1]")
		Expect(cmd1.Err()).NotTo(HaveOccurred())
		Expect(cmd1.Val()).To(Equal("OK"))

		cmd2 := client.JSONNumIncrBy(ctx, "incr1", "$.[0]", float64(1))
		Expect(cmd2.Err()).NotTo(HaveOccurred())
		Expect(cmd2.Val()).To(Equal([]interface{}{float64(2)}))
	})

	It("can increment a single value by a floating point value", func() {
		cmd1 := client.JSONSet(ctx, "incr2", "$", "[1]")
		Expect(cmd1.Err()).NotTo(HaveOccurred())
		Expect(cmd1.Val()).To(Equal("OK"))

		cmd2 := client.JSONNumIncrBy(ctx, "incr2", "$.[0]", float64(2.75))
		Expect(cmd2.Err()).NotTo(HaveOccurred())
		Expect(cmd2.Val()).To(Equal([]interface{}{float64(3.75)}))
	})

	It("can increment multiple paths in a single key", func() {
		cmd1 := client.JSONSet(ctx, "incr3", "$", `{"a": [1, 2], "b": {"a": [0, -1]}}`)
		Expect(cmd1.Err()).NotTo(HaveOccurred())
		Expect(cmd1.Val()).To(Equal("OK"))

		cmd2 := client.JSONNumIncrBy(ctx, "incr3", "$..a[1]", float64(1))
		Expect(cmd2.Err()).NotTo(HaveOccurred())
		Expect(cmd2.Val()).To(Equal([]interface{}{float64(3), float64(0)}))
	})

	It("can increment partially when one of the path matches is not a number", func() {
		cmd1 := client.JSONSet(ctx, "incr4", "$", `{"a": 3, "b": {"a": "redis"}, "c": {"a": 4}}`)
		Expect(cmd1.Err()).NotTo(HaveOccurred())
		Expect(cmd1.Val()).To(Equal("OK"))

		cmd2 := client.JSONNumIncrBy(ctx, "incr4", "$..a", float64(1))
		Expect(cmd2.Err()).NotTo(HaveOccurred())
		Expect(cmd2.Val()).To(Equal([]interface{}{float64(4), nil, float64(5)}))
	})
})

var _ = Describe("Get Object Keys", Label("misc", "json.objkeys"), func() {

	It("can get the keys in a simple object", func() {
		cmd1 := client.JSONSet(ctx, "objkeys1", "$", `{"a": [1, 2], "b": {"a": [0, -1]}}`)
		Expect(cmd1.Err()).NotTo(HaveOccurred())
		Expect(cmd1.Val()).To(Equal("OK"))

		cmd2 := client.JSONObjKeys(ctx, "objkeys1", "$")
		Expect(cmd2.Err()).NotTo(HaveOccurred())
		Expect(cmd2.Val()).To(HaveLen(1))
		Expect(cmd2.Val()[0]).To(BeAssignableToTypeOf([]interface{}{"a"}))
		Expect(cmd2.Val()[0]).To(Equal([]interface{}{"a", "b"}))
	})

	It("can get the keys in a complex object", func() {
		cmd1 := client.JSONSet(ctx, "objkeys1", "$", `{"a": [1, 2], "b": {"a": [0, -1]}}`)
		Expect(cmd1.Err()).NotTo(HaveOccurred())
		Expect(cmd1.Val()).To(Equal("OK"))

		cmd2 := client.JSONObjKeys(ctx, "objkeys1", "$..*")
		Expect(cmd2.Err()).NotTo(HaveOccurred())
		Expect(cmd2.Val()).To(HaveLen(7))
		Expect(cmd2.Val()).To(Equal([]interface{}{nil, []interface{}{"a"}, nil, nil, nil, nil, nil}))
	})

})

var _ = Describe("Get Object Length", Label("misc", "json.objlen"), func() {

	It("can get the number of keys in a simple object", func() {
		cmd1 := client.JSONSet(ctx, "objlen1", "$", `{"a": [1, 2], "b": {"a": [0, -1]}}`)
		Expect(cmd1.Err()).NotTo(HaveOccurred())
		Expect(cmd1.Val()).To(Equal("OK"))

		cmd2 := client.JSONObjLen(ctx, "objlen1", "$")
		Expect(cmd2.Err()).NotTo(HaveOccurred())
		Expect(cmd2.Val()).To(HaveLen(1))
		Expect(*cmd2.Val()[0]).To(Equal(int64(2)))
	})

	It("can get the number of keys in a complex object", func() {
		cmd1 := client.JSONSet(ctx, "objlen2", "$", `{"a": [1, 2], "b": {"a": [0, -1]}}`)
		Expect(cmd1.Err()).NotTo(HaveOccurred())
		Expect(cmd1.Val()).To(Equal("OK"))

		cmd2 := client.JSONObjLen(ctx, "objlen2", "$..*")
		Expect(cmd2.Err()).NotTo(HaveOccurred())
		Expect(cmd2.Val()).To(HaveLen(7))
		Expect(cmd2.Val()[0]).To(BeNil())
		Expect(*cmd2.Val()[1]).To(Equal(int64(1)))
	})

})

var _ = Describe("Get the length of strings", Label("misc", "json.strlen"), func() {

	It("can get the length of a single string", func() {
		cmd1 := client.JSONSet(ctx, "strlen1", "$", `{"a": "alice", "b": "bob"}`)
		Expect(cmd1.Err()).NotTo(HaveOccurred())
		Expect(cmd1.Val()).To(Equal("OK"))

		cmd2 := client.JSONStrLen(ctx, "strlen1", "$.a")
		Expect(cmd2.Err()).NotTo(HaveOccurred())
		Expect(cmd2.Val()).To(HaveLen(1))
		var tmp int64 = 20
		Expect(cmd2.Val()[0]).To(BeAssignableToTypeOf(&tmp))
		Expect(*cmd2.Val()[0]).To(Equal(int64(5)))
	})

	It("can get the length of all the strings", func() {
		cmd1 := client.JSONSet(ctx, "strlen2", "$", `{"a": "alice", "b": "bob", "c": {"a": "alice", "b": "bob"}}`)
		Expect(cmd1.Err()).NotTo(HaveOccurred())
		Expect(cmd1.Val()).To(Equal("OK"))

		cmd2 := client.JSONStrLen(ctx, "strlen2", "$..*")
		Expect(cmd2.Err()).NotTo(HaveOccurred())
		Expect(cmd2.Val()).To(HaveLen(5))
		var tmp int64 = 20
		Expect(cmd2.Val()[0]).To(BeAssignableToTypeOf(&tmp))
		Expect(*cmd2.Val()[0]).To(Equal(int64(5)))
		Expect(*cmd2.Val()[1]).To(Equal(int64(3)))
		Expect(cmd2.Val()[2]).To(BeNil())
		Expect(*cmd2.Val()[3]).To(Equal(int64(5)))
		Expect(*cmd2.Val()[4]).To(Equal(int64(3)))
	})

})

var _ = Describe("Append to Strings", Label("misc", "json.strappend"), func() {

	It("append to a simple string", func() {
		cmd1 := client.JSONSet(ctx, "strappend1", "$", `{"a": "alice", "b": "bob"}`)
		Expect(cmd1.Err()).NotTo(HaveOccurred())
		Expect(cmd1.Val()).To(Equal("OK"))

		cmd2 := client.JSONStrAppend(ctx, "strappend1", "$.a", `" bob"`)
		Expect(cmd2.Err()).NotTo(HaveOccurred())
		Expect(cmd2.Val()).To(HaveLen(1))
		var tmp int64 = 20
		Expect(cmd2.Val()[0]).To(BeAssignableToTypeOf(&tmp))
		Expect(*cmd2.Val()[0]).To(Equal(int64(9)))
	})

	It("append to multiple keys", func() {
		cmd1 := client.JSONSet(ctx, "strappend2", "$", `{"a": "alice", "b": "bob", "c": {"a": "alice", "b": "bob"}}`)
		Expect(cmd1.Err()).NotTo(HaveOccurred())
		Expect(cmd1.Val()).To(Equal("OK"))

		cmd2 := client.JSONStrAppend(ctx, "strappend2", "$..*", `" bob"`)
		Expect(cmd2.Err()).NotTo(HaveOccurred())
		Expect(cmd2.Val()).To(HaveLen(5))
		var tmp int64 = 20
		Expect(cmd2.Val()[0]).To(BeAssignableToTypeOf(&tmp))
		Expect(*cmd2.Val()[0]).To(Equal(int64(9)))
		Expect(*cmd2.Val()[1]).To(Equal(int64(7)))
		Expect(cmd2.Val()[2]).To(BeNil())
		Expect(*cmd2.Val()[3]).To(Equal(int64(9)))
		Expect(*cmd2.Val()[4]).To(Equal(int64(7)))
	})

})
