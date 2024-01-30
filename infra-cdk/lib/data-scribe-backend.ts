#!/usr/bin/env node
import "source-map-support/register";
import * as cdk from "aws-cdk-lib";
import { accountID } from "./constants/env-constants";
import { DynamoDBStack } from "./dynamodb/dynamodb-stack";
import { LambdasStack } from "./lambda/lambda-stack";
import { GatewayStack } from "./gateway/gateway-stack";
import { CognitoUserPoolStack } from "./cognito/cognito-stack";
import { S3BucketStack } from "./s3/s3-stack";

const app = new cdk.App();
const env = {
  account: accountID,
  region: "us-east-2",
};

const dynamoDBStack = new DynamoDBStack(app, "DynamoDBStack", { env });

// Instantiate the Cognito User Pool Stack
const cognitoStack = new CognitoUserPoolStack(app, "CognitoStack", {
  env: env, // Specify the account and region
});

const s3BucketStack = new S3BucketStack(app, "CSVBucketStack", {
  env: env, // Specify the account and region
});

const lambdaFunctionsStack = new LambdasStack(app, "LambdaStack", {
  env,
  reportTable: dynamoDBStack.reportTable,
  templateTable: dynamoDBStack.templateTable,
  userPool: cognitoStack.userPool,
  csvBucket: s3BucketStack.csvBucket,
});

const apiGatewayStack = new GatewayStack(app, "GatewayStack", {
  env,
  // Report Lambdas
  getReportByIDLambda: lambdaFunctionsStack.getReportByIDLambda,
  getAllReportsLambda: lambdaFunctionsStack.getAllReportsLambda,
  createReportLambda: lambdaFunctionsStack.createReportLambda,
  generateSectionLambda: lambdaFunctionsStack.generateSectionLambda,
  getAllReportTypesLambda: lambdaFunctionsStack.getAllReportTypesLambda,
  uploadCSVLambda: lambdaFunctionsStack.uploadCSVLambda,

  // Template Lambdas
  getTemplateByIDLambda: lambdaFunctionsStack.getTemplateByIDLambda,
  getAllTemplatesLambda: lambdaFunctionsStack.getAllTemplatesLambda,
  createTemplateLambda: lambdaFunctionsStack.createTemplateLambda,

  // Shared Lambdas
  addPartLambda: lambdaFunctionsStack.addPartLambda,
  deletePartLambda: lambdaFunctionsStack.deletePartLambda,
  addSectionLambda: lambdaFunctionsStack.addSectionLambda,
  deleteSectionLambda: lambdaFunctionsStack.deleteSectionLambda,
  updatePartLambda: lambdaFunctionsStack.updatePartLambda,
  updateSectionLambda: lambdaFunctionsStack.updateSectionLambda,
  updateItemTitleLambda: lambdaFunctionsStack.updateItemTitleLambda,
  shareItemLambda: lambdaFunctionsStack.shareItemLambda,
  convertItemLambda: lambdaFunctionsStack.convertItemLambda,
  deleteItemLambda: lambdaFunctionsStack.deleteItemLambda,

  // User Lambdas
  getUserIDLambda: lambdaFunctionsStack.getUserIDLambda,
  getAllUsersLambda: lambdaFunctionsStack.getAllUsersLambda,

  // User Pool
  userPool: cognitoStack.userPool,
});
