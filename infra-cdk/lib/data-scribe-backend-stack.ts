import * as cdk from "aws-cdk-lib";
import { Construct } from "constructs";
import * as lambda from "aws-cdk-lib/aws-lambda";
import * as dynamodb from "aws-cdk-lib/aws-dynamodb";
import * as apigateway from "aws-cdk-lib/aws-apigateway";
import path = require("path");
import * as fs from "fs";
import { DynamoDBTable, ReportTable } from "./constants/dynamodb-constants";
import * as iam from "aws-cdk-lib/aws-iam";
import { exit } from "process";

export class DataScribeBackendStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    const reportTable = new dynamodb.Table(this, "ReportTable", {
      partitionKey: { name: "ReportID", type: dynamodb.AttributeType.STRING },
      billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
    });

    const createReportLambda = new lambda.Function(this, "CreateReportLambda", {
      code: lambda.Code.fromAsset(
        path.join(__dirname, "../bin/lambdas/post-create-report")
      ),
      handler: "main",
      runtime: lambda.Runtime.PROVIDED_AL2023,
      environment: {
        REPORT_TABLE: reportTable.tableName,
      },
    });

    const getReportByIDLambda = new lambda.Function(
      this,
      "GetReportByIDLambda",
      {
        code: lambda.Code.fromAsset(
          path.join(__dirname, "../bin/lambdas/get-report-by-id")
        ),
        handler: "main",
        runtime: lambda.Runtime.PROVIDED_AL2023,
        environment: {
          REPORT_TABLE: reportTable.tableName,
        },
      }
    );

    const getAllReportsLambda = new lambda.Function(
      this,
      "GetAllReportsLambda",
      {
        code: lambda.Code.fromAsset(
          path.join(__dirname, "../bin/lambdas/get-all-reports")
        ),
        handler: "main",
        runtime: lambda.Runtime.PROVIDED_AL2023,
        environment: {
          REPORT_TABLE: reportTable.tableName,
        },
      }
    );

    const getAllReportTypesLambda = new lambda.Function(
      this,
      "GetAllReportTypesLambda",
      {
        code: lambda.Code.fromAsset(
          path.join(__dirname, "../bin/lambdas/get-all-report-types")
        ),
        handler: "main",
        runtime: lambda.Runtime.PROVIDED_AL2023,
        environment: {
          REPORT_TABLE: reportTable.tableName,
        },
      }
    );

    const addPartLambda = new lambda.Function(this, "AddPartLambda", {
      code: lambda.Code.fromAsset(
        path.join(__dirname, "../bin/lambdas/post-add-part")
      ),
      handler: "main",
      runtime: lambda.Runtime.PROVIDED_AL2023,
      environment: {
        REPORT_TABLE: reportTable.tableName,
      },
    });

    const addSectionLambda = new lambda.Function(this, "AddSectionLambda", {
      code: lambda.Code.fromAsset(
        path.join(__dirname, "../bin/lambdas/post-add-section")
      ),
      handler: "main",
      runtime: lambda.Runtime.PROVIDED_AL2023,
      environment: {
        REPORT_TABLE: reportTable.tableName,
      },
    });

    // Set openAI key as env variable
    let openAIKeyPath = path.join(__dirname, "../../keys/openai-key.txt");
    let openAIKey;

    fs.readFile(openAIKeyPath, "utf8", (err, data) => {
      if (err) {
        console.error(
          "\x1b[31m%s\x1b[0m",
          "You are most likely missing your openai key. Place it in /keys/openai-key.txt"
        );
        process.exit();
      }

      openAIKeyPath = data;
    });
    const generateSectionLambda = new lambda.Function(
      this,
      "GenerateSectionLambda",
      {
        code: lambda.Code.fromAsset(
          path.join(__dirname, "../bin/lambdas/post-generate-section")
        ),
        handler: "main",
        runtime: lambda.Runtime.PROVIDED_AL2023,
        environment: {
          REPORT_TABLE: reportTable.tableName,
          OPENAI_API_KEY: openAIKey!,
        },
      }
    );

    reportTable.grantReadData(getReportByIDLambda);
    reportTable.grantReadData(getAllReportsLambda);
    reportTable.grantWriteData(createReportLambda);
    reportTable.grantReadWriteData(addPartLambda);
    reportTable.grantReadWriteData(addSectionLambda);
    reportTable.grantReadWriteData(generateSectionLambda);

    const gateway = new apigateway.RestApi(this, "DataScribeGateway", {
      defaultCorsPreflightOptions: {
        allowOrigins: ["*"],
        allowMethods: ["POST", "GET"],
      },
    });

    const reportResource = gateway.root.addResource("reports");
    const partResource = reportResource.addResource("parts");
    const sectionsResource = partResource.addResource("sections");

    // Get endpoints

    const getReportByIDEndpoint = reportResource.addResource("get");
    getReportByIDEndpoint.addMethod(
      "GET",
      new apigateway.LambdaIntegration(getReportByIDLambda),
      {
        requestParameters: {
          "method.request.querystring.reportID": true,
        },
      }
    );

    const getAllReportsEndpoint = reportResource.addResource("all");
    getAllReportsEndpoint.addMethod(
      "GET",
      new apigateway.LambdaIntegration(getAllReportsLambda)
    );

    const getAllReportTypesEndpoint = reportResource.addResource("types");
    getAllReportTypesEndpoint.addMethod(
      "GET",
      new apigateway.LambdaIntegration(getAllReportTypesLambda)
    );

    // Post endpoints

    const createNewReportEndpoint = reportResource.addResource("create");
    createNewReportEndpoint.addMethod(
      "POST",
      new apigateway.LambdaIntegration(createReportLambda)
    );

    const addPartEndpoint = partResource.addResource("add");
    addPartEndpoint.addMethod(
      "POST",
      new apigateway.LambdaIntegration(addPartLambda)
    );

    const addSectionEndpoint = sectionsResource.addResource("add");
    addSectionEndpoint.addMethod(
      "POST",
      new apigateway.LambdaIntegration(addSectionLambda)
    );

    const generateSectionEndpoint = sectionsResource.addResource("generate");
    generateSectionEndpoint.addMethod(
      "POST",
      new apigateway.LambdaIntegration(generateSectionLambda)
    );
  }
}
