package models

type ChartOperation string

const (
	NumericalSum         ChartOperation = "NumericalSum"
	Average              ChartOperation = "Average"
	UniqueOccurrences    ChartOperation = "UniqueOccurrences"
	SetElementOccurences ChartOperation = "SetElementOccurences"
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

	FilterColumns map[string][]string
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
