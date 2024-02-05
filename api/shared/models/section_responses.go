package models

type Answer struct {
	Answer string
}

type OneDimConfigResponse struct {
	Column         string   // The actual column in the csv
	AcceptedValues []string // Optional
}

type ChartOutputResponse struct {
	IndependentColumn               string
	IndependentColumnAcceptedValues []string

	DependentColumns []OneDimConfigResponse
	FilterColumns    map[string][]string
}

type CsvDataResponse struct {
	OperationColumn               string   // The actual column in the csv
	OperationColumnAcceptedValues []string // Optional

	FilterColumns map[string][]string // Has a map of filter columns and their accepted values
}
