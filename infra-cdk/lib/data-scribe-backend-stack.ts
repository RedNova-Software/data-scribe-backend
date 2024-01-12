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


    const getReportByIDLambda = new lambda.Function(this, 'GetReportByIDLambda', {
      code: lambda.Code.fromAsset(path.join(__dirname, '../bin/lambdas/get-report-by-id')),
      handler: 'main',
      runtime: lambda.Runtime.PROVIDED_AL2023,
      environment: {
        REPORT_TABLE: reportTable.tableName,
      },
    });

    const getAllReportsLambda = new lambda.Function(this, 'GetAllReportsLambda', {
      code: lambda.Code.fromAsset(path.join(__dirname, '../bin/lambdas/get-all-reports')),
      handler: 'main',
      runtime: lambda.Runtime.PROVIDED_AL2023,
      environment: {
        REPORT_TABLE: reportTable.tableName,
      },
    });

    const getAllReportTypesLambda = new lambda.Function(this, 'GetAllReportTypesLambda', {
      code: lambda.Code.fromAsset(path.join(__dirname, '../bin/lambdas/get-all-report-types')),
      handler: 'main',
      runtime: lambda.Runtime.PROVIDED_AL2023,
    });

    reportTable.grantReadData(getReportByIDLambda)
    reportTable.grantReadData(getAllReportsLambda)
    reportTable.grantReadWriteData(addReportLambda)
    reportTable.grantReadWriteData(addPartLambda)
    
    
    const gateway = new apigateway.RestApi(this, "DataScribeGateway", {
      defaultCorsPreflightOptions: {
        allowOrigins: ["*"],
        allowMethods: ["POST", "GET"],
      }
    });


    const addReportIntegration = new apigateway.LambdaIntegration(addReportLambda);
    const addPartIntegration = new apigateway.LambdaIntegration(addPartLambda);

    const reportResource = gateway.root.addResource('report');
    const addResource = reportResource.addResource('add');
    const partResource = addResource.addResource('part');
    const getReportByIDResource = reportResource.addResource('get')

    reportResource.addMethod("POST", addReportIntegration); 
    partResource.addMethod("POST", addPartIntegration);
    getReportByIDResource.addMethod("GET", new apigateway.LambdaIntegration(getReportByIDLambda), {
      requestParameters: {
        'method.request.querystring.reportID': true,
      },
    })

    const getAllReportsEndpoint = reportResource.addResource('getAll')
    getAllReportsEndpoint.addMethod("GET", new apigateway.LambdaIntegration(getAllReportsLambda))
  
    
    const getAllReportTypesEndpoint = reportResource.addResource('types')
    getAllReportTypesEndpoint.addMethod("GET",new apigateway.LambdaIntegration(getAllReportTypesLambda))
  }

}
