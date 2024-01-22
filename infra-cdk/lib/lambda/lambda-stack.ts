import * as cdk from 'aws-cdk-lib'
import { type Construct } from 'constructs'
import * as lambda from 'aws-cdk-lib/aws-lambda'
import type * as dynamodb from 'aws-cdk-lib/aws-dynamodb'
import path = require('path')
import * as fs from 'fs'

interface LambdasStackProps extends cdk.StackProps {
  reportTable: dynamodb.Table
}

export class LambdasStack extends cdk.Stack {
  public readonly getReportByIDLambda: lambda.IFunction
  public readonly getAllReportsLambda: lambda.IFunction
  public readonly createReportLambda: lambda.IFunction
  public readonly addPartLambda: lambda.IFunction
  public readonly addSectionLambda: lambda.IFunction
  public readonly generateSectionLambda: lambda.IFunction
  public readonly getAllReportTypesLambda: lambda.IFunction

  constructor (scope: Construct, id: string, props: LambdasStackProps) {
    super(scope, id, props)

    this.createReportLambda = new lambda.Function(this, 'CreateReportLambda', {
      code: lambda.Code.fromAsset(
        path.join(__dirname, '../../bin/lambdas/post-create-report')
      ),
      handler: 'main',
      runtime: lambda.Runtime.PROVIDED_AL2023,
      environment: {
        REPORT_TABLE: props.reportTable.tableName
      }
    })

    this.getReportByIDLambda = new lambda.Function(
      this,
      'GetReportByIDLambda',
      {
        code: lambda.Code.fromAsset(
          path.join(__dirname, '../../bin/lambdas/get-report-by-id')
        ),
        handler: 'main',
        runtime: lambda.Runtime.PROVIDED_AL2023,
        environment: {
          REPORT_TABLE: props.reportTable.tableName
        }
      }
    )

    this.getAllReportsLambda = new lambda.Function(
      this,
      'GetAllReportsLambda',
      {
        code: lambda.Code.fromAsset(
          path.join(__dirname, '../../bin/lambdas/get-all-reports')
        ),
        handler: 'main',
        runtime: lambda.Runtime.PROVIDED_AL2023,
        environment: {
          REPORT_TABLE: props.reportTable.tableName
        }
      }
    )

    this.getAllReportTypesLambda = new lambda.Function(
      this,
      'GetAllReportTypesLambda',
      {
        code: lambda.Code.fromAsset(
          path.join(__dirname, '../../bin/lambdas/get-all-report-types')
        ),
        handler: 'main',
        runtime: lambda.Runtime.PROVIDED_AL2023,
        environment: {
          REPORT_TABLE: props.reportTable.tableName
        }
      }
    )

    this.addPartLambda = new lambda.Function(this, 'AddPartLambda', {
      code: lambda.Code.fromAsset(
        path.join(__dirname, '../../bin/lambdas/post-add-part')
      ),
      handler: 'main',
      runtime: lambda.Runtime.PROVIDED_AL2023,
      environment: {
        REPORT_TABLE: props.reportTable.tableName
      }
    })

    this.addSectionLambda = new lambda.Function(this, 'AddSectionLambda', {
      code: lambda.Code.fromAsset(
        path.join(__dirname, '../../bin/lambdas/post-add-section')
      ),
      handler: 'main',
      runtime: lambda.Runtime.PROVIDED_AL2023,
      environment: {
        REPORT_TABLE: props.reportTable.tableName
      }
    })

    // Set openAI key as env variable
    const openAIKeyPath = path.join(__dirname, '../../../keys/openai-key.txt')
    const openAIKey = fs.readFileSync(openAIKeyPath, {
      encoding: 'utf8',
      flag: 'r'
    })

    this.generateSectionLambda = new lambda.Function(
      this,
      'GenerateSectionLambda',
      {
        code: lambda.Code.fromAsset(
          path.join(__dirname, '../../bin/lambdas/post-generate-section')
        ),
        handler: 'main',
        runtime: lambda.Runtime.PROVIDED_AL2023,
        environment: {
          REPORT_TABLE: props.reportTable.tableName,
          OPENAI_API_KEY: openAIKey
        },
        timeout: cdk.Duration.minutes(5)
      }
    )

    props.reportTable.grantReadData(this.getReportByIDLambda)
    props.reportTable.grantReadData(this.getAllReportsLambda)
    props.reportTable.grantWriteData(this.createReportLambda)
    props.reportTable.grantReadWriteData(this.addPartLambda)
    props.reportTable.grantReadWriteData(this.addSectionLambda)
    props.reportTable.grantReadWriteData(this.generateSectionLambda)
  }
}
