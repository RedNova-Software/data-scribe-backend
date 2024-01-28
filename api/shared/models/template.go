package models

type TemplateQuestion struct {
	Label    string
	Question string
}

type TemplateTextOutput struct {
	Title string
	Type  TextOutputType
	Input string
}

type TemplateSection struct {
	Title       string
	Questions   []TemplateQuestion
	TextOutputs []TemplateTextOutput
}

type TemplatePart struct {
	Title    string
	Sections []TemplateSection
}

type Template struct {
	TemplateID    string
	Title         string
	Parts         []TemplatePart
	OwnedBy       User
	SharedWithIDs []string
}

type TemplateMetadata struct {
	TemplateID string
	Title      string
	OwnedBy    User
	SharedWith []User
}
