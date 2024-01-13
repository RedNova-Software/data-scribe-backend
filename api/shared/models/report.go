package models

type Question struct {
	Question  string
	Answer    string
	DataLabel string
}

type SubSection struct {
	Title     string
	Questions []Question
}

type Section struct {
	Title       string
	SubSections []SubSection
}

type Part struct {
	Title    string
	Index    uint16
	Sections []Section
}

type TextOutput struct {
	Name   string
	Output string
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
