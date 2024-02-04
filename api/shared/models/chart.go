package models

type ChartOperation string

const (
	NumericalSum          ChartOperation = "NumericalSum"
	Average               ChartOperation = "Average"
	UniqueOccurrences     ChartOperation = "UniqueOccurrences"
	SetElementOccurrences ChartOperation = "SetElementOccurrences"
)

type CSVDataType string

const (
	OneDim CSVDataType = "OneDim"
	TwoDim CSVDataType = "TwoDim"
)

type TemplateOneDimConfig struct {
	AggregateValueLabel string

	Description string

	OperationType ChartOperation
}

type ReportOneDimConfig struct {
	AggregateValueLabel string // The name of the label in the output

	Column string // The actual column in the csv

	Description string

	OperationType ChartOperation

	AcceptedValues []string // Optional
}

type TemplateTwoDimConfig struct {
	IndependentColumnLabel string
	DependentColumns       []TemplateOneDimConfig
}

type ReportTwoDimConfig struct {
	IndependentColumnLabel          string
	IndependentColumn               string
	IndependantColumnAcceptedValues []string // Optional
	DependentColumns                []ReportOneDimConfig
}

type ChartType string

const (
	Line    ChartType = "Line"
	Area    ChartType = "Area"
	Bar     ChartType = "Bar"
	Scatter ChartType = "Scatter"
	Pie     ChartType = "Pie"
	Radar   ChartType = "Radar"
)

type CsvDataColumnUniqueValuesMap map[string][]string
