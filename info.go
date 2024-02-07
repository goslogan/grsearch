package grsearch

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/goslogan/grsearch/internal"
)

/****
* data structures for the FT.info command
*****/

// Info represents the parsed results of a call to FT.INFO.
type Info struct {
	IndexName                  string
	Index                      *IndexOptions
	NumDocs                    int64
	MaxDocId                   int64
	NumTerms                   int64
	NumRecords                 int64
	Indexing                   float64
	PercentIndexed             float64
	HashIndexingFailures       int64
	TotalInvertedIndexBlocks   int64
	InvertedSize               float64
	VectorIndexSize            float64
	DocTableSize               float64
	OffsetVectorsSize          float64
	SortableValuesSize         float64
	KeyTableSize               float64
	AverageRecordsPerDoc       float64
	AverageBytesPerRecord      float64
	AverageOffsetsPerTerm      float64
	AverageOffsetBitsPerRecord float64
	TotalIndexingTime          time.Duration
	NumberOfUses               int64
	GCStats                    struct {
		BytesCollected       int64
		TotalMsRun           time.Duration
		TotalCycles          int64
		AverageCycleTime     time.Duration
		LastRunTime          time.Duration
		GCNumericTreesMissed int64
		GCBlocksDenied       int64
	}
	CursorStats struct {
		GlobalIdle    int64
		GlobalTotal   int64
		IndexCapacity int64
		IndexTotal    int64
	}
	DialectStats struct {
		Dialect1 int64
		Dialect2 int64
		Dialect3 int64
	}
}

// parse takes the results of an FT.INFO command and creates and Info
// struct from it.
func (info *Info) parse(result map[interface{}]interface{}) error {

	info.IndexName = result["index_name"].(string)
	info.NumDocs, _ = internal.Int64(result["num_docs"])
	info.MaxDocId, _ = internal.Int64(result["max_doc_id"])
	info.NumTerms, _ = internal.Int64(result["num_terms"])
	info.NumRecords, _ = internal.Int64(result["num_records"])
	info.Indexing, _ = internal.Float64(result["indexing"])
	info.PercentIndexed, _ = internal.Float64(result["percent_indexed"])
	info.HashIndexingFailures, _ = internal.Int64(result["hash_indexing_failures"])
	info.TotalInvertedIndexBlocks, _ = internal.Int64(result["total_inverted_index_blocks"])
	info.InvertedSize, _ = internal.Float64(result["inverted_sz_mb"])
	info.VectorIndexSize, _ = internal.Float64(result["vector_index_sz_mb"])
	info.DocTableSize, _ = internal.Float64(result["doc_table_size_mb"])
	info.OffsetVectorsSize, _ = internal.Float64(result["offset_vectors_sz_mb"])
	info.SortableValuesSize, _ = internal.Float64(result["sortable_values_size_mb"])
	info.KeyTableSize, _ = internal.Float64(result["key_table_size_mb"])
	info.AverageRecordsPerDoc, _ = internal.Float64(result["records_per_doc_avg"])
	info.AverageBytesPerRecord, _ = internal.Float64(result["bytes_per_record_avg"])
	info.AverageOffsetsPerTerm, _ = internal.Float64(result["offsets_per_term_avg"])
	info.AverageOffsetBitsPerRecord, _ = internal.Float64(result["offset_bits_per_record_avg"])
	info.NumberOfUses, _ = internal.Int64(result["number_of_uses"])

	// Given no other evidence, we assume this is seconds
	t, _ := internal.Float64(result["total_indexing_time"])
	info.TotalIndexingTime, _ = time.ParseDuration(fmt.Sprintf("%fs", t))

	// Parse the index stats.
	info.parseIndexStats(result)

	// Parse out the options.
	return info.parseIndexOptionsFromInfo(result)

}

