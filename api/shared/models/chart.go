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

type OneDimConfig struct {
	Column string // The actual column in the csv

	Description string

	OperationType ChartOperation

	AggregateValueLabel string // The name of the label in the output

	AcceptedValues []string // Optional
}

type TwoDimConfig struct {
	IndependentColumn               string
	IndependantColumnAcceptedValues []string // Optional
	DependentColumns                []OneDimConfig
}
