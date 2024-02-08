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
	OwnedByField       string = "OwnedBy"
	OwnedByUserIDField string = OwnedByField + "." + "UserID"
	SharedWithIDsField string = "SharedWithIDs"
)

const (
	CreatedAtField      string = "CreatedAt"
	LastModifiedAtField string = "LastModifiedAt"
)

const (
	CSVIDField           string = "CSVID"
	CSVColumnsS3KeyField string = "CSVColumnsS3Key"
)

const (
	IsDeletedField string = "IsDeleted"
	DeleteAtField  string = "DeleteAt"
)

const (
	OperationIDField        string = "OperationID"
	OperationCompletedField string = "Completed"
)

const GlobalQuestionsField string = "GlobalQuestions"
