package models

type Answer struct {
	QuestionIndex uint16
	Answer        string
}

type Question struct {
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

type TextOutput struct {
	Title  string
	Index  uint16
	Type   TextOutputType
	Input  string
	Result string
}

type Section struct {
	Title           string
	Index           uint16
	OutputGenerated bool
	Questions       []Question
	TextOutputs     []TextOutput
}

type Part struct {
	Title    string
	Index    uint16
	Sections []Section
}

type ModelInfo struct {
}

type Model struct {
	ModelType string
	ModelInfo ModelInfo
}

type ModelOutput struct {
	TextOutputs []TextOutput
	Models      []Model
}

type Report struct {
	ReportID   string
	ReportType string
	Title      string
	City       string
	Parts      []Part
}
