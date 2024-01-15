package constants

type DynamoDBTable string
type DynamoDBField string
type AWSRegions string

const (
	ReportTable DynamoDBTable = "REPORT_TABLE"
)

const (
	USEast2 AWSRegions = "us-east-2"
)

const (
	ReportIDField   DynamoDBField = "ReportID"
	ReportTypeField DynamoDBField = "ReportType"
	TitleField      DynamoDBField = "Title"
	CityField       DynamoDBField = "City"
	PartsField      DynamoDBField = "Parts"
	HeadersField    DynamoDBField = "Headers"
	SectionsField   DynamoDBField = "Sections"
	QuestionsField  DynamoDBField = "Questions"
	IndexField      DynamoDBField = "Index"
)

func (field DynamoDBField) String() string {
	return string(field)
}
