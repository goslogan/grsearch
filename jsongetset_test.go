package grstack_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type TestResult struct {
	Hello string `json:"hello"`
}

var _ = Describe("JSON Set", func() {
	It("can add a new document", func() {
		cmd := client.JSONSet(ctx, "testdoc1", "$", `{"a": 1, "b": 2, "hello": "world"}`)
		Expect(cmd.Err()).NotTo(HaveOccurred())
		Expect(cmd.Val()).To(Equal("OK"))
	})

	It("can update a docuiment", func() {
		cmd1 := client.JSONSet(ctx, "testdoc2", "$", `{"a": 1, "b": 2, "hello": "world"}`)
		Expect(cmd1.Err()).NotTo(HaveOccurred())
		Expect(cmd1.Val()).To(Equal("OK"))

		cmd2 := client.JSONSet(ctx, "testdoc2", "$.hello", `"WORLD"`)
		Expect(cmd2.Err()).NotTo(HaveOccurred())
		Expect(cmd2.Val()).To(Equal("OK"))

	})
})

var _ = Describe("JSON Get", func() {
	It("can retrieve a single result from a document", func() {
		cmd1 := client.JSONSet(ctx, "testdoc3", "$", `{"a": 1, "b": 2, "hello": "world"}`)
		Expect(cmd1.Err()).NotTo(HaveOccurred())
		Expect(cmd1.Val()).To(Equal("OK"))

		cmd2 := client.JSONGet(ctx, "testdoc3", "$")
		Expect(cmd2.Err()).NotTo(HaveOccurred())
		Expect(len(cmd2.Val())).To(Equal(1))
	})

	It("can retrieve multiple results from a document", func() {
		cmd1 := client.JSONSet(ctx, "testdoc3", "$", `{"a": 1, "b": 2, "hello": "world"}`)
		Expect(cmd1.Err()).NotTo(HaveOccurred())
		Expect(cmd1.Val()).To(Equal("OK"))

		cmd2 := client.JSONGet(ctx, "testdoc3", "$.*")
		Expect(cmd2.Err()).NotTo(HaveOccurred())
		Expect(len(cmd2.Val())).To(Equal(3))
	})

	It("can retrieve complex results from a document", func() {
		cmd1 := client.JSONSet(ctx, "testdoc4", "$", `{"a": 1, "b": 2, "c": {"hello": "world"}}`)
		Expect(cmd1.Err()).NotTo(HaveOccurred())
		Expect(cmd1.Val()).To(Equal("OK"))

		cmd2 := client.JSONGet(ctx, "testdoc4", "$.*")
		Expect(cmd2.Err()).NotTo(HaveOccurred())
		Expect(len(cmd2.Val())).To(Equal(3))
	})

	It("can scan a result from a document", func() {
		cmd1 := client.JSONSet(ctx, "testdoc5", "$", `{"a": 1, "b": 2, "c": {"hello": "golang"}}`)
		Expect(cmd1.Err()).NotTo(HaveOccurred())
		Expect(cmd1.Val()).To(Equal("OK"))

		cmd2 := client.JSONGet(ctx, "testdoc5", "$.*")
		Expect(cmd2.Err()).NotTo(HaveOccurred())
		Expect(len(cmd2.Val())).To(Equal(3))

		test := TestResult{}
		err := cmd2.Scan(2, &test)
		Expect(err).NotTo(HaveOccurred())
		Expect(test.Hello).To(Equal("golang"))
	})

})
