package ftsearch_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Dictionary", func() {
	It("can add entries to a dictionary", func() {
		cmd1 := client.FTDictDump(ctx, "testdict1") // should not exist
		Expect(cmd1.Err()).To(HaveOccurred())
		Expect(len(cmd1.Val())).To(BeZero())
		cmd2 := client.FTDictAdd(ctx, "testdict1", "foo", "bar")
		Expect(cmd2.Err()).NotTo(HaveOccurred())
		Expect(cmd2.Val()).To(BeNumerically("==", 2))
	})

	It("can remove entries from a dictionary", func() {
		cmd1 := client.FTDictDump(ctx, "testdict2")
		Expect(cmd1.Err()).To(HaveOccurred()) // should not exist
		Expect(len(cmd1.Val())).To(BeZero())
		cmd2 := client.FTDictAdd(ctx, "testdict2", "foo", "bar")
		Expect(cmd2.Err()).NotTo(HaveOccurred())
		Expect(cmd2.Val()).To(BeNumerically("==", 2))
		cmd3 := client.FTDictDel(ctx, "testdict2", "bar")
		Expect(cmd3.Err()).NotTo(HaveOccurred())
		Expect(cmd3.Val()).To(BeNumerically("==", 1))
		cmd4 := client.FTDictDump(ctx, "testdict2")
		Expect(cmd4.Err()).NotTo(HaveOccurred())
		Expect(cmd4.Val()).To(Equal([]string{"foo"}))
	})

	It("can dump a dictionary", func() {
		cmd1 := client.FTDictDump(ctx, "testdict3")
		Expect(cmd1.Err()).To(HaveOccurred()) // should not exist
		Expect(len(cmd1.Val())).To(BeZero())
		cmd2 := client.FTDictAdd(ctx, "testdict3", "foo", "bar", "baz")
		Expect(cmd2.Err()).NotTo(HaveOccurred())
		Expect(cmd2.Val()).To(BeNumerically("==", 3))
		cmd3 := client.FTDictDump(ctx, "testdict3")
		Expect(cmd3.Err()).NotTo(HaveOccurred())
		Expect(cmd3.Val()).To(ContainElements([]string{"foo", "bar", "baz"}))
	})
})
