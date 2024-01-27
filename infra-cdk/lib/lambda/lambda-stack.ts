import * as cdk from "aws-cdk-lib";
import { type Construct } from "constructs";
import * as lambda from "aws-cdk-lib/aws-lambda";
import type * as dynamodb from "aws-cdk-lib/aws-dynamodb";
import path = require("path");
import * as fs from "fs";
import { memoryUsage } from "process";

interface LambdasStackProps extends cdk.StackProps {
  reportTable: dynamodb.Table;
  templateTable: dynamodb.Table;
}

export class LambdasStack extends cdk.Stack {
  // Report Lambdas
  public readonly getReportByIDLambda: lambda.IFunction;
  public readonly getAllReportsLambda: lambda.IFunction;
  public readonly createReportLambda: lambda.IFunction;
  public readonly generateSectionLambda: lambda.IFunction;
  public readonly getAllReportTypesLambda: lambda.IFunction;

  // Template Lambas
  public readonly getTemplateByIDLambda: lambda.IFunction;
  public readonly getAllTemplatesLambda: lambda.IFunction;
  public readonly createTemplateLambda: lambda.IFunction;

  // Shared Lambdas
  public readonly addPartLambda: lambda.IFunction;
  public readonly addSectionLambda: lambda.IFunction;
  public readonly updatePartLambda: lambda.IFunction;
  public readonly updateSectionLambda: lambda.IFunction;

  public readonly updateItemTitleLambda: lambda.IFunction;

  // --------------------------------------------------------- //

