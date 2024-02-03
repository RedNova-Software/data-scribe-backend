export enum DynamoDBTable {
  ReportTable = "ReportTable",
  TemplateTable = "TemplateTable",
}

export enum TableFields {
  ReportID = "ReportID",
  TemplateID = "TemplateID",
  // Attribute that stores a timestamp for when an item should actually be deleted
  DeleteAt = "DeleteAt",
  CSVID = "CSVID",
}