// Parse and store index stats
func (info *Info) parseIndexStats(result map[interface{}]interface{}) {

	gc := internal.ToMap(result["gc_stats"])
	v, _ := internal.Int64(gc["average_cycle_time_ms"])
	info.GCStats.AverageCycleTime = time.Duration(v) * time.Millisecond
	v, _ = internal.Int64(gc["total_ms_run"])
	info.GCStats.TotalMsRun = time.Duration(v) * time.Millisecond
	v, _ = internal.Int64(gc["last_run_time_ms"])
	info.GCStats.LastRunTime = time.Duration(v) * time.Millisecond
	info.GCStats.BytesCollected, _ = internal.Int64(gc["bytes_collected"])
	info.GCStats.TotalCycles, _ = internal.Int64(gc["total_cycles"])
	info.GCStats.GCBlocksDenied, _ = internal.Int64(gc["gc_blocks_denied"])
	info.GCStats.GCNumericTreesMissed, _ = internal.Int64(gc["gc_numeric_trees_missed"])

	c := internal.ToMap(result["cursor_stats"])
	info.CursorStats.GlobalIdle, _ = internal.Int64(c["global_idle"])
	info.CursorStats.GlobalTotal, _ = internal.Int64(c["global_total"])
	info.CursorStats.IndexCapacity, _ = internal.Int64(c["index_capacity"])
	info.CursorStats.IndexTotal, _ = internal.Int64(c["index_total"])

	d := internal.ToMap(result["dialect_stats"])
	info.DialectStats.Dialect1, _ = internal.Int64(d["dialect_1"])
	info.DialectStats.Dialect2, _ = internal.Int64(d["dialect_2"])
	info.DialectStats.Dialect3, _ = internal.Int64(d["dialect_3"])
}

// Create IndexOptions from ft.info output
func (info *Info) parseIndexOptionsFromInfo(input map[interface{}]interface{}) error {

	i := &IndexOptions{}

	if data, ok := input["index_definition"]; ok {
		mapped := internal.ToMap(data)
		if mapped["key_type"] == "JSON" {
			i.On = "JSON"
		} else {
			i.On = "HASH"
		}
		i.Score, _ = internal.Float64(mapped["default_score"])
		i.Prefix = make([]string, len(mapped["prefixes"].([]interface{})))
		for n, p := range mapped["prefixes"].([]interface{}) {
			i.Prefix[n] = p.(string)
		}
	}

	if i.Schema == nil {
		i.Schema = make([]SchemaAttribute, 0)
	}
	if data, ok := input["attributes"]; ok {
		for _, a := range data.([]interface{}) {
			attribInfo := info.attribInfoMap(a)
			// deal with flags and RESP2

			attribType := attribInfo["type"].(string)

			var attribute SchemaAttribute
			switch strings.ToLower(attribType) {
			case "tag":
				attribute = &TagAttribute{}
			case "text":
				attribute = &TextAttribute{}
			case "numeric":
				attribute = &NumericAttribute{}
			case "geo":
				attribute = &GeoAttribute{}
			case "geometry":
				attribute = &GeometryAttribute{}
			case "vector":
				attribute = &VectorAttribute{}
			default:
				return fmt.Errorf("grsearch: unhandled attribute type: %s", attribInfo[5].(string))
			}
			if attribInfo != nil {
				attribute.parseFromInfo(attribInfo)
				i.Schema = append(i.Schema, attribute)
			}
		}

	}
	info.Index = i
	return nil
}

// attribInfoMap converts the attribute definition from FT.INFO into a
// a map, handling the use of SORTABLE and UNF in RESP2 and RESP3 properly
func (info *Info) attribInfoMap(result interface{}) map[interface{}]interface{} {

	var attrib map[interface{}]interface{}

	switch val := result.(type) {
	case []interface{}: // RESP2
		attrib = internal.ToMap(result)
		for _, item := range val {
			if s, ok := item.(string); ok {
				if strings.ToLower(s) == "sortable" {
					attrib["sortable"] = true
				}
				if strings.ToLower(s) == "unf" {
					attrib["unf"] = true
				}
			}
		}
	case map[interface{}]interface{}:
		attrib = val
		if flags, ok := val["flags"]; ok {
			for _, a := range flags.([]interface{}) {
				if s, ok := a.(string); ok {
					if strings.ToLower(s) == "sortable" {
						attrib["sortable"] = true
					}
					if strings.ToLower(s) == "unf" {
						attrib["unf"] = true
					}
				}
			}
		}
	}

	return attrib
}

