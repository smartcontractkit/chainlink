// A programming model for software-defined state
package constructs


// Options for `construct.addMetadata()`.
type MetadataOptions struct {
	// Include stack trace with metadata entry.
	StackTrace *bool `field:"optional" json:"stackTrace" yaml:"stackTrace"`
	// A JavaScript function to begin tracing from.
	//
	// This option is ignored unless `stackTrace` is `true`.
	TraceFromFunction interface{} `field:"optional" json:"traceFromFunction" yaml:"traceFromFunction"`
}

