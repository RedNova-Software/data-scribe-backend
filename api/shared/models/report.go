package models

type Question struct {
	Question  string
	Answer    string
	DataLabel string
}

type Section struct {
	Title     string
	Questions []Question
}

type Header struct {
	Title string
	Sections []Section
}

type Part struct {
	Title    string
	Headers  []Header
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
	ReportID    string
	ReportType  string
	Title 		string
	City        string
	Parts       []Part
}