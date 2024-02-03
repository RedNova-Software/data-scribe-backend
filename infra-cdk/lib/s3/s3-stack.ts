import * as cdk from "aws-cdk-lib";
import * as cognito from "aws-cdk-lib/aws-cognito";
import { Construct } from "constructs";
import * as lambda from "aws-cdk-lib/aws-lambda";
import * as iam from "aws-cdk-lib/aws-iam";
import * as s3 from "aws-cdk-lib/aws-s3";
import * as s3n from "aws-cdk-lib/aws-s3-notifications";
import type * as dynamodb from "aws-cdk-lib/aws-dynamodb";
import path = require("path");

interface S3BucketStackProps extends cdk.StackProps {
  reportTable: dynamodb.Table;
  operationTable: dynamodb.Table;
}

export class S3BucketStack extends cdk.Stack {
  public readonly csvBucket: s3.Bucket;
  public readonly columnDataBucket: s3.Bucket;

  constructor(scope: Construct, id: string, props: S3BucketStackProps) {
    super(scope, id, props);

    this.csvBucket = new s3.Bucket(this, "CsvBucket", {
      bucketName: "scribe-csv-bucket",
      publicReadAccess: false,
      encryption: s3.BucketEncryption.S3_MANAGED,
      blockPublicAccess: s3.BlockPublicAccess.BLOCK_ALL,
    });

    this.columnDataBucket = new s3.Bucket(this, "ColumnDataBucket", {
      bucketName: "scribe-column-data-bucket",
      publicReadAccess: false,
      encryption: s3.BucketEncryption.S3_MANAGED,
      blockPublicAccess: s3.BlockPublicAccess.BLOCK_ALL,
    });

    // Add CORS rule. This is needed for pre-signed urls
    // that we use to upload csvs
    this.csvBucket.addCorsRule({
      allowedMethods: [
        s3.HttpMethods.GET,
        s3.HttpMethods.PUT,
        s3.HttpMethods.POST,
        s3.HttpMethods.DELETE,
      ],
      allowedOrigins: ["*"],
      allowedHeaders: ["Content-Type"],
    });

    // A lambda to process a new CSV file when its uploaded
    // It will set the columns and unique columns values in a Report
    const readCsvColumnsLambda = new lambda.Function(
      this,
      "ReadCsvColumnsLambda",
      {
        code: lambda.Code.fromAsset(
          path.join(__dirname, "../../bin/lambdas/read-csv-columns")
        ),
        handler: "main",
        runtime: lambda.Runtime.PROVIDED_AL2023,
        environment: {
          REPORT_TABLE: props.reportTable.tableName,
          OPERATION_TABLE: props.operationTable.tableName,
          CSV_BUCKET_NAME: this.csvBucket.bucketName,
          COLUMN_DATA_BUCKET_NAME: this.columnDataBucket.bucketName,
        },
        memorySize: 1024,
      }
    );
    props.reportTable.grantReadWriteData(readCsvColumnsLambda);
    props.operationTable.grantWriteData(readCsvColumnsLambda);
    this.csvBucket.grantRead(readCsvColumnsLambda);
    this.columnDataBucket.grantReadWrite(readCsvColumnsLambda);

    this.csvBucket.addEventNotification(
      s3.EventType.OBJECT_CREATED,
      new s3n.LambdaDestination(readCsvColumnsLambda)
    );

    // Optional: Output the bucket name
    new cdk.CfnOutput(this, "BucketName", { value: this.csvBucket.bucketName });
  }
}
