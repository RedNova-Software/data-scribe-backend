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

type ReportSection struct {
	Title           string
	OutputGenerated bool
	Questions       []ReportQuestion
	TextOutputs     []ReportTextOutput
}

type ReportPart struct {
	Title    string
	Sections []ReportSection
}

type ModelInfo struct {
}

type Model struct {
	ModelType string
	ModelInfo ModelInfo
}

type ModelOutput struct {
	TextOutputs []ReportTextOutput
	Models      []Model
}

type Report struct {
	ReportID      string
	ReportType    string
	Title         string
	City          string
	Parts         []ReportPart
	OwnedBy       User
	SharedWithIDs []string
	Created       int64
	LastModified  int64
}

type ReportMetadata struct {
	ReportID     string
	ReportType   string
	Title        string
	City         string
	OwnedBy      User
	SharedWith   []User
	Created      int64
	LastModified int64
}
