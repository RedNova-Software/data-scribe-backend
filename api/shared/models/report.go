package models

type ReportQuestion struct {
	Label    string
	Question string
	Answer   string
}

type TextOutputType string

const (
	Generator TextOutputType = "Generator"
	Static    TextOutputType = "Static"
)

type ReportTextOutput struct {
	Title  string
	Type   TextOutputType
	Input  string
	Result string
}

type ReportChartOutput struct {
	Title                  string
	Type                   ChartType
	Description            string
	XAxisTitle             string
	YAxisTitle             string // Optional
	CartesianGrid          bool
	IndependentColumnLabel string

	IndependentColumn string   // Actual column
	AcceptedValues    []string // Optional

	FilterColumns map[string][]string

	DependentColumns []ReportOneDimConfig

	Results []map[string]interface{}
}

type ReportCSVData struct {
	Label           string
	Description     string
	OperationType   ChartOperation
	OperationColumn string
	AcceptedValues  []string
	FilterColumns   map[string][]string
	Result          string
}

type ReportSection struct {
	Title           string
	OutputGenerated bool
	Questions       []ReportQuestion
	CSVData         []ReportCSVData
	TextOutputs     []ReportTextOutput
	ChartOutputs    []ReportChartOutput
}

type ReportPart struct {
	Title    string
	Sections []ReportSection
}

type Report struct {
	ReportID       string
	ReportType     string
	Title          string
	City           string
	Parts          []ReportPart
	OwnedBy        User
	SharedWithIDs  []string
	CreatedAt      int64
	LastModifiedAt int64
	IsDeleted      bool
	DeleteAt       int64
	CSVID          string

	CSVColumnsS3Key string

	GlobalQuestions []ReportQuestion
}

type ReportMetadata struct {
	ReportID       string
	ReportType     string
	Title          string
	City           string
	OwnedBy        User
	SharedWith     []User
	CreatedAt      int64
	LastModifiedAt int64
}
