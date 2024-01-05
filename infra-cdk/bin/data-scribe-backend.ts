#!/usr/bin/env node
import 'source-map-support/register';
import * as cdk from 'aws-cdk-lib';
import { DataScribeBackendStack } from '../lib/data-scribe-backend-stack';

const app = new cdk.App();
new DataScribeBackendStack(app, 'DataScribeBackendStack', {
  /* If you don't specify 'env', this stack will be environment-agnostic.
   * Account/Region-dependent features and context lookups will not work,
   * but a single synthesized template can be deployed anywhere. */

  /* Uncomment the next line to specialize this stack for the AWS Account
   * and Region that are implied by the current CLI configuration. */
  // env: { account: process.env.CDK_DEFAULT_ACCOUNT, region: process.env.CDK_DEFAULT_REGION },

  env: { account: '905418134223', region: 'us-east-2' },

  /* For more information, see https://docs.aws.amazon.com/cdk/latest/guide/environments.html */
});