package models

type Answer struct {
	QuestionIndex uint16
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
	Questions       []ReportQuestion   `dynamodbav:"questions"`
	TextOutputs     []ReportTextOutput `dynamodbav:"textOutputs"`
}

type ReportPart struct {
	Title    string
	Sections []ReportSection `dynamodbav:"sections"`
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
	Parts      []ReportPart `dynamodbav:"parts"`
}
