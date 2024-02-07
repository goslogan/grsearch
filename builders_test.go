package grsearch_test

import (
	"time"

	"github.com/goslogan/grsearch"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("We can build query options", Label("builders", "ft.search"), func() {

	It("can construct default query options", func() {
		base := grsearch.NewQueryOptions()
		built := grsearch.NewQueryBuilder()

		Expect(base).To(Equal(built.Options()))
	})

	It("can construct options with simple flags set", func() {
		base := grsearch.NewQueryOptions()
		base.Dialect = 2
		base.ExplainScore = true
		base.NoContent = true
		base.Timeout = time.Duration(1000)
		base.NoStopWords = true
		base.Verbatim = true

		built := grsearch.NewQueryBuilder().
			Dialect(2).
			ExplainScore().
			NoContent().
			Timeout(time.Duration(1000)).
			Verbatim().
			NoStopWords()

		Expect(base).To(Equal(built.Options()))

	})

	It("can construct queries with parameters", func() {
		base := grsearch.NewQueryOptions()
		base.Params = map[string]interface{}{
			"foo": "one",
			"bar": 2,
		}

		built := grsearch.NewQueryBuilder().
			Param("foo", "one").
			Param("bar", 2)

		Expect(base).To(Equal(built.Options()))
	})

	It("can construct queries with geofilters", func() {
		base := grsearch.NewQueryOptions()
		base.GeoFilters = []grsearch.GeoFilter{
			{
				Attribute: "test",
				Long:      100,
				Lat:       200,
				Radius:    300,
				Units:     "m",
			},
		}
		built := grsearch.NewQueryBuilder().
			GeoFilter("test", 100, 200, 300, "m")

		Expect(base).To(Equal(built.Options()))
	})

	It("can construct queries with filters", func() {
		base := grsearch.NewQueryOptions()
		base.Filters = []grsearch.QueryFilter{
			{
				Attribute: "test",
				Min:       -100,
				Max:       "+inf",
			},
		}
		built := grsearch.NewQueryBuilder().
			Filter("test", -100, "+inf")

		Expect(base).To(Equal(built.Options()))
	})

})

var _ = Describe("We can build aggregate options", Label("builders", "ft.aggregate"), func() {

	It("can construct default query options", func() {
		base := grsearch.NewAggregateOptions()
		built := grsearch.NewAggregateBuilder()

		Expect(base).To(Equal(built.Options()))
	})

	It("can construct options with simple flags set", func() {
		base := grsearch.NewAggregateOptions()
		base.Dialect = 2
		base.Steps = append(base.Steps, grsearch.AggregateFilter("@test != 3"))
		base.Timeout = time.Duration(1000)
		base.Verbatim = true

		built := grsearch.NewAggregateBuilder().
			Dialect(2).
			Timeout(time.Duration(1000)).
			Verbatim().
			Filter("@test != 3")

		Expect(base).To(Equal(built.Options()))

	})

	It("can construct queries with parameters", func() {
		base := grsearch.NewAggregateOptions()
		base.Params = map[string]interface{}{
			"foo": "one",
			"bar": 2,
		}

		built := grsearch.NewAggregateBuilder().
			Param("foo", "one").
			Param("bar", 2)

		Expect(base).To(Equal(built.Options()))
	})

	It("can construct a single group by", func() {
		base := grsearch.NewAggregateOptions()
		base.Steps = append(base.Steps, &grsearch.AggregateGroupBy{
			Properties: []string{"@name"},
			Reducers: []grsearch.AggregateReducer{
				{
					Name: "count",
					As:   "nameCount",
				},
			},
		})

		built := grsearch.NewAggregateBuilder().
			GroupBy(grsearch.NewGroupByBuilder().
				Properties([]string{"@name"}).
				Reduce(grsearch.ReduceCount("nameCount")).
				GroupBy())

		Expect(base).To(Equal(built.Options()))
	})

	It("can build a complex aggregate", func() {
		base := grsearch.NewAggregateOptions()
		base.Steps = append(base.Steps, &grsearch.AggregateApply{
			Expression: "@timestamp - (@timestamp % 86400)",
			As:         "day",
		})
		base.Steps = append(base.Steps, &grsearch.AggregateGroupBy{
			Properties: []string{"@day", "@country"},
			Reducers: []grsearch.AggregateReducer{{
				Name: "count",
				As:   "num_visits",
			}}})

		base.Steps = append(base.Steps, &grsearch.AggregateSort{
			Keys: []grsearch.AggregateSortKey{{
				Name:  "@day",
				Order: grsearch.SortAsc,
			}, {
				Name:  "@country",
				Order: grsearch.SortDesc,
			}},
		})

		built := grsearch.NewAggregateBuilder().
			Apply("@timestamp - (@timestamp % 86400)", "day").
			GroupBy(grsearch.NewGroupByBuilder().
				Properties([]string{"@day", "@country"}).
				Reduce(grsearch.ReduceCount("num_visits")).
				GroupBy()).
			SortBy([]grsearch.AggregateSortKey{{Name: "@day", Order: grsearch.SortAsc}, {Name: "@country", Order: grsearch.SortDesc}})

		Expect(base).To(Equal(built.Options()))

	})

})
