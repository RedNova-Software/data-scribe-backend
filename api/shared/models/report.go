package models

type Answer struct {
	QuestionIndex uint16
	Answer        string
}

type ReportQuestion struct {
	Label    string
	Index    uint16
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
	Index  uint16
	Type   TextOutputType
	Input  string
	Result string
}

type ReportSection struct {
	Title           string
	Index           uint16
	OutputGenerated bool
	Questions       []ReportQuestion
	TextOutputs     []ReportTextOutput
}

type ReportPart struct {
	Title    string
	Index    uint16
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
	ReportID   string
	ReportType string
	Title      string
	City       string
	Parts      []ReportPart
}
