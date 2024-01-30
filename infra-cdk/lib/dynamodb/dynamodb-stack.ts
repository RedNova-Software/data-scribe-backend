import * as cdk from "aws-cdk-lib";
import * as dynamodb from "aws-cdk-lib/aws-dynamodb";
import { type Construct } from "constructs";
import { DynamoDBTable, Tables } from "../constants/dynamodb-constants";

export class DynamoDBStack extends cdk.Stack {
  public readonly reportTable: dynamodb.Table;
  public readonly templateTable: dynamodb.Table;

  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    this.reportTable = new dynamodb.Table(this, DynamoDBTable.ReportTable, {
      partitionKey: {
        name: Tables.ReportID,
        type: dynamodb.AttributeType.STRING,
      },
      timeToLiveAttribute: Tables.DeleteAt,
      billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
      pointInTimeRecovery: true,
    });

    this.templateTable = new dynamodb.Table(this, DynamoDBTable.TemplateTable, {
      partitionKey: {
        name: Tables.TemplateID,
        type: dynamodb.AttributeType.STRING,
      },
      timeToLiveAttribute: Tables.DeleteAt,
      billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
      pointInTimeRecovery: true,
    });
  }
}
