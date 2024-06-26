import * as cdk from "aws-cdk-lib";
import { type Construct } from "constructs";
import * as lambda from "aws-cdk-lib/aws-lambda";
import * as cognito from "aws-cdk-lib/aws-cognito";
import * as s3 from "aws-cdk-lib/aws-s3";
import type * as dynamodb from "aws-cdk-lib/aws-dynamodb";
import path = require("path");
import * as fs from "fs";

interface LambdasStackProps extends cdk.StackProps {
  reportTable: dynamodb.Table;
  templateTable: dynamodb.Table;
  operationsTable: dynamodb.Table;
  userPool: cognito.UserPool;
  readonly csvBucket: s3.Bucket;
  readonly columnDataBucket: s3.Bucket;
}

export class LambdasStack extends cdk.Stack {
  // Report Lambdas
  public readonly getReportByIDLambda: lambda.IFunction;
  public readonly getAllReportsLambda: lambda.IFunction;
  public readonly createReportLambda: lambda.IFunction;
  public readonly generateSectionLambda: lambda.IFunction;
  public readonly getAllReportTypesLambda: lambda.IFunction;
  public readonly uploadCSVLambda: lambda.IFunction;
  public readonly getCSVUniqueColumnsMapLambda: lambda.IFunction;
  public readonly setSectionResponsesLambda: lambda.IFunction;

  // Template Lambas
  public readonly getTemplateByIDLambda: lambda.IFunction;
  public readonly getAllTemplatesLambda: lambda.IFunction;
  public readonly createTemplateLambda: lambda.IFunction;

  // Shared Lambdas
  public readonly addPartLambda: lambda.IFunction;
  public readonly deletePartLambda: lambda.IFunction;
  public readonly addSectionLambda: lambda.IFunction;
  public readonly deleteSectionLambda: lambda.IFunction;
  public readonly updatePartLambda: lambda.IFunction;
  public readonly updateSectionLambda: lambda.IFunction;
  public readonly updateItemTitleLambda: lambda.IFunction;
  public readonly shareItemLambda: lambda.IFunction;
  public readonly convertItemLambda: lambda.IFunction;
  public readonly deleteItemLambda: lambda.IFunction;
  public readonly restoreItemLambda: lambda.IFunction;
  public readonly updateItemGlobalQuestionsLambda: lambda.IFunction;

  // User Lambdas
  public readonly getUserIDLambda: lambda.IFunction;
  public readonly getAllUsersLambda: lambda.IFunction;