  constructor(scope: Construct, id: string, props: LambdasStackProps) {
    super(scope, id, props);

    // Report Lambdas

    this.createReportLambda = new lambda.Function(this, "CreateReportLambda", {
      code: lambda.Code.fromAsset(
        path.join(__dirname, "../../bin/lambdas/create-report")
      ),
      handler: "main",
      runtime: lambda.Runtime.PROVIDED_AL2023,
      environment: {
        REPORT_TABLE: props.reportTable.tableName,
      },
      memorySize: 1024,
    });
    props.reportTable.grantWriteData(this.createReportLambda);

    this.getReportByIDLambda = new lambda.Function(
      this,
      "GetReportByIDLambda",
      {
        code: lambda.Code.fromAsset(
          path.join(__dirname, "../../bin/lambdas/get-report-by-id")
        ),
        handler: "main",
        runtime: lambda.Runtime.PROVIDED_AL2023,
        environment: {
          REPORT_TABLE: props.reportTable.tableName,
        },
        memorySize: 1024,
      }
    );
    props.reportTable.grantReadData(this.getReportByIDLambda);

    this.getAllReportsLambda = new lambda.Function(
      this,
      "GetAllReportsLambda",
      {
        code: lambda.Code.fromAsset(
          path.join(__dirname, "../../bin/lambdas/get-all-reports")
        ),
        handler: "main",
        runtime: lambda.Runtime.PROVIDED_AL2023,
        environment: {
          REPORT_TABLE: props.reportTable.tableName,
        },
        memorySize: 1024,
      }
    );
    props.reportTable.grantReadData(this.getAllReportsLambda);

    this.getAllReportTypesLambda = new lambda.Function(
      this,
      "GetAllReportTypesLambda",
      {
        code: lambda.Code.fromAsset(
          path.join(__dirname, "../../bin/lambdas/get-all-report-types")
        ),
        handler: "main",
        runtime: lambda.Runtime.PROVIDED_AL2023,
        memorySize: 1024,
      }
    );

    // Set openAI key as env variable
    const openAIKeyPath = path.join(__dirname, "../../../keys/openai-key.txt");
    const openAIKey = fs.readFileSync(openAIKeyPath, {
      encoding: "utf8",
      flag: "r",
    });

    this.generateSectionLambda = new lambda.Function(
      this,
      "GenerateSectionLambda",
      {
        code: lambda.Code.fromAsset(
          path.join(__dirname, "../../bin/lambdas/generate-section")
        ),
        handler: "main",
        runtime: lambda.Runtime.PROVIDED_AL2023,
        environment: {
          REPORT_TABLE: props.reportTable.tableName,
          OPENAI_API_KEY: openAIKey,
        },
        timeout: cdk.Duration.minutes(5),
        memorySize: 1024,
      }
    );
    props.reportTable.grantReadWriteData(this.generateSectionLambda);
    // --------------------------------------------------------- //
    // Template Lambdas

    this.createTemplateLambda = new lambda.Function(
      this,
      "CreateTemplateLambda",
      {
        code: lambda.Code.fromAsset(
          path.join(__dirname, "../../bin/lambdas/create-template")
        ),
        handler: "main",
        runtime: lambda.Runtime.PROVIDED_AL2023,
        environment: {
          TEMPLATE_TABLE: props.templateTable.tableName,
        },
        memorySize: 1024,
      }
    );
    props.templateTable.grantReadWriteData(this.createTemplateLambda);

    this.getTemplateByIDLambda = new lambda.Function(
      this,
      "GetTemplateByIDLambda",
      {
        code: lambda.Code.fromAsset(
          path.join(__dirname, "../../bin/lambdas/get-template-by-id")
        ),
        handler: "main",
        runtime: lambda.Runtime.PROVIDED_AL2023,
        environment: {
          TEMPLATE_TABLE: props.templateTable.tableName,
        },
        memorySize: 1024,
      }
    );
    props.templateTable.grantReadData(this.getTemplateByIDLambda);

    this.getAllTemplatesLambda = new lambda.Function(
      this,
      "GetAllTemplatesLambda",
      {
        code: lambda.Code.fromAsset(
          path.join(__dirname, "../../bin/lambdas/get-all-templates")
        ),
        handler: "main",
        runtime: lambda.Runtime.PROVIDED_AL2023,
        environment: {
          TEMPLATE_TABLE: props.templateTable.tableName,
        },
        memorySize: 1024,
      }
    );
    props.templateTable.grantReadData(this.getAllTemplatesLambda);

    // --------------------------------------------------------- //
    // Shared Lambdas

    this.addPartLambda = new lambda.Function(this, "AddPartLambda", {
      code: lambda.Code.fromAsset(
        path.join(__dirname, "../../bin/lambdas/add-part")
      ),
      handler: "main",
      runtime: lambda.Runtime.PROVIDED_AL2023,
      environment: {
        REPORT_TABLE: props.reportTable.tableName,
        TEMPLATE_TABLE: props.templateTable.tableName,
      },
      memorySize: 1024,
    });
    props.reportTable.grantReadWriteData(this.addPartLambda);
    props.templateTable.grantReadWriteData(this.addPartLambda);

    this.addSectionLambda = new lambda.Function(this, "AddSectionLambda", {
      code: lambda.Code.fromAsset(
        path.join(__dirname, "../../bin/lambdas/add-section")
      ),
      handler: "main",
      runtime: lambda.Runtime.PROVIDED_AL2023,
      environment: {
        REPORT_TABLE: props.reportTable.tableName,
        TEMPLATE_TABLE: props.templateTable.tableName,
      },
      memorySize: 1024,
    });
    props.reportTable.grantReadWriteData(this.addSectionLambda);
    props.templateTable.grantReadWriteData(this.addSectionLambda);

    this.updatePartLambda = new lambda.Function(this, "UpdatePartLambda", {
      code: lambda.Code.fromAsset(
        path.join(__dirname, "../../bin/lambdas/update-part")
      ),
      handler: "main",
      runtime: lambda.Runtime.PROVIDED_AL2023,
      environment: {
        REPORT_TABLE: props.reportTable.tableName,
        TEMPLATE_TABLE: props.templateTable.tableName,
      },
      memorySize: 1024,
    });
    props.reportTable.grantReadWriteData(this.updatePartLambda);
    props.templateTable.grantReadWriteData(this.updatePartLambda);

    this.updateSectionLambda = new lambda.Function(
      this,
      "UpdateSectionLambda",
      {
        code: lambda.Code.fromAsset(
          path.join(__dirname, "../../bin/lambdas/update-section")
        ),
        handler: "main",
        runtime: lambda.Runtime.PROVIDED_AL2023,
        environment: {
          REPORT_TABLE: props.reportTable.tableName,
          TEMPLATE_TABLE: props.templateTable.tableName,
        },
        memorySize: 1024,
      }
    );
    props.reportTable.grantReadWriteData(this.updateSectionLambda);
    props.templateTable.grantReadWriteData(this.updateSectionLambda);

    this.updateItemTitleLambda = new lambda.Function(
      this,
      "UpdateItemTitleLambda",
      {
        code: lambda.Code.fromAsset(
          path.join(__dirname, "../../bin/lambdas/update-item-title")
        ),
        handler: "main",
        runtime: lambda.Runtime.PROVIDED_AL2023,
        environment: {
          REPORT_TABLE: props.reportTable.tableName,
          TEMPLATE_TABLE: props.templateTable.tableName,
        },
        memorySize: 1024,
      }
    );
    props.reportTable.grantReadWriteData(this.updateItemTitleLambda);
    props.templateTable.grantReadWriteData(this.updateItemTitleLambda);
    // --------------------------------------------------------- //
  }
}
