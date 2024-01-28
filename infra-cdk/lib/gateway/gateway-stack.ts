import * as cdk from "aws-cdk-lib";
import * as cognito from "aws-cdk-lib/aws-cognito";
import { type Construct } from "constructs";
import type * as lambda from "aws-cdk-lib/aws-lambda";
import * as apigateway from "aws-cdk-lib/aws-apigateway";
import * as logs from "aws-cdk-lib/aws-logs";
import * as iam from "aws-cdk-lib/aws-iam";
import path = require("path");

interface GatewayStackProps extends cdk.StackProps {
  // Report Lambdas
  getReportByIDLambda: lambda.IFunction;
  getAllReportsLambda: lambda.IFunction;
  createReportLambda: lambda.IFunction;
  generateSectionLambda: lambda.IFunction;
  getAllReportTypesLambda: lambda.IFunction;

  // Template Lambas
  getTemplateByIDLambda: lambda.IFunction;
  getAllTemplatesLambda: lambda.IFunction;
  createTemplateLambda: lambda.IFunction;

  // Shared Lambdas
  addPartLambda: lambda.IFunction;
  addSectionLambda: lambda.IFunction;
  updatePartLambda: lambda.IFunction;
  updateSectionLambda: lambda.IFunction;
  updateItemTitleLambda: lambda.IFunction;
  shareItemLambda: lambda.IFunction;

  // User Lambdas
  getUserIDLambda: lambda.IFunction;
  getAllUsersLambda: lambda.IFunction;

  // Cognito User Pool
  userPool: cognito.UserPool;
}

export class GatewayStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props: GatewayStackProps) {
    super(scope, id, props);

    const logGroup = new logs.LogGroup(this, "ApiGatewayLogGroup", {
      retention: logs.RetentionDays.ONE_MONTH, // Set the retention as needed
    });

    const stageOptions: apigateway.StageOptions = {
      loggingLevel: apigateway.MethodLoggingLevel.INFO,
      dataTraceEnabled: false,
      metricsEnabled: true,
      tracingEnabled: false, // For X-Ray tracing
      accessLogDestination: new apigateway.LogGroupLogDestination(logGroup),
      accessLogFormat: apigateway.AccessLogFormat.clf(), // Common Log Format
    };

    const gateway = new apigateway.RestApi(this, "DataScribeGateway", {
      defaultCorsPreflightOptions: {
        allowOrigins: apigateway.Cors.ALL_ORIGINS,
        allowMethods: apigateway.Cors.ALL_METHODS,
      },
      cloudWatchRole: true, // Needed to output logs
      deployOptions: stageOptions,
    });

    const userPool = cognito.UserPool.fromUserPoolId(
      this,
      "DataScribeUserPool",
      props.userPool.userPoolId
    );

    const authorizer = new apigateway.CognitoUserPoolsAuthorizer(
      this,
      "CognitoAuthorizer",
      {
        cognitoUserPools: [userPool],
        identitySource: "method.request.header.Authorization", // default
      }
    );

    const userResource = gateway.root.addResource("users");

    const reportResource = gateway.root.addResource("reports");
    const reportPartResource = reportResource.addResource("parts");
    const reportSectionsResource = reportPartResource.addResource("sections");

    const templateResource = gateway.root.addResource("templates");
    const templatePartResource = templateResource.addResource("parts");
    const templateSectionsResource =
      templatePartResource.addResource("sections");

    const sharedResource = gateway.root.addResource("shared");
    const sharedPartResource = sharedResource.addResource("parts");
    const sharedSectionResource = sharedPartResource.addResource("sections");

    // Report Endpoints

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

    const createNewReportEndpoint = reportResource.addResource("create");
    createNewReportEndpoint.addMethod(
      "POST",
      new apigateway.LambdaIntegration(props.createReportLambda),
      {
        authorizer,
        authorizationType: apigateway.AuthorizationType.COGNITO,
      }
    );

