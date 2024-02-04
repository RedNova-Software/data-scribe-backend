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
}

type CsvDataResponse struct {
	GroupColumn              string
	GroupColumnAcceptedValue string // Needed

	DependentColumn OneDimConfigResponse
}
