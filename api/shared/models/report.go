package models

type Answer struct {
	QuestionIndex int
	Answer        string
}

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
	Title         string
	XAxisTitle    string
	YAxisTitle    string // Optional
	CartesianGrid bool

	Config TwoDimConfig

	Results []map[string]interface{}
}

type ReportCSVData struct {
	Label        string
	Type         CSVDataType
	ConfigOneDim OneDimConfig
	ConfigTwoDim TwoDimConfig
	Result       string
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

	CSVColumns map[string][]string
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
