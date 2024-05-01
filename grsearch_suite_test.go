package grsearch_test

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"regexp"
	"strings"
	"testing"
	"text/template"
	"time"

	_ "embed"

	grsearch "github.com/goslogan/grsearch"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/redis/go-redis/v9"
)

//go:embed testdata/customers_test.csv
var customerData string

//go:embed testdata/customers_test.template
var customerJSON string

var client *grsearch.Client
var ctx = context.Background()

// convert strings that need to stay "as one for tokenising"
func escapeForHash(s string) string {
	re := regexp.MustCompile(`([,.<>{}\[\]"':;!@#$%^&*()\-+=~ ])`)
	return re.ReplaceAllString(s, `\$1`)
}

// convert strings that need to stay "as one for tokenising"
func escapeForJSON(s string) string {
	re := regexp.MustCompile(`([,.<>{}\[\]"':;!@#$%^&*()\-+=~ ])`)
	return re.ReplaceAllString(s, "\\\\$1")
}

// initialise the test data we use throughout
func createJSONTestData() {

	fmt.Println("Generating JSON Data...")

	csvData := strings.NewReader(customerData)
	csvReader := csv.NewReader(csvData)
	t, err := template.New("customer").Parse(customerJSON)
	Expect(err).NotTo((HaveOccurred()))

	// skip the first one so we can have something that is found in one index and
	// not the other
	first := true
	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if first {
			first = false
			continue
		}
		row[2] = escapeForJSON(row[2])
		row[3] = escapeForJSON(row[3])
		Expect(err).NotTo(HaveOccurred())
		var js bytes.Buffer
		Expect(t.ExecuteTemplate(&js, "customer", row)).NotTo(HaveOccurred())
		Expect(client.JSONSet(ctx, fmt.Sprintf("jaccount:%s", row[4]), "$", js.String()).Err()).NotTo(HaveOccurred())
	}

}

// initialise the test data we use throughout
func createHashTestData() {

	fmt.Println("Generating Hash Data...")

	csvData := strings.NewReader(customerData)
	csvReader := csv.NewReader(csvData)

	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		Expect(err).NotTo(HaveOccurred())

		Expect(client.HSet(ctx, fmt.Sprintf("haccount:%s", row[4]),
			"customer", row[0]+" "+row[1],
			"email", escapeForHash(row[2]),
			"ip", escapeForHash(row[3]),
			"account_id", row[4],
			"account_owner", row[5],
			"balance", row[6],
			"country", row[7],
		).Err()).NotTo(HaveOccurred())

	}

}

func createHashIndexes() {
	fmt.Println("Generating Hash Indexes...")
	Expect(client.FTCreate(ctx, "hcustomers", grsearch.NewIndexBuilder().
		Prefix("haccount:").
		Schema(&grsearch.TagAttribute{
			Name:     "account_id",
			Alias:    "id",
			Sortable: true}).
		Schema(&grsearch.TextAttribute{
			Name:     "customer",
			Sortable: true}).
		Schema(&grsearch.TextAttribute{
			Name:     "email",
			Sortable: false}).
		Schema(&grsearch.TagAttribute{
			Name:     "account_owner",
			Alias:    "owner",
			Sortable: true}).
		Schema(&grsearch.NumericAttribute{
			Name:     "balance",
			Sortable: true}).
		Schema(&grsearch.TagAttribute{
			Name:     "country",
			Sortable: true,
		}).
		Options()).Err()).NotTo(HaveOccurred())
}

func createJSONIndexes() {

	fmt.Println("Generating JSON Indexes...")
	cmd := client.FTCreate(ctx, "jcustomers", grsearch.NewIndexBuilder().
		On("json").
		Prefix("jaccount:").
		Schema(&grsearch.TagAttribute{
			Name:     "$.account_id",
			Alias:    "id",
			Sortable: true}).
		Schema(&grsearch.TextAttribute{
			Name:     "$.customer",
			Alias:    "customer",
			Sortable: true}).
		Schema(&grsearch.TextAttribute{
			Name:     "$.email",
			Alias:    "email",
			Sortable: false}).
		Schema(&grsearch.TagAttribute{
			Name:     "$.account_owner",
			Alias:    "owner",
			Sortable: true}).
		Schema(&grsearch.NumericAttribute{
			Name:     "$.balance",
			Alias:    "balance",
			Sortable: true,
		}).
		Schema(&grsearch.TagAttribute{
			Name:     "$.country",
			Alias:    "country",
			Sortable: true,
		}).
		Options())
	Expect(cmd.Err()).NotTo(HaveOccurred())

	// complex documents for testing JSON edge cases
	doc1 := `{"data": 1, "test1": {"data": 2 }, "test2": {"data": 1}}`
	doc2 := `{"data": 1, "test": {"data": 1 }}`

	Expect(client.JSONSet(ctx, "jcomplex1", "$", doc1).Err()).NotTo(HaveOccurred())
	Expect(client.JSONSet(ctx, "jcomplex2", "$", doc2).Err()).NotTo(HaveOccurred())
	Expect(client.FTCreate(ctx, "jsoncomplex",
		grsearch.NewIndexBuilder().
			On("json").
			Prefix("jcomplex").
			Schema(&grsearch.NumericAttribute{
				Name:     "$..data",
				Alias:    "datum",
				Sortable: true,
			}).Options()).Err()).NotTo(HaveOccurred())

}

func TestFtsearch(t *testing.T) {
	RegisterFailHandler(Fail)
	suiteConfig, reportConfig := GinkgoConfiguration()
	//suiteConfig.LabelFilter = "json && FT.SEARCH && valid"
	RunSpecs(t, "Ftsearch Suite", suiteConfig, reportConfig)
}

var _ = BeforeSuite(func() {
	client = grsearch.NewClient(&redis.Options{})
	Expect(client.Ping(ctx).Err()).NotTo(HaveOccurred())
	Expect(client.FlushAll(ctx).Err()).NotTo(HaveOccurred())
	createHashTestData()
	createHashIndexes()
	createJSONTestData()
	createJSONIndexes()
	time.Sleep(time.Second * 5)
})

/* var _ = AfterSuite(func() {
	Expect(client.FlushDB(ctx).Err()).NotTo(HaveOccurred())
}) */
