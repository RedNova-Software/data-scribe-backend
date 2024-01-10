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
      billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
    });

    const addReportLambda = new lambda.Function(this, 'AddReportLambda', {
      code: lambda.Code.fromAsset(path.join(__dirname, '../bin/lambdas/post-add-report')),
      handler: 'main',
      runtime: lambda.Runtime.PROVIDED_AL2023,
      environment: {
        REPORT_TABLE: reportTable.tableName,
      },
    });

    const addPartLambda = new lambda.Function(this, 'AddPartLambda',  {
      code: lambda.Code.fromAsset(path.join(__dirname, '../bin/lambdas/post-add-part')),
        handler: 'main',
        runtime: lambda.Runtime.PROVIDED_AL2023,
        environment: {
          REPORT_TABLE: reportTable.tableName,
        },
    });

    reportTable.grantReadWriteData(addReportLambda)
    reportTable.grantReadWriteData(addPartLambda)
    
    const gateway = new apigateway.RestApi(this, "RedNovaGateway", {
      defaultCorsPreflightOptions: {
        allowOrigins: ["*"],
        allowMethods: ["POST"],
      }
    });

    const addReportIntegration = new apigateway.LambdaIntegration(addReportLambda);
    const addPartIntegration = new apigateway.LambdaIntegration(addPartLambda);

    const reportResource = gateway.root.addResource('report');
    const addResource = reportResource.addResource('add');
    const partResource = addResource.addResource('part');

    reportResource.addMethod("POST", addReportIntegration); 
    partResource.addMethod("POST", addPartIntegration);
    
  }

}
