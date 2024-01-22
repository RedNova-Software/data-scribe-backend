#!/usr/bin/env node
import "source-map-support/register";
import * as cdk from "aws-cdk-lib";
import { accountID } from "../lib/constants/env-constants";
import { DynamoDBStack } from "../lib/dynamodb/dynamodb-stack";
import { LambdasStack } from "../lib/lambda/lambda-stack";
import { GatewayStack } from "../lib/gateway/gateway-stack";

const app = new cdk.App();
const env = {
  account: accountID,
  region: "us-east-2",
};

const dynamoDBStack = new DynamoDBStack(app, "DynamoDBStack", { env });

const lambdaFunctionsStack = new LambdasStack(app, "LambdaStack", {
  env,
  reportTable: dynamoDBStack.reportTable,
});

const apiGatewayStack = new GatewayStack(app, "GatewayStack", {
  env,
  getReportByIDLambda: lambdaFunctionsStack.getReportByIDLambda,
  getAllReportsLambda: lambdaFunctionsStack.getAllReportsLambda,
  createReportLambda: lambdaFunctionsStack.createReportLambda,
  addPartLambda: lambdaFunctionsStack.addPartLambda,
  addSectionLambda: lambdaFunctionsStack.addSectionLambda,
  generateSectionLambda: lambdaFunctionsStack.generateSectionLambda,
  getAllReportTypesLambda: lambdaFunctionsStack.getAllReportTypesLambda,
});
