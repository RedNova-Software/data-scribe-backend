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

type TemplateChartOutput struct {
	Title         string
	XAxisTitle    string
	YAxisTitle    string // Optional
	CartesianGrid bool

	Config TwoDimConfig
}

type TemplateCSVData struct {
	Label        string
	Type         CSVDataType
	ConfigOneDim OneDimConfig
	ConfigTwoDim TwoDimConfig
}

type TemplateSection struct {
	Title        string
	Questions    []TemplateQuestion
	CSVData      []TemplateCSVData
	TextOutputs  []TemplateTextOutput
	ChartOutputs []TemplateChartOutput
}

type TemplatePart struct {
	Title    string
	Sections []TemplateSection
}

type Template struct {
	TemplateID     string
	Title          string
	Parts          []TemplatePart
	OwnedBy        User
	SharedWithIDs  []string
	LastModifiedAt int64
	CreatedAt      int64
	IsDeleted      bool
	DeleteAt       int64

	CSVColumns map[string]string
}

type TemplateMetadata struct {
	TemplateID     string
	Title          string
	OwnedBy        User
	SharedWith     []User
	CreatedAt      int64
	LastModifiedAt int64
}
