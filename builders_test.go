package grstack_test

import (
	"time"

	"github.com/goslogan/grstack"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("We can build query options", Label("builders", "ft.search"), func() {

	It("can construct default query options", func() {
		base := grstack.NewQueryOptions()
		built := grstack.NewQueryBuilder()

		Expect(base).To(Equal(built.Options()))
	})

	It("can construct options with simple flags set", func() {
		base := grstack.NewQueryOptions()
		base.Dialect = 2
		base.ExplainScore = true
		base.NoContent = true
		base.Timeout = time.Duration(1000)
		base.NoStopWords = true
		base.Verbatim = true

		built := grstack.NewQueryBuilder().
			Dialect(2).
			ExplainScore().
			NoContent().
			Timeout(time.Duration(1000)).
			Verbatim().
			NoStopWords()

		Expect(base).To(Equal(built.Options()))

	})

	It("can construct queries with parameters", func() {
		base := grstack.NewQueryOptions()
		base.Params = map[string]interface{}{
			"foo": "one",
			"bar": 2,
		}

		built := grstack.NewQueryBuilder().
			Param("foo", "one").
			Param("bar", 2)

		Expect(base).To(Equal(built.Options()))
	})

	It("can construct queries with geofilters", func() {
		base := grstack.NewQueryOptions()
		base.GeoFilters = []grstack.GeoFilter{
			{
				Attribute: "test",
				Long:      100,
				Lat:       200,
				Radius:    300,
				Units:     "m",
			},
		}
		built := grstack.NewQueryBuilder().
			GeoFilter("test", 100, 200, 300, "m")

		Expect(base).To(Equal(built.Options()))
	})

	It("can construct queries with filters", func() {
		base := grstack.NewQueryOptions()
		base.Filters = []grstack.QueryFilter{
			{
				Attribute: "test",
				Min:       -100,
				Max:       "+inf",
			},
		}
		built := grstack.NewQueryBuilder().
			Filter("test", -100, "+inf")

		Expect(base).To(Equal(built.Options()))
	})

})

var _ = Describe("We can build aggregate options", Label("builders", "ft.aggregate"), func() {

	It("can construct default query options", func() {
		base := grstack.NewAggregateOptions()
		built := grstack.NewAggregateBuilder()

		Expect(base).To(Equal(built.Options()))
	})

	It("can construct options with simple flags set", func() {
		base := grstack.NewAggregateOptions()
		base.Dialect = 2
		base.Filter = "@test != 3"
		base.Timeout = time.Duration(1000)
		base.Verbatim = true

		built := grstack.NewAggregateBuilder().
			Dialect(2).
			Timeout(time.Duration(1000)).
			Verbatim().
			Filter("@test != 3")

		Expect(base).To(Equal(built.Options()))

	})

	It("can construct queries with parameters", func() {
		base := grstack.NewAggregateOptions()
		base.Params = map[string]interface{}{
			"foo": "one",
			"bar": 2,
		}

		built := grstack.NewAggregateBuilder().
			Param("foo", "one").
			Param("bar", 2)

		Expect(base).To(Equal(built.Options()))
	})

	It("can construct a single group by", func() {
		base := grstack.NewAggregateOptions()
		base.GroupBy = []grstack.AggregateGroupBy{{
			Properties: []string{"@name"},
			Reducers: []grstack.AggregateReducer{
				{
					Name: "count",
					As:   "nameCount",
				},
			},
		}}

		built := grstack.NewAggregateBuilder().
			GroupBy(grstack.NewGroupByBuilder().
				Properties([]string{"@name"}).
				Reduce(grstack.ReduceCount("nameCount")).
				GroupBy())

		Expect(base).To(Equal(built.Options()))
	})

	It("can build a complex aggregate", func() {
		base := grstack.NewAggregateOptions()
		base.Apply = []grstack.AggregateApply{{
			Expression: "@timestamp - (@timestamp % 86400)",
			As:         "day",
		}}
		base.GroupBy = []grstack.AggregateGroupBy{{
			Properties: []string{"@day", "@country"},
			Reducers: []grstack.AggregateReducer{{
				Name: "count",
				As:   "num_visits",
			}},
		}}
		base.SortBy = &grstack.AggregateSort{
			Keys: []grstack.AggregateSortKey{{
				Name:  "@day",
				Order: grstack.SortAsc,
			}, {
				Name:  "@country",
				Order: grstack.SortDesc,
			}},
		}

		built := grstack.NewAggregateBuilder().
			Apply("@timestamp - (@timestamp % 86400)", "day").
			GroupBy(grstack.NewGroupByBuilder().
				Properties([]string{"@day", "@country"}).
				Reduce(grstack.ReduceCount("num_visits")).
				GroupBy()).
			SortBy("@day", grstack.SortAsc).
			SortBy("@country", grstack.SortDesc)

		Expect(base).To(Equal(built.Options()))

	})

})
