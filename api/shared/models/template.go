package models

type TemplateQuestion struct {
	Label    string
	Index    uint16
	Question string
}

type TemplateTextOutput struct {
	Title string
	Index uint16
	Type  TextOutputType
	Input string
}

type TemplateSection struct {
	Title       string
	Index       uint16
	Questions   []TemplateQuestion
	TextOutputs []TemplateTextOutput
}

type TemplatePart struct {
	Title    string
	Index    uint16
	Sections []TemplateSection
}

type Template struct {
	TemplateID string
	Title      string
	Parts      []TemplatePart
}