func (a *TagAttribute) parseFromInfo(source map[interface{}]interface{}) {

	for key, val := range source {
		switch strings.ToLower(key.(string)) {
		case "identifier":
			a.Name = val.(string)
		case "attribute":
			a.Alias = val.(string)
		case "separator":
			a.Separator = val.(string)
		case "sortable":
			a.Sortable = true
		case "unf":
			a.UnNormalized = true
		case "casesensitive":
			a.CaseSensitive = true
		case "withsuffixtrie":
			a.WithSuffixTrie = true
		case "noindex":
			a.NoIndex = true
		}
	}

}

func (a *TextAttribute) parseFromInfo(source map[interface{}]interface{}) {
	for key, val := range source {

		switch strings.ToLower(key.(string)) {
		case "identifier":
			a.Name = val.(string)
		case "attribute":
			a.Alias = val.(string)
		case "sortable":
			a.Sortable = true
		case "unf":
			a.UnNormalized = true
		case "withsuffixtrie":
			a.WithSuffixTrie = true
		case "noindex":
			a.NoIndex = true
		case "weight":
			a.Weight, _ = internal.Float64(val)
		case "phonetic":
			a.Phonetic = val.(string)
		case "nostem":
			a.NoStem = true
		}
	}
}

func (a *NumericAttribute) parseFromInfo(source map[interface{}]interface{}) {
	for key, val := range source {

		switch strings.ToLower(key.(string)) {
		case "identifier":
			a.Name = val.(string)
		case "attribute":
			a.Alias = val.(string)
		case "sortable":
			a.Sortable = true
		case "noindex":
			a.NoIndex = true
		}
	}
}

func (a *GeometryAttribute) parseFromInfo(source map[interface{}]interface{}) {
	for key, val := range source {

		switch strings.ToLower(key.(string)) {
		case "identifier":
			a.Name = val.(string)
		case "attribute":
			a.Alias = val.(string)
		}
	}
}

func (a *GeoAttribute) parseFromInfo(source map[interface{}]interface{}) {
	for key, val := range source {

		switch strings.ToLower(key.(string)) {
		case "identifier":
			a.Name = val.(string)
		case "attribute":
			a.Alias = val.(string)
		case "sortable":
			a.Sortable = true
		case "noindex":
			a.NoIndex = true
		}
	}
}

func (a *VectorAttribute) parseFromInfo(source map[interface{}]interface{}) {
	for key, val := range source {

		switch strings.ToLower(key.(string)) {
		case "flat", "hnsw":
			a.Algorithm = key.(string)
		case "type":
			a.Type = strings.ToLower(val.(string))
		case "dim":
			i, _ := strconv.ParseUint(val.(string), 10, 64)
			a.Dim = i
		case "distance_metric":
			a.DistanceMetric = strings.ToLower(val.(string))
		case "initial_cap":
			cap, _ := strconv.ParseUint(val.(string), 10, 64)
			a.InitialCap = cap
		case "block_size":
			size, _ := strconv.ParseUint(val.(string), 10, 64)
			a.BlockSize = size
		case "m":
			m, _ := strconv.ParseUint(val.(string), 10, 64)
			a.M = m
		case "ef_construction":
			ef, _ := strconv.ParseUint(val.(string), 10, 64)
			a.EFConstruction = ef
		case "ef_runtime":
			ef, _ := strconv.ParseUint(val.(string), 10, 64)
			a.EFRuntime = ef
		case "epsilon":
			ef, _ := strconv.ParseFloat(val.(string), 64)
			a.Epsilon = ef
		}
	}
}
