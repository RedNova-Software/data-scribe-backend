import * as cdk from "aws-cdk-lib";
import * as cognito from "aws-cdk-lib/aws-cognito";
import { Construct } from "constructs";
import * as lambda from "aws-cdk-lib/aws-lambda";
import * as iam from "aws-cdk-lib/aws-iam";
import path = require("path");
import { buildLambda } from "./helper/build-lambda";

export class CognitoUserPoolStack extends cdk.Stack {
  public readonly userPool: cognito.UserPool;

  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // Build the post confirmation lambda
    buildLambda();

    // Define the lambda function
    const userConfirmationLambda = new lambda.Function(
      this,
      "UserConfirmationLambda",
      {
        runtime: lambda.Runtime.NODEJS_20_X,
        handler: "index.handler",
        code: lambda.Code.fromAsset(
          path.join(__dirname, "./lambda/build/lambda.zip")
        ),
      }
    );

    // Define the policy statement that allows 'AdminDisableUser' action
    const policyStatement = new iam.PolicyStatement({
      actions: ["cognito-idp:AdminDisableUser"],
      effect: iam.Effect.ALLOW,
      resources: [
        // This must be updated if redeploying the Cognito Stack as the arn of the
        // userpool will change
        "arn:aws:cognito-idp:us-east-2:905418134223:userpool/us-east-2_w8Ssq0MxW",
      ],
    });

    userConfirmationLambda.role?.attachInlinePolicy(
      new iam.Policy(this, "AllowAdminDisableUser", {
        statements: [policyStatement],
      })
    );

    // Create the Cognito user pool with Lambda trigger
    this.userPool = new cognito.UserPool(this, "DataScribeUserPool", {
      selfSignUpEnabled: true,
      userVerification: {
        emailStyle: cognito.VerificationEmailStyle.CODE,
        emailSubject: "Data Scribe Sign Up Confirmation",
      },
      standardAttributes: {
        nickname: {
          mutable: false,
          required: true,
        },
      },
      signInAliases: {
        email: true,
      },
      lambdaTriggers: {
        postConfirmation: userConfirmationLambda,
      },
      passwordPolicy: {
        minLength: 8,
        requireDigits: true,
        requireLowercase: true,
        requireUppercase: true,
        requireSymbols: true,
      },
      mfa: cognito.Mfa.OPTIONAL,
      mfaSecondFactor: {
        sms: true,
        otp: false,
      },
    });

    const scribeClient = this.userPool.addClient("DataScribeClient", {
      generateSecret: false,
      authFlows: {
        userPassword: true,
        userSrp: true,
      },
    });

    // Output User Pool ID
    new cdk.CfnOutput(this, "UserPoolId", {
      value: this.userPool.userPoolId,
    });

    new cdk.CfnOutput(this, "UserPoolClientId", {
      value: scribeClient.userPoolClientId,
    });
  }
}
