import * as cdk from 'aws-cdk-lib'
import * as dynamodb from 'aws-cdk-lib/aws-dynamodb'
import { type Construct } from 'constructs'

export class DynamoDBStack extends cdk.Stack {
  public readonly reportTable: dynamodb.Table

  constructor (scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props)

    this.reportTable = new dynamodb.Table(this, 'ReportTable', {
      partitionKey: { name: 'ReportID', type: dynamodb.AttributeType.STRING },
      billingMode: dynamodb.BillingMode.PAY_PER_REQUEST
    })
  }
}
