import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as lambda from "aws-cdk-lib/aws-lambda"
import * as apigateway from "aws-cdk-lib/aws-apigateway";
import path = require('path');

export class DataScribeBackendStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    const myFunction = new lambda.Function(this, "MyLambda", {
      code: lambda.Code.fromAsset(path.join(__dirname, "../bin/lambdas/test-endpoint")),
      handler: "main",
      runtime: lambda.Runtime.PROVIDED_AL2023,
    })

    const myNewFunction = new lambda.Function(this, "NewLambda", {
      code: lambda.Code.fromAsset(path.join(__dirname, "../bin/lambdas/new-endpoint")),
      handler: "main",
      runtime: lambda.Runtime.PROVIDED_AL2023,
    })

    const gateway = new apigateway.RestApi(this, "myGateway", {
      defaultCorsPreflightOptions: {
        allowOrigins: ["*"],
        allowMethods: ["GET", "POST", "OPTIONS", "DELETE", "PUT"],
      }
    })


    const integration = new apigateway.LambdaIntegration(myFunction)
    const newIntegration = new apigateway.LambdaIntegration(myNewFunction)


    const testEndpoint = gateway.root.addResource("test")
    const newEndpoint = gateway.root.addResource("new")


    testEndpoint.addMethod("GET", integration)
    newEndpoint.addMethod("GET", newIntegration)
  }


}
