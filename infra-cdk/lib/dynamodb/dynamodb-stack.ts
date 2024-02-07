import * as cdk from "aws-cdk-lib";
import * as dynamodb from "aws-cdk-lib/aws-dynamodb";
import { type Construct } from "constructs";
import { DynamoDBTable, TableFields } from "../constants/dynamodb-constants";

export class DynamoDBStack extends cdk.Stack {
  public readonly reportTable: dynamodb.Table;
  public readonly templateTable: dynamodb.Table;
  public readonly operationTable: dynamodb.Table;

  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    this.reportTable = new dynamodb.Table(this, DynamoDBTable.ReportTable, {
      partitionKey: {
        name: TableFields.ReportID,
        type: dynamodb.AttributeType.STRING,
      },
      timeToLiveAttribute: TableFields.DeleteAt,
      billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
      pointInTimeRecovery: true,
      deletionProtection: true,
    });

    this.reportTable.addGlobalSecondaryIndex({
      indexName: TableFields.CSVID,
      partitionKey: {
        name: TableFields.CSVID,
        type: dynamodb.AttributeType.STRING,
      },
      projectionType: dynamodb.ProjectionType.ALL,
    });

    this.templateTable = new dynamodb.Table(this, DynamoDBTable.TemplateTable, {
      partitionKey: {
        name: TableFields.TemplateID,
        type: dynamodb.AttributeType.STRING,
      },
      timeToLiveAttribute: TableFields.DeleteAt,
      billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
      pointInTimeRecovery: true,
      deletionProtection: true,
    });

    // This table stores ongoing operations for polling functions to check
    this.operationTable = new dynamodb.Table(
      this,
      DynamoDBTable.OperationsTable,
      {
        partitionKey: {
          name: TableFields.OperationID,
          type: dynamodb.AttributeType.STRING,
        },
        timeToLiveAttribute: TableFields.DeleteAt,
        billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
        deletionProtection: true,
      }
    );
  }
}
