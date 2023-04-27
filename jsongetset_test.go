package grstack_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type TestResult struct {
	Hello string `json:"hello"`
}

var _ = Describe("JSON Set", Label("getset", "json.set"), func() {
	It("can add a new document", func() {
		cmd := client.JSONSet(ctx, "set1", "$", `{"a": 1, "b": 2, "hello": "world"}`)
		Expect(cmd.Err()).NotTo(HaveOccurred())
		Expect(cmd.Val()).To(Equal("OK"))
	})

	It("can update a docuiment", func() {
		cmd1 := client.JSONSet(ctx, "set2", "$", `{"a": 1, "b": 2, "hello": "world"}`)
		Expect(cmd1.Err()).NotTo(HaveOccurred())
		Expect(cmd1.Val()).To(Equal("OK"))

		cmd2 := client.JSONSet(ctx, "set2", "$.hello", `"WORLD"`)
		Expect(cmd2.Err()).NotTo(HaveOccurred())
		Expect(cmd2.Val()).To(Equal("OK"))

	})
})

var _ = Describe("JSON Get", Label("getset", "json.get"), func() {
	It("can retrieve a single result from a document", func() {
		cmd1 := client.JSONSet(ctx, "get1", "$", `{"a": 1, "b": 2, "hello": "world"}`)
		Expect(cmd1.Err()).NotTo(HaveOccurred())
		Expect(cmd1.Val()).To(Equal("OK"))

		cmd2 := client.JSONGet(ctx, "get1", "$")
		Expect(cmd2.Err()).NotTo(HaveOccurred())
		Expect(len(cmd2.Val())).To(Equal(1))
	})

	It("can retrieve multiple results from a document", func() {
		cmd1 := client.JSONSet(ctx, "get2", "$", `{"a": 1, "b": 2, "hello": "world"}`)
		Expect(cmd1.Err()).NotTo(HaveOccurred())
		Expect(cmd1.Val()).To(Equal("OK"))

		cmd2 := client.JSONGet(ctx, "get2", "$.*")
		Expect(cmd2.Err()).NotTo(HaveOccurred())
		Expect(len(cmd2.Val())).To(Equal(3))
	})

	It("can retrieve complex results from a document", func() {
		cmd1 := client.JSONSet(ctx, "get3", "$", `{"a": 1, "b": 2, "c": {"hello": "world"}}`)
		Expect(cmd1.Err()).NotTo(HaveOccurred())
		Expect(cmd1.Val()).To(Equal("OK"))

		cmd2 := client.JSONGet(ctx, "get3", "$.*")
		Expect(cmd2.Err()).NotTo(HaveOccurred())
		Expect(len(cmd2.Val())).To(Equal(3))
	})

	It("can scan a result from a document", func() {
		cmd1 := client.JSONSet(ctx, "get4", "$", `{"a": 1, "b": 2, "c": {"hello": "golang"}}`)
		Expect(cmd1.Err()).NotTo(HaveOccurred())
		Expect(cmd1.Val()).To(Equal("OK"))

		cmd2 := client.JSONGet(ctx, "get4", "$.*")
		Expect(cmd2.Err()).NotTo(HaveOccurred())
		Expect(len(cmd2.Val())).To(Equal(3))

		test := TestResult{}
		err := cmd2.Scan(2, &test)
		Expect(err).NotTo(HaveOccurred())
		Expect(test.Hello).To(Equal("golang"))
	})

})

