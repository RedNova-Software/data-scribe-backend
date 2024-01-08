import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as lambda from "aws-cdk-lib/aws-lambda";
import * as dynamodb from "aws-cdk-lib/aws-dynamodb";
import * as apigateway from "aws-cdk-lib/aws-apigateway";
import path = require('path');
import { DynamoDBTable, ReportTable } from './constants/dynamodb-constants';
import * as iam from 'aws-cdk-lib/aws-iam';

export class DataScribeBackendStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    const reportTable = new dynamodb.Table(this, 'ReportTable', {
      partitionKey: { name: 'ReportID', type: dynamodb.AttributeType.STRING },
      sortKey: { name: 'ReportType', type: dynamodb.AttributeType.STRING },
      billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
    });

    const reportLambda = new lambda.Function(this, 'ReportLambda', {
      code: lambda.Code.fromAsset(path.join(__dirname, '../bin/lambdas/create-empty-report-endpoint')),
      handler: 'main',
      runtime: lambda.Runtime.PROVIDED_AL2023,
      environment: {
        REPORT_TABLE: reportTable.tableName,
      },
    });

    reportTable.grantReadWriteData(reportLambda)
    
    const gateway = new apigateway.RestApi(this, "RedNovaGateway", {
      defaultCorsPreflightOptions: {
        allowOrigins: ["*"],
        allowMethods: ["POST"],
      }
    })

    const integration = new apigateway.LambdaIntegration(reportLambda)
    const reportEndpoint = gateway.root.addResource('report').addResource('create')
   
    reportEndpoint.addMethod("POST", integration)
    
  }

}
