package ftsearch_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/RedisLabs-Solution-Architects/go-search/ftsearch"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/redis/go-redis/v9"
)

var client *ftsearch.Client
var ctx = context.Background()

var testData = [][]string{
	{"Davon", "Audley", `daudley3\@alexa.com`, `174\.113\.231\.230`, "536299", "daniel.preiskel", "113.00"},
	{"Emmet", "Jowers", `ejowers0\@unblog\.fr`, `167\.230\.3\.244`, "467734", "igor.jovanovic", "-999.99"},
	{"Essy", "Kiddle", `ekiddle1\@si\.edu`, `230\.144\.170\.69`, "1044493", "igor.jovanovic", "0"},
	{"Vlad", "Darrigrand", `vdarrigrand2\@kickstarter\.com`, `45\.89\.83\.231`, "508709", "igor.jovanovic", "0"},
	{"Remington", "Ponceford", `rponceford4\@topsy\.com`, `236\.175\.61\.40`, "1505510", "igor.jovanovic", "0"},
	{"Aili", "Brahms", `abrahms5\@wikia\.com`, `182\.115\.99\.238`, "1339089", `nic.gibson`, "0"},
	{"Lynn", "Beed", `lbeed6\@jalbum\.net`, `142\.218\.3\.176`, "886088", "nic.gibson", "0"},
	{"Corabelle", "Bertelmot", `cbertelmot7\@amazon\.com`, `122\.247\.52\.99`, "1222097", "igor.jovanovic", "-3402.00"},
	{"Demetri", "Vigors", `dvigors8\@w3\.org`, `40\.75\.113\.150`, "1952347", "igor.jovanovic", "1371.00"},
	{"Rafaelita", "Wisam", `rwisam9\@vk\.com`, `175\.40\.66\.248`, "507187", "nic.gibson", "0"},
	{"Jsandye", "Sprackling", `jspracklinga\@ow\.ly`, `177\.67\.221\.138`, "419113", "igor.jovanovic", "0"},
	{"Chen", "Clilverd", `cclilverdb\@stanford\.edu`, `164\.230\.108\.100`, "765279", "nic.gibson", "0"},
	{"Kandace", "Korneichuk", `kkorneichukc\@cpanel\.net`, `148\.140\.255\.235`, "1121175", "nic.gibson", "927.00"},
	{"Brandy", "Gustus", `bgustusd\@loc\.gov`, `233\.98\.10\.248`, "575072", "nic.gibson", "3.24"},
	{"Annabal", "O'Carran", `aocarrane\@instagram\.com`, `234\.240\.12\.81`, "1888382", "igor.jovanovic", "-560.88"},
	{"Gizela", "Rolph", `grolphf\@theguardian\.com`, `91\.16\.124\.34`, "1371128", "nic.gibson", "0"},
	{"Curtice", "Iscowitz", `ciscowitzg\@newyorker\.com`, `118\.83\.100\.5`, "1826581", "nic.gibson", "0"},
	{"Meggy", "Sheward", `mshewardh\@alexa\.com`, `56\.21\.83\.123`, "806396", "nic.gibson", "0"},
	{"Binnie", "Sowerby", `bsowerbyi\@china\.com.cn`, `185\.1\.56\.15`, "239155", "nic.gibson", "0"},
	{"Tamarah", "Hallybone", `thallybonej\@wisc\.edu`, `115\.140\.35\.151`, "376460", "igor.jovanovic", "-6503.33"},
	{"Niko", "Drillingcourt", `ndrillingcourtk\@nydailynews\.com`, `37\.18\.13\.16`, "1443633", "igor.jovanovic", "-1200.23"},
	{"Martynne", "Shovell", `mshovelll\@google\.com.br`, `73\.142\.109\.212`, "1840413", "igor.jovanovic", "0"},
	{"Trev", "Todman", `ttodmanm\@rambler\.ru`, `137\.130\.15\.201`, "232226", "igor.jovanovic", "-991.31"},
	{"Fields", "Baldry", `fbaldryn\@weibo\.com`, `58\.157\.227\.177`, "355348", "igor.jovanovic", "0"},
	{"Tracy", "Pauly", `tpaulyo\@myspace\.com`, `75\.189\.188\.225`, "1838082", "igor.jovanovic", "650.01"},
}

func TestFtsearch(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Ftsearch Suite")
}

// initialise the test data we use throughout
func createTestData() {
	Expect(client.FlushAll(ctx).Err()).NotTo(HaveOccurred())

	for _, row := range testData {
		Expect(client.HSet(ctx, fmt.Sprintf("account:%s", row[4]),
			"customer", row[0]+" "+row[1],
			"email", row[2],
			"ip", row[3],
			"account_id", row[4],
			"account_owner", row[5],
			"balance", row[6],
		).Err()).NotTo(HaveOccurred())
	}

	Expect(client.FTCreateIndex(ctx, "customers", ftsearch.NewIndexOptions().
		AddPrefix("account:").
		AddSchemaAttribute(ftsearch.TagAttribute{
			Name:     "account_id",
			Alias:    "id",
			Sortable: true}).AddSchemaAttribute(ftsearch.TextAttribute{
		Name:     "customer",
		Sortable: true}).AddSchemaAttribute(ftsearch.TextAttribute{
		Name:     "email",
		Sortable: true}).AddSchemaAttribute(ftsearch.TagAttribute{
		Name:     "account_owner",
		Alias:    "owner",
		Sortable: true,
	})).Err()).NotTo(HaveOccurred())
}

var _ = BeforeSuite(func() {
	client = ftsearch.NewClient(&redis.Options{})
	Expect(client.Ping(ctx).Err()).NotTo(HaveOccurred())
	Expect(client.FlushAll(ctx).Err()).NotTo(HaveOccurred())
	createTestData()
	time.Sleep(time.Second * 5)
})

/* var _ = AfterSuite(func() {
	Expect(client.FlushDB(ctx).Err()).NotTo(HaveOccurred())
}) */