    const generateSectionEndpoint =
      reportSectionsResource.addResource("generate");
    generateSectionEndpoint.addMethod(
      "POST",
      new apigateway.LambdaIntegration(props.generateSectionLambda),
      {
        authorizer,
        authorizationType: apigateway.AuthorizationType.COGNITO,
      }
    );
    // --------------------------------------------------------- //
    // Template Endpoints

    const getTemplateByIDEndpoint = templateResource.addResource("get");
    getTemplateByIDEndpoint.addMethod(
      "GET",
      new apigateway.LambdaIntegration(props.getTemplateByIDLambda),
      {
        authorizer,
        authorizationType: apigateway.AuthorizationType.COGNITO,
        requestParameters: {
          "method.request.querystring.templateID": true,
        },
      }
    );

    const getAllTemplatesEndpoint = templateResource.addResource("all");
    getAllTemplatesEndpoint.addMethod(
      "GET",
      new apigateway.LambdaIntegration(props.getAllTemplatesLambda),
      {
        authorizer,
        authorizationType: apigateway.AuthorizationType.COGNITO,
      }
    );

    const createNewTemplateEndpoint = templateResource.addResource("create");
    createNewTemplateEndpoint.addMethod(
      "POST",
      new apigateway.LambdaIntegration(props.createTemplateLambda),
      {
        authorizer,
        authorizationType: apigateway.AuthorizationType.COGNITO,
      }
    );

    // --------------------------------------------------------- //
    // Shared Endpoints.

    const updateItemTitleEndpoint = sharedResource.addResource("title");
    updateItemTitleEndpoint.addMethod(
      "PUT",
      new apigateway.LambdaIntegration(props.updateItemTitleLambda),
      {
        authorizer,
        authorizationType: apigateway.AuthorizationType.COGNITO,
      }
    );

    const addPartEndpoint = sharedPartResource.addResource("add");
    addPartEndpoint.addMethod(
      "POST",
      new apigateway.LambdaIntegration(props.addPartLambda),
      {
        authorizer,
        authorizationType: apigateway.AuthorizationType.COGNITO,
      }
    );

    const addSectionEndpoint = sharedSectionResource.addResource("add");
    addSectionEndpoint.addMethod(
      "POST",
      new apigateway.LambdaIntegration(props.addSectionLambda),
      {
        authorizer,
        authorizationType: apigateway.AuthorizationType.COGNITO,
      }
    );

    const updatePartEndpoint = sharedPartResource.addResource("update");
    updatePartEndpoint.addMethod(
      "PUT",
      new apigateway.LambdaIntegration(props.updatePartLambda),
      {
        authorizer,
        authorizationType: apigateway.AuthorizationType.COGNITO,
      }
    );

    const updateSectionEndpoint = sharedSectionResource.addResource("update");
    updateSectionEndpoint.addMethod(
      "PUT",
      new apigateway.LambdaIntegration(props.updateSectionLambda),
      {
        authorizer,
        authorizationType: apigateway.AuthorizationType.COGNITO,
      }
    );

    const shareItemEndpoint = sharedResource.addResource("share");
    shareItemEndpoint.addMethod(
      "PUT",
      new apigateway.LambdaIntegration(props.shareItemLambda),
      {
        authorizer,
        authorizationType: apigateway.AuthorizationType.COGNITO,
      }
    );

    // --------------------------------------------------------- //
    // User Endpoints

    const getUserIDEndpoint = userResource.addResource("getCurrentID");
    getUserIDEndpoint.addMethod(
      "GET",
      new apigateway.LambdaIntegration(props.getUserIDLambda),
      {
        authorizer,
        authorizationType: apigateway.AuthorizationType.COGNITO,
      }
    );

    const getAllUsersEndpoint = userResource.addResource("all");
    getAllUsersEndpoint.addMethod(
      "GET",
      new apigateway.LambdaIntegration(props.getAllUsersLambda),
      {
        authorizer,
        authorizationType: apigateway.AuthorizationType.COGNITO,
      }
    );

    // --------------------------------------------------------- //
  }
}
