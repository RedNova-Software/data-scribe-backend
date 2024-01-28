package constants

type ItemType string

const (
	USEast2 string = "us-east-2"
)

// Types of items
const (
	Report   ItemType = "report"
	Template ItemType = "template"
)

// Item Fields
const (
	ReportIDField   string = "ReportID"
	TemplateIDField string = "TemplateID"
	ReportTypeField string = "ReportType"
	TitleField      string = "Title"
	CityField       string = "City"
	PartsField      string = "Parts"
	HeadersField    string = "Headers"
	SectionsField   string = "Sections"
	QuestionsField  string = "Questions"
	IndexField      string = "Index"
)

const (
	OwnerUserIDField   string = "OwnedBy.UserID"
	SharedWithIDsField string = "SharedWithIDs"
)

const (
	CreatedField      string = "Created"
	LastModifiedField string = "LastModified"
)
