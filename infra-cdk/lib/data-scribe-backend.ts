#!/usr/bin/env node
import "source-map-support/register";
import * as cdk from "aws-cdk-lib";
import { accountID } from "./constants/env-constants";
import { DynamoDBStack } from "./dynamodb/dynamodb-stack";
import { LambdasStack } from "./lambda/lambda-stack";
import { GatewayStack } from "./gateway/gateway-stack";

const app = new cdk.App();
const env = {
  account: accountID,
  region: "us-east-2",
};

const dynamoDBStack = new DynamoDBStack(app, "DynamoDBStack", { env });

const lambdaFunctionsStack = new LambdasStack(app, "LambdaStack", {
  env,
  reportTable: dynamoDBStack.reportTable,
  templateTable: dynamoDBStack.templateTable,
});

const apiGatewayStack = new GatewayStack(app, "GatewayStack", {
  env,
  // Report Lambdas
  getReportByIDLambda: lambdaFunctionsStack.getReportByIDLambda,
  getAllReportsLambda: lambdaFunctionsStack.getAllReportsLambda,
  createReportLambda: lambdaFunctionsStack.createReportLambda,
  generateSectionLambda: lambdaFunctionsStack.generateSectionLambda,
  getAllReportTypesLambda: lambdaFunctionsStack.getAllReportTypesLambda,

  // Template Lambdas
  getTemplateByIDLambda: lambdaFunctionsStack.getTemplateByIDLambda,
  getAllTemplatesLambda: lambdaFunctionsStack.getAllTemplatesLambda,
  createTemplateLambda: lambdaFunctionsStack.createTemplateLambda,

  // Shared Lambdas
  addPartLambda: lambdaFunctionsStack.addPartLambda,
  addSectionLambda: lambdaFunctionsStack.addSectionLambda,
});
