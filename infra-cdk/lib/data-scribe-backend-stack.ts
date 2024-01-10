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

    const createNewReportLambda = new lambda.Function(this, 'CreateNewReportLambda', {
      code: lambda.Code.fromAsset(path.join(__dirname, '../bin/lambdas/create-empty-report')),
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

    reportTable.grantReadWriteData(createNewReportLambda)
    reportTable.grantReadData(getReportByIDLambda)
    reportTable.grantReadData(getAllReportsLambda)
    
    
    const gateway = new apigateway.RestApi(this, "DataScribeGateway", {
      defaultCorsPreflightOptions: {
        allowOrigins: ["*"],
        allowMethods: ["POST", "GET"],
      }
    })

    const reportEndpoints = gateway.root.addResource('reports')


    const createNewReportEndpoint = reportEndpoints.addResource('create')
    createNewReportEndpoint.addMethod("POST", new apigateway.LambdaIntegration(createNewReportLambda))
    

    const getReportByIDEndpoint = reportEndpoints.addResource('get')
    getReportByIDEndpoint.addMethod("GET", new apigateway.LambdaIntegration(getReportByIDLambda), {
      requestParameters: {
        'method.request.querystring.reportID': true,
      },
    })

    const getAllReportsEndpoint = reportEndpoints.addResource('getAll')
    getAllReportsEndpoint.addMethod("GET", new apigateway.LambdaIntegration(getAllReportsLambda))
  }

}
