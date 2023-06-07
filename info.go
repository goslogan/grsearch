package grstack

import "time"

/****
* data structures for the FT.info command
*****/

type Info struct {
	IndexName                  string           `mapstructure:"index_name"`
	Index                      IndexOptions     `mapstructure:"-"`
	NumDocs                    int64            `mapstructure:"num_docs"`
	MaxDocId                   int64            `mapstructure:"max_doc_id"`
	MaxTerms                   int64            `mapstructure:"max_terms"`
	NumRecords                 int64            `mapstructure:"num_records"`
	Indexing                   float64          `mapstructure:"indexing"`
	PercentIndexed             float64          `mapstructure:"percent_indexed"`
	HashIndexingFailures       int64            `mapstructure:"hash_indexing_failures"`
	TotalInvertedIndexBlocks   int64            `mapstructure:"total_inverted_index_blocks"`
	InvertSize                 float64          `mapstructure:"invert_sz_mb"`
	VectorIndexSize            float64          `mapstructure:"vector_index_sz_mb"`
	DocTableSize               float64          `mapstructure:"doc_table_size_mb"`
	OffsetVectorsSize          float64          `mapstructure:"offset_vectors_sz_mb"`
	SortableValuesSize         float64          `mapstructure:"sortable_values_size_mb"`
	KeyTableSize               float64          `mapstructure:"key_table_size_mb"`
	AverageRecordsPerDoc       float64          `mapstructure:"records_per_doc_avg"`
	AverageBytesPerRecord      float64          `mapstructure:"bytes_per_record_avg"`
	AverageOffsetsPerTerm      float64          `mapstructure:"offsets_per_term_avg"`
	AverageOffsetBitsPerRecord float64          `mapstructure:"offset_bits_per_record_avg"`
	TotalIndexingTime          time.Duration    `mapstructure:"total_indexing_time"`
	NumberOfUses               int64            `mapstructure:"number_of_uses"`
	GCStats                    InfoGCStats      `mapstructure:"gc_stats"`
	InfoCursorStats            InfoCursorStats  `mapstructure:"cursor_stats"`
	DialectStats               InfoDialectStats `mapstructure:"dialect_stats"`
}

type InfoGCStats struct {
	BytesCollected       int64         `mapstructure:"bytes_collected"`
	TotalMsRun           time.Duration `mapstructure:"total_ms_run"`
	TotalCycles          int64         `mapstructure:"total_cycle"`
	AverageCycleTime     time.Duration `mapstructure:"average_cycle_time"`
	LastRunTime          time.Duration `mapstructure:"last_run_time"`
	GCNumericTreesMissed int64         `mapstructure:"gc_numeric_trees_missed"`
	GCBlocksDenied       int64         `mapstructure:"gc_blocks_denied"`
}

type InfoDialectStats struct {
	Dialect1 int64 `mapstructure:"dialect_1"`
	Dialect2 int64 `mapstructure:"dialect_2"`
	Dialect3 int64 `mapstructure:"dialect_3"`
}

type InfoCursorStats struct {
	GlobalIdle    int64 `mapstructure:"global_idle"`
	GlobalTotal   int64 `mapstructure:"global_total"`
	IndexCapacity int64 `mapstructure:"index_capacity"`
	IndexTotal    int64 //index_total
}