  // Operation lambdas
  public readonly getOperationStatusLambda: lambda.IFunction;

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
        USER_POOL_ID: props.userPool.userPoolId,
      },
      memorySize: 1024,
    });
    props.reportTable.grantWriteData(this.createReportLambda);
    props.userPool.grant(this.createReportLambda, "cognito-idp:AdminGetUser");

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
          USER_POOL_ID: props.userPool.userPoolId,
        },
        memorySize: 1024,
      }
    );
    props.reportTable.grantReadData(this.getReportByIDLambda);
    props.userPool.grant(this.getReportByIDLambda, "cognito-idp:AdminGetUser");

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
          USER_POOL_ID: props.userPool.userPoolId,
        },
        memorySize: 1024,
      }
    );
    props.reportTable.grantReadData(this.getAllReportsLambda);
    props.userPool.grant(this.getAllReportsLambda, "cognito-idp:AdminGetUser");

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
          OPERATION_TABLE: props.operationsTable.tableName,
          CSV_BUCKET_NAME: props.csvBucket.bucketName,
          OPENAI_API_KEY: openAIKey,
        },
        timeout: cdk.Duration.minutes(2.5),
        memorySize: 2048,
      }
    );
    props.reportTable.grantReadWriteData(this.generateSectionLambda);
    props.userPool.grant(
      this.generateSectionLambda,
      "cognito-idp:AdminGetUser"
    );
    props.csvBucket.grantReadWrite(this.generateSectionLambda);
    props.operationsTable.grantReadWriteData(this.generateSectionLambda);

    this.uploadCSVLambda = new lambda.Function(this, "UploadCSVLambda", {
      code: lambda.Code.fromAsset(
        path.join(__dirname, "../../bin/lambdas/upload-csv")
      ),
      handler: "main",
      runtime: lambda.Runtime.PROVIDED_AL2023,
      memorySize: 1024,
      environment: {
        REPORT_TABLE: props.reportTable.tableName,
        OPERATION_TABLE: props.operationsTable.tableName,
        CSV_BUCKET_NAME: props.csvBucket.bucketName,
      },
      timeout: cdk.Duration.seconds(30),
    });
    props.reportTable.grantReadWriteData(this.uploadCSVLambda);
    props.csvBucket.grantReadWrite(this.uploadCSVLambda);
    props.operationsTable.grantReadWriteData(this.uploadCSVLambda);

    this.getCSVUniqueColumnsMapLambda = new lambda.Function(
      this,
      "GetCSVUniqueColumnsMapLambda",
      {
        code: lambda.Code.fromAsset(
          path.join(__dirname, "../../bin/lambdas/get-csv-unique-columns-map")
        ),
        handler: "main",
        runtime: lambda.Runtime.PROVIDED_AL2023,
        memorySize: 2048,
        environment: {
          REPORT_TABLE: props.reportTable.tableName,
          OPERATION_TABLE: props.operationsTable.tableName,
          COLUMN_DATA_BUCKET_NAME: props.columnDataBucket.bucketName,
        },
        timeout: cdk.Duration.seconds(30),
      }
    );
    props.reportTable.grantReadWriteData(this.getCSVUniqueColumnsMapLambda);
    props.columnDataBucket.grantReadWrite(this.getCSVUniqueColumnsMapLambda);
    props.operationsTable.grantReadWriteData(this.getCSVUniqueColumnsMapLambda);

    this.setSectionResponsesLambda = new lambda.Function(
      this,
      "SetSectionResponsesLambda",
      {
        code: lambda.Code.fromAsset(
          path.join(__dirname, "../../bin/lambdas/set-section-responses")
        ),
        handler: "main",
        runtime: lambda.Runtime.PROVIDED_AL2023,
        memorySize: 1024,
        environment: {
          REPORT_TABLE: props.reportTable.tableName,
        },
        timeout: cdk.Duration.seconds(30),
      }
    );
    props.reportTable.grantReadWriteData(this.setSectionResponsesLambda);

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
          USER_POOL_ID: props.userPool.userPoolId,
        },
        memorySize: 1024,
      }
    );
    props.templateTable.grantReadWriteData(this.createTemplateLambda);
    props.userPool.grant(this.createTemplateLambda, "cognito-idp:AdminGetUser");

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
          USER_POOL_ID: props.userPool.userPoolId,
        },
        memorySize: 1024,
      }
    );
    props.templateTable.grantReadData(this.getTemplateByIDLambda);
    props.userPool.grant(
      this.getTemplateByIDLambda,
      "cognito-idp:AdminGetUser"
    );

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
          USER_POOL_ID: props.userPool.userPoolId,
        },
        memorySize: 1024,
      }
    );
    props.templateTable.grantReadData(this.getAllTemplatesLambda);
    props.userPool.grant(
      this.getAllTemplatesLambda,
      "cognito-idp:AdminGetUser"
    );

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
    props.userPool.grant(this.addPartLambda, "cognito-idp:AdminGetUser");

    this.deletePartLambda = new lambda.Function(this, "DeletePartLambda", {
      code: lambda.Code.fromAsset(
        path.join(__dirname, "../../bin/lambdas/delete-part")
      ),
      handler: "main",
      runtime: lambda.Runtime.PROVIDED_AL2023,
      environment: {
        REPORT_TABLE: props.reportTable.tableName,
        TEMPLATE_TABLE: props.templateTable.tableName,
      },
      memorySize: 1024,
    });
    props.reportTable.grantReadWriteData(this.deletePartLambda);
    props.templateTable.grantReadWriteData(this.deletePartLambda);
    props.userPool.grant(this.deletePartLambda, "cognito-idp:AdminGetUser");

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
    props.userPool.grant(this.addSectionLambda, "cognito-idp:AdminGetUser");

    this.deleteSectionLambda = new lambda.Function(
      this,
      "DeleteSectionLambda",
      {
        code: lambda.Code.fromAsset(
          path.join(__dirname, "../../bin/lambdas/delete-section")
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
    props.reportTable.grantReadWriteData(this.deleteSectionLambda);
    props.templateTable.grantReadWriteData(this.deleteSectionLambda);
    props.userPool.grant(this.deleteSectionLambda, "cognito-idp:AdminGetUser");

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
    props.userPool.grant(this.updatePartLambda, "cognito-idp:AdminGetUser");

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
    props.userPool.grant(this.updateSectionLambda, "cognito-idp:AdminGetUser");

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
    props.userPool.grant(
      this.updateItemTitleLambda,
      "cognito-idp:AdminGetUser"
    );

    this.shareItemLambda = new lambda.Function(this, "ShareItemLambda", {
      code: lambda.Code.fromAsset(
        path.join(__dirname, "../../bin/lambdas/share-item")
      ),
      handler: "main",
      runtime: lambda.Runtime.PROVIDED_AL2023,
      memorySize: 1024,
      environment: {
        REPORT_TABLE: props.reportTable.tableName,
        TEMPLATE_TABLE: props.templateTable.tableName,
      },
    });
    props.reportTable.grantReadWriteData(this.shareItemLambda);
    props.templateTable.grantReadWriteData(this.shareItemLambda);

    this.convertItemLambda = new lambda.Function(this, "ConvertItemLambda", {
      code: lambda.Code.fromAsset(
        path.join(__dirname, "../../bin/lambdas/convert-item")
      ),
      handler: "main",
      runtime: lambda.Runtime.PROVIDED_AL2023,
      memorySize: 1024,
      environment: {
        USER_POOL_ID: props.userPool.userPoolId,
        REPORT_TABLE: props.reportTable.tableName,
        TEMPLATE_TABLE: props.templateTable.tableName,
      },
    });
    props.reportTable.grantReadWriteData(this.convertItemLambda);
    props.templateTable.grantReadWriteData(this.convertItemLambda);

    this.deleteItemLambda = new lambda.Function(this, "DeleteItemLambda", {
      code: lambda.Code.fromAsset(
        path.join(__dirname, "../../bin/lambdas/delete-item")
      ),
      handler: "main",
      runtime: lambda.Runtime.PROVIDED_AL2023,
      environment: {
        REPORT_TABLE: props.reportTable.tableName,
        TEMPLATE_TABLE: props.templateTable.tableName,
      },
      memorySize: 1024,
    });
    props.reportTable.grantReadWriteData(this.deleteItemLambda);
    props.templateTable.grantReadWriteData(this.deleteItemLambda);

    this.restoreItemLambda = new lambda.Function(this, "RestoreItemLambda", {
      code: lambda.Code.fromAsset(
        path.join(__dirname, "../../bin/lambdas/restore-item")
      ),
      handler: "main",
      runtime: lambda.Runtime.PROVIDED_AL2023,
      environment: {
        REPORT_TABLE: props.reportTable.tableName,
        TEMPLATE_TABLE: props.templateTable.tableName,
      },
      memorySize: 1024,
    });
    props.reportTable.grantReadWriteData(this.restoreItemLambda);
    props.templateTable.grantReadWriteData(this.restoreItemLambda);

    this.updateItemGlobalQuestionsLambda = new lambda.Function(
      this,
      "UpdateItemGlobalQuestions",
      {
        code: lambda.Code.fromAsset(
          path.join(__dirname, "../../bin/lambdas/update-item-global-questions")
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
    props.reportTable.grantReadWriteData(this.updateItemGlobalQuestionsLambda);
    props.templateTable.grantReadWriteData(
      this.updateItemGlobalQuestionsLambda
    );

    // --------------------------------------------------------- //

    // User Lambdas

    this.getUserIDLambda = new lambda.Function(this, "GetUserIDLambda", {
      code: lambda.Code.fromAsset(
        path.join(__dirname, "../../bin/lambdas/get-user-id")
      ),
      handler: "main",
      runtime: lambda.Runtime.PROVIDED_AL2023,
      memorySize: 1024,
    });

    this.getAllUsersLambda = new lambda.Function(this, "GetAllUsersLambda", {
      code: lambda.Code.fromAsset(
        path.join(__dirname, "../../bin/lambdas/get-all-users")
      ),
      handler: "main",
      runtime: lambda.Runtime.PROVIDED_AL2023,
      memorySize: 1024,
      environment: {
        USER_POOL_ID: props.userPool.userPoolId,
        REPORT_TABLE: props.reportTable.tableName,
        TEMPLATE_TABLE: props.templateTable.tableName,
      },
    });
    props.userPool.grant(this.getAllUsersLambda, "cognito-idp:ListUsers");

    // --------------------------------------------------------- //

    // Operation Lambdas

    this.getOperationStatusLambda = new lambda.Function(
      this,
      "GetOperationStatusLambda",
      {
        code: lambda.Code.fromAsset(
          path.join(__dirname, "../../bin/lambdas/get-operation-status")
        ),
        handler: "main",
        runtime: lambda.Runtime.PROVIDED_AL2023,
        memorySize: 1024,
        environment: {
          OPERATION_TABLE: props.operationsTable.tableName,
        },
      }
    );
    props.operationsTable.grantReadData(this.getOperationStatusLambda);
  }
}
