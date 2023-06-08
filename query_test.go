package grstack_test

import (
	"math"

	grstack "github.com/goslogan/grstack"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// ft.search accounts "@tam:nic\.gibson"
var _ = Describe("Query basics", Label("hash", "query", "ft.search"), func() {

	It("doesn't raise an error on a valid query", func() {
		Expect(client.FTSearch(ctx, "hcustomers", `@email:ejowers0\@unblog\.fr`, nil).
			Err()).To(Not(HaveOccurred()))
	})
	It("can generate a valid query", func() {
		Expect(client.FTSearch(ctx, "hcustomers", `@email:ejowers0@unblog.fr`, nil).
			String()).To(ContainSubstring(`ft.search hcustomers @email:ejowers0@unblog.fr`))
	})

	It("can search for a specific result by attribute", func() {
		cmd := client.FTSearch(ctx, "hcustomers", `@email:ejowers0\@unblog\.fr`, nil)
		Expect(cmd.Err()).NotTo(HaveOccurred())
	})

	It("can return a single result", func() {
		Expect(client.FTSearch(ctx, "hcustomers", `@email:ejowers0\@unblog\.fr`, nil).Len()).To(Equal(1))
	})

	It("can return a map result", func() {
		Expect(client.FTSearch(ctx, "hcustomers", `@id:{1121175}`, nil).Val()).To(Equal(
			grstack.QueryResults{&grstack.QueryResult{
				Key:         "haccount:1121175",
				Score:       0,
				Explanation: nil,
				Values: &grstack.HashQueryValue{
					Value: map[string]string{
						"account_owner": "lara.croft",
						"balance":       "927.00",
						"customer":      "Kandace Korneichuk",
						"email":         `kkorneichukc\@cpanel\.net`,
						"ip":            `148\.140\.255\.235`,
						"account_id":    "1121175",
					}}}}))
	})

	It("will fail quietly with no search defined", func() {
		cmd := client.FTSearch(ctx, "hcustomers", "", nil)
		Expect(cmd.Err()).NotTo(HaveOccurred())
		Expect(cmd.Len()).To(Equal(0))
	})

	It("can return all the results for a given tag", func() {
		Expect(client.FTSearch(ctx, "hcustomers", `@owner:{lara\.croft}`, nil).Len()).To(Equal(10))
	})

})

var _ = Describe("Query options", Label("hash", "query", "ft.search"), func() {

	It("will return empty results - NOCONTENT", func() {
		opts := grstack.NewQueryOptions()
		opts.NoContent = true
		cmd := client.FTSearch(ctx, "hcustomers", `@id:{1121175}`, opts)
		Expect(cmd.Val().Key("haccount:1121175")).NotTo(BeNil())
		Expect(cmd.Val().Key("haccount:1121175")).To(BeAssignableToTypeOf(&grstack.QueryResult{}))
		Expect(cmd.Val().Key("haccount:1121175")).To(Equal(
			&grstack.QueryResult{
				Key:         "haccount:1121175",
				Score:       0,
				Explanation: nil,
				Values:      nil,
			}))
	})

	It("will return scores - WITHSCORES", func() {
		opts := grstack.NewQueryOptions()
		opts.WithScores = true
		cmd := client.FTSearch(ctx, "hcustomers", `@id:{1121175}`, opts)
		Expect(cmd.Err()).NotTo(HaveOccurred())
		Expect(cmd.Val()).To(BeAssignableToTypeOf(grstack.QueryResults{}))
		Expect(cmd.Val().Key("haccount:1121175").Score).Should(BeNumerically(">=", 1))
	})

	It("will return filtered results - FILTER (numeric)", func() {
		opts := grstack.NewQueryOptions()
		opts.NoContent = true
		opts.Filters = []grstack.QueryFilter{
			{
				Attribute: "balance",
				Min:       grstack.FilterValue(0, true),
				Max:       grstack.FilterValue(math.Inf(1), false),
			},
		}
		cmd := client.FTSearch(ctx, "hcustomers", `@owner:{lara\.croft}`, opts)
		Expect(cmd.Err()).NotTo(HaveOccurred())
		Expect(len(cmd.Val())).To(Equal(2))
	})

	It("can explain a score", func() {
		opts := grstack.NewQueryOptions()
		opts.WithScores = true
		opts.ExplainScore = true
		cmd := client.FTSearch(ctx, "hcustomers", `@owner:{lara\.croft}`, opts)
		Expect(cmd.Err()).NotTo(HaveOccurred())
		Expect(len(cmd.Val())).NotTo(BeZero())
		Expect(cmd.Val().Key("haccount:806396").Explanation).NotTo(BeNil())
	})

	It("can sort results", Label("sort"), func() {
		results := []string{}
		opts := grstack.NewQueryOptions()
		opts.NoContent = true
		opts.SortBy = "customer"
		cmd := client.FTSearch(ctx, "hcustomers", `@owner:{lara\.croft}`, opts)
		Expect(cmd.Err()).NotTo(HaveOccurred())
		for _, k := range cmd.Val() {
			results = append(results, k.Key)
		}
		Expect(results).To(ConsistOf([]string{"haccount:1339089", "haccount:239155", "haccount:575072", "haccount:765279", "haccount:1826581", "haccount:1371128", "haccount:1121175", "haccount:886088", "haccount:806396", "haccount:507187"}))

	})

	It("can limit to the first n results", Label("limit"), func() {
		results := []string{}
		opts := grstack.NewQueryOptions()
		opts.NoContent = true
		opts.Limit = &grstack.Limit{Num: 3, Offset: 0}
		opts.SortBy = "customer"
		cmd := client.FTSearch(ctx, "hcustomers", `@owner:{lara\.croft}`, opts)
		Expect(cmd.Err()).NotTo(HaveOccurred())
		Expect(cmd.Val()).To(HaveLen(3))
		for _, k := range cmd.Val() {
			results = append(results, k.Key)
		}
		Expect(results).To(ConsistOf([]string{"haccount:1339089", "haccount:239155", "haccount:575072"}))

	})

	It("can return n results from an offset", Label("limit"), func() {
		results := []string{}
		opts := grstack.NewQueryOptions()
		opts.NoContent = true
		opts.Limit = &grstack.Limit{Num: 3, Offset: 3}
		opts.SortBy = "customer"
		cmd := client.FTSearch(ctx, "hcustomers", `@owner:{lara\.croft}`, opts)
		Expect(cmd.Err()).NotTo(HaveOccurred())
		Expect(cmd.Val()).To(HaveLen(3))
		for _, k := range cmd.Val() {
			results = append(results, k.Key)
		}
		Expect(results).To(ConsistOf([]string{"haccount:765279", "haccount:1826581", "haccount:1371128"}))
	})

	It("can handle the RETURN subcommand", func() {
		opts := grstack.NewQueryOptions()
		opts.Return = append(opts.Return, grstack.QueryReturn{
			Name: "owner",
		})
		opts.Return = append(opts.Return, grstack.QueryReturn{Name: "balance"})
		cmd := client.FTSearch(ctx, "hcustomers", `@owner:{ellen\.ripley}`, opts)
		Expect(cmd.Err()).NotTo(HaveOccurred())
		Expect(cmd.Val()).To(Equal(grstack.QueryResults{
			&grstack.QueryResult{
				Key:         "haccount:536299",
				Score:       0,
				Explanation: nil,
				Values: &grstack.HashQueryValue{

					Value: map[string]string{"owner": "ellen.ripley", "balance": "113"},
				}}}))
	})

})

var _ = Describe("JSON query basics", Label("json", "query", "ft.search"), func() {

	It("doesn't raise an error on a valid query", func() {
		Expect(client.FTSearch(ctx, "jcustomers", `@email:ejowers0\@unblog\.fr`, nil).
			Err()).To(Not(HaveOccurred()))
	})
	It("can generate a valid query", func() {
		Expect(client.FTSearch(ctx, "jcustomers", `@email:ejowers0@unblog.fr`, nil).
			String()).To(ContainSubstring(`ft.search jcustomers @email:ejowers0@unblog.fr`))
	})

})

var _ = Describe("JSON searches", Label("json", "query", "ft.search"), func() {

	It("will return valid results", func() {
		opts := grstack.NewQueryOptions()
		cmd := client.FTSearchJSON(ctx, "jcustomers", `@id:{1121175}`, opts)
		Expect(cmd.Err()).NotTo(HaveOccurred())
		Expect(cmd.Val().Key("jaccount:1121175")).NotTo(BeNil())
		Expect(cmd.Val().Key("jaccount:1121175").Values).To(BeAssignableToTypeOf(&grstack.JSONQueryValue{}))
	})

	It("can return multiple search results", func() {

		cmd := client.FTSearchJSON(ctx, "jsoncomplex", `@datum:[1 1]`, nil)
		Expect(cmd.Err()).NotTo(HaveOccurred())
		Expect(cmd.Val().Key("jcomplex1")).NotTo(BeNil())
		Expect(cmd.Val().Key("jcomplex2")).NotTo(BeNil())
		Expect(cmd.Val().Key("jcomplex2").Values.(*grstack.JSONQueryValue).Value).To(Equal(map[string]interface{}{
			"$": map[string]interface{}{
				"data": float64(1),
				"test": map[string]interface{}{
					"data": float64(1),
				},
			}},
		))
		Expect(cmd.Val().Key("jcomplex1").Values.(*grstack.JSONQueryValue).Value).To(Equal(map[string]interface{}{
			"$": map[string]interface{}{

				"data": float64(1),
				"test1": map[string]interface{}{
					"data": float64(2),
				},
				"test2": map[string]interface{}{
					"data": float64(1),
				},
			}},
		))
	})

	It("can handle multiple search results with DIALECT 3", func() {

		options := grstack.NewQueryOptions()
		options.Dialect = 3
		cmd := client.FTSearchJSON(ctx, "jsoncomplex", `@datum:[1 1]`, options)
		Expect(cmd.Err()).NotTo(HaveOccurred())
		Expect(cmd.Val().Key("jcomplex1")).NotTo(BeNil())
		Expect(cmd.Val().Key("jcomplex2")).NotTo(BeNil())
		Expect(cmd.Val().Key("jcomplex2").Values.(*grstack.JSONQueryValue).Value).To(Equal(map[string]interface{}{
			"$": []interface{}{map[string]interface{}{
				"data": float64(1),
				"test": map[string]interface{}{
					"data": float64(1),
				},
			}},
		}))
		Expect(cmd.Val().Key("jcomplex1").Values.(*grstack.JSONQueryValue).Value).To(Equal(map[string]interface{}{
			"$": []interface{}{map[string]interface{}{
				"data": float64(1),
				"test1": map[string]interface{}{
					"data": float64(2),
				},
				"test2": map[string]interface{}{
					"data": float64(1),
				},
			}},
		}))
	})

	It("can return values when we use RETURN", func() {

		options := grstack.NewQueryOptions()
		options.Return = append(options.Return, grstack.QueryReturn{
			Name: "$..data",
			As:   "answer",
		})
		cmd := client.FTSearchJSON(ctx, "jsoncomplex", `@datum:[1 1]`, options)
		Expect(cmd.Err()).NotTo(HaveOccurred())
		Expect(cmd.Val().Key("jcomplex1")).NotTo(BeNil())
		Expect(cmd.Val().Key("jcomplex2")).NotTo(BeNil())
		Expect(cmd.Val().Key("jcomplex2").Values.(*grstack.JSONQueryValue).Value).To(Equal(map[string]interface{}{"answer": float64(1)}))
		Expect(cmd.Val().Key("jcomplex1").Values.(*grstack.JSONQueryValue).Value).To(Equal(map[string]interface{}{"answer": float64(1)}))
	})

	It("can return values when we use RETURN and DIALECT v3", func() {

		options := grstack.NewQueryOptions()
		options.Return = append(options.Return, grstack.QueryReturn{
			Name: "$..data",
			As:   "answer",
		})
		options.Dialect = 3
		cmd := client.FTSearchJSON(ctx, "jsoncomplex", `@datum:[1 1]`, options)
		Expect(cmd.Err()).NotTo(HaveOccurred())
		Expect(cmd.Val().Key("jcomplex1")).NotTo(BeNil())
		Expect(cmd.Val().Key("jcomplex2")).NotTo(BeNil())
		Expect(cmd.Val().Key("jcomplex1").Values.(*grstack.JSONQueryValue).Value).To(Equal(map[string]interface{}{"answer": []interface{}{float64(1), float64(2), float64(1)}}))
		Expect(cmd.Val().Key("jcomplex2").Values.(*grstack.JSONQueryValue).Value).To(Equal(map[string]interface{}{"answer": []interface{}{float64(1), float64(1)}}))
	})

	It("can scan a result", func() {

		type Customer struct {
			Customer     string  `json:"customer"`
			Email        string  `json:"email"`
			IP           string  `json:"ip"`
			AccountId    string  `json:"account_id"`
			AccountOwner string  `json:"account_owner"`
			Balance      float64 `json:"balance"`
		}

		var customer = Customer{}

		cmd := client.FTSearchJSON(ctx, "jcustomers", `@owner:{ellen\.ripley}`, grstack.NewQueryOptions())
		Expect(cmd.Err()).NotTo(HaveOccurred())
		Expect(cmd.Val().Key("jaccount:536299")).NotTo(BeNil())
		Expect(cmd.Val().Key("jaccount:536299").Values.(*grstack.JSONQueryValue).Scan("$", &customer)).NotTo(HaveOccurred())
		Expect(customer.Customer).To(Equal("Davon Audley"))
		Expect(customer.Email).To(Equal(`daudley3\@alexa\.com`))
	})

	It("can return all the keys from a query result", func() {
		cmd := client.FTSearchJSON(ctx, "jsoncomplex", `@datum:[1 1]`, nil)
		Expect(cmd.Err()).NotTo(HaveOccurred())
		Expect(cmd.Val().Keys()).To(Equal([]string{"jcomplex2", "jcomplex1"}))
	})

})