var _ = Describe("JSON MGet", Label("getset", "json.mget"), func() {

	It("can get a single value from multiple keys", func() {
		cmd1 := client.JSONSet(ctx, "mget1a", "$", `{"a": "alice"}`)
		Expect(cmd1.Err()).NotTo(HaveOccurred())
		Expect(cmd1.Val()).To(Equal("OK"))
		cmd2 := client.JSONSet(ctx, "mget1b", "$", `{"a": "bob"}`)
		Expect(cmd2.Err()).NotTo(HaveOccurred())
		Expect(cmd2.Val()).To(Equal("OK"))

		cmd3 := client.JSONMGet(ctx, "$", "mget1a", "mget1b")
		Expect(cmd3.Err()).NotTo(HaveOccurred())
		Expect(cmd3.Val()).To(HaveLen(2))
		Expect(cmd3.Val()[0]).To(BeAssignableToTypeOf([]interface{}{1}))
		Expect(cmd3.Val()[0]).To(Equal([]interface{}{map[string]interface{}{"a": "alice"}}))
		Expect(cmd3.Val()[1]).To(BeAssignableToTypeOf([]interface{}{1}))
		Expect(cmd3.Val()[1]).To(Equal([]interface{}{map[string]interface{}{"a": "bob"}}))
	})

	It("can get a multiple values from multiple keys", func() {
		cmd1 := client.JSONSet(ctx, "mget2a", "$", `{"a": ["aa", "ab", "ac", "ad"], "b": {"a": ["ba", "bb", "bc", "bd"]}}`)
		Expect(cmd1.Err()).NotTo(HaveOccurred())
		Expect(cmd1.Val()).To(Equal("OK"))
		cmd2 := client.JSONSet(ctx, "mget2b", "$", `{"a": [100, 200, 300, 200], "b": {"a": [100, 200, 300, 200]}}`)
		Expect(cmd2.Err()).NotTo(HaveOccurred())
		Expect(cmd2.Val()).To(Equal("OK"))

		cmd3 := client.JSONMGet(ctx, "$..a", "mget2a", "mget2b")
		Expect(cmd3.Err()).NotTo(HaveOccurred())
		Expect(cmd3.Val()).To(HaveLen(2))
		Expect(cmd3.Val()[0]).To(BeAssignableToTypeOf([]interface{}{1}))
		Expect(cmd3.Val()[0]).To(Equal([]interface{}{
			[]interface{}{"aa", "ab", "ac", "ad"},
			[]interface{}{"ba", "bb", "bc", "bd"},
		}))
		Expect(cmd3.Val()[1]).To(BeAssignableToTypeOf([]interface{}{1}))
		Expect(cmd3.Val()[1]).To(Equal([]interface{}{
			[]interface{}{float64(100), float64(200), float64(300), float64(200)},
			[]interface{}{float64(100), float64(200), float64(300), float64(200)},
		}))
	})

	It("can get some values when only some docs match", func() {
		cmd1 := client.JSONSet(ctx, "mget3a", "$", `{"a": ["aa", "ab", "ac", "ad"], "b": {"a": ["ba", "bb", "bc", "bd"]}}`)
		Expect(cmd1.Err()).NotTo(HaveOccurred())
		Expect(cmd1.Val()).To(Equal("OK"))
		cmd2 := client.JSONSet(ctx, "mget3b", "$", `{"a": [100, 200, 300, 200], "b": {"a": [100, 200, 300, 200]}}`)
		Expect(cmd2.Err()).NotTo(HaveOccurred())
		Expect(cmd2.Val()).To(Equal("OK"))
		cmd3 := client.JSONSet(ctx, "mget3c", "$", `{"person": "bob"}`)
		Expect(cmd3.Err()).NotTo(HaveOccurred())
		Expect(cmd3.Val()).To(Equal("OK"))

		// mget4c does not exist
		cmd4 := client.JSONMGet(ctx, "$..a", "mget3a", "mget3b", "mget3c", "mget4c")
		Expect(cmd4.Err()).NotTo(HaveOccurred())
		Expect(cmd4.Val()).To(HaveLen(4))
		Expect(cmd4.Val()[0]).To(BeAssignableToTypeOf([]interface{}{1}))
		Expect(cmd4.Val()[0]).To(Equal([]interface{}{
			[]interface{}{"aa", "ab", "ac", "ad"},
			[]interface{}{"ba", "bb", "bc", "bd"},
		}))
		Expect(cmd4.Val()[1]).To(BeAssignableToTypeOf([]interface{}{1}))
		Expect(cmd4.Val()[1]).To(Equal([]interface{}{
			[]interface{}{float64(100), float64(200), float64(300), float64(200)},
			[]interface{}{float64(100), float64(200), float64(300), float64(200)},
		}))
		Expect(cmd4.Val()[2]).To(BeAssignableToTypeOf([]interface{}{1}))
		Expect(cmd4.Val()[2]).To(Equal([]interface{}{}))

		Expect(cmd4.Val()[3]).To(BeAssignableToTypeOf([]interface{}{1}))
		Expect(cmd4.Val()[3]).To(BeNil())

	})
})
