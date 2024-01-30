import * as cdk from "aws-cdk-lib";
import * as cognito from "aws-cdk-lib/aws-cognito";
import { Construct } from "constructs";
import * as lambda from "aws-cdk-lib/aws-lambda";
import * as iam from "aws-cdk-lib/aws-iam";
import * as s3 from "aws-cdk-lib/aws-s3";
import path = require("path");

export class S3BucketStack extends cdk.Stack {
  public readonly csvBucket: s3.Bucket;

  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // Define the S3 bucket
    this.csvBucket = new s3.Bucket(this, "CsvBucket", {
      bucketName: "scribe-csv-bucket", // Replace with your desired bucket name
      publicReadAccess: false,
      encryption: s3.BucketEncryption.S3_MANAGED,
      blockPublicAccess: s3.BlockPublicAccess.BLOCK_ALL,
    });

    // Optional: Output the bucket name
    new cdk.CfnOutput(this, "BucketName", { value: this.csvBucket.bucketName });
  }
}
