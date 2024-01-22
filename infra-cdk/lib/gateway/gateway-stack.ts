import * as cdk from "aws-cdk-lib";
import * as cognito from "aws-cdk-lib/aws-cognito";
import { type Construct } from "constructs";
import type * as lambda from "aws-cdk-lib/aws-lambda";
import * as apigateway from "aws-cdk-lib/aws-apigateway";
import path = require("path");
import { userPoolId } from "../constants/cognito-constants";

interface GatewayStackProps extends cdk.StackProps {
  getReportByIDLambda: lambda.IFunction;
  getAllReportsLambda: lambda.IFunction;
  createReportLambda: lambda.IFunction;
  addPartLambda: lambda.IFunction;
  addSectionLambda: lambda.IFunction;
  generateSectionLambda: lambda.IFunction;
  getAllReportTypesLambda: lambda.IFunction;
}

export class GatewayStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props: GatewayStackProps) {
    super(scope, id, props);

    const gateway = new apigateway.RestApi(this, "DataScribeGateway", {
      defaultCorsPreflightOptions: {
        allowOrigins: ["*"],
        allowMethods: ["POST", "GET"],
      },
    });

    const userPool = cognito.UserPool.fromUserPoolId(
      this,
      "DataScribeUserPool",
      userPoolId
    );

    const authorizer = new apigateway.CognitoUserPoolsAuthorizer(
      this,
      "CognitoAuthorizer",
      {
        cognitoUserPools: [userPool],
        identitySource: "method.request.header.Authorization", // default
      }
    );

    const reportResource = gateway.root.addResource("reports");
    const partResource = reportResource.addResource("parts");
    const sectionsResource = partResource.addResource("sections");

    // Get endpoints

    const getReportByIDEndpoint = reportResource.addResource("get");
    getReportByIDEndpoint.addMethod(
      "GET",
      new apigateway.LambdaIntegration(props.getReportByIDLambda),
      {
        authorizer,
        authorizationType: apigateway.AuthorizationType.COGNITO,
        requestParameters: {
          "method.request.querystring.reportID": true,
        },
      }
    );

    const getAllReportsEndpoint = reportResource.addResource("all");
    getAllReportsEndpoint.addMethod(
      "GET",
      new apigateway.LambdaIntegration(props.getAllReportsLambda),
      {
        authorizer,
        authorizationType: apigateway.AuthorizationType.COGNITO,
      }
    );

    const getAllReportTypesEndpoint = reportResource.addResource("types");
    getAllReportTypesEndpoint.addMethod(
      "GET",
      new apigateway.LambdaIntegration(props.getAllReportTypesLambda),
      {
        authorizer,
        authorizationType: apigateway.AuthorizationType.COGNITO,
      }
    );

    // Post endpoints

    const createNewReportEndpoint = reportResource.addResource("create");
    createNewReportEndpoint.addMethod(
      "POST",
      new apigateway.LambdaIntegration(props.createReportLambda),
      {
        authorizer,
        authorizationType: apigateway.AuthorizationType.COGNITO,
      }
    );

    const addPartEndpoint = partResource.addResource("add");
    addPartEndpoint.addMethod(
      "POST",
      new apigateway.LambdaIntegration(props.addPartLambda),
      {
        authorizer,
        authorizationType: apigateway.AuthorizationType.COGNITO,
      }
    );

    const addSectionEndpoint = sectionsResource.addResource("add");
    addSectionEndpoint.addMethod(
      "POST",
      new apigateway.LambdaIntegration(props.addSectionLambda),
      {
        authorizer,
        authorizationType: apigateway.AuthorizationType.COGNITO,
      }
    );

    const generateSectionEndpoint = sectionsResource.addResource("generate");
    generateSectionEndpoint.addMethod(
      "POST",
      new apigateway.LambdaIntegration(props.generateSectionLambda),
      {
        authorizer,
        authorizationType: apigateway.AuthorizationType.COGNITO,
      }
    );
  }
}
