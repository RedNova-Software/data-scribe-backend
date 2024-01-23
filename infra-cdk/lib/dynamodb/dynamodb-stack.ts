import * as cdk from "aws-cdk-lib";
import * as dynamodb from "aws-cdk-lib/aws-dynamodb";
import { type Construct } from "constructs";
import { DynamoDBTable, ReportTable } from "../constants/dynamodb-constants";

export class DynamoDBStack extends cdk.Stack {
  public readonly reportTable: dynamodb.Table;
  public readonly templateTable: dynamodb.Table;

  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    this.reportTable = new dynamodb.Table(this, DynamoDBTable.ReportTable, {
      partitionKey: {
        name: ReportTable.ReportID,
        type: dynamodb.AttributeType.STRING,
      },
      billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
    });

    this.templateTable = new dynamodb.Table(this, DynamoDBTable.TemplateTable, {
      partitionKey: {
        name: ReportTable.TemplateID,
        type: dynamodb.AttributeType.STRING,
      },
      billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
    });
  }
}
