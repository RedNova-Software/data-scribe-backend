# Serverless Golang API for AWS Lambda

## Overview

This repository contains a serverless Golang API designed to run on AWS Lambda. The structure is optimized for minimal binary sizes so that each Lambda function includes only what is necessary for its operation, ensuring efficient execution when triggered by API Gateway.

## Project Structure

- `./api`: The main directory for all API-related code.
  - `/lambdas`: Each subdirectory represents a separate Lambda function. To create a new Lambda endpoint:
    - Create a new folder under `/lambdas`.
    - Add a `handler.go` file in this folder, which will be the entry point for the Lambda function.
  - `/shared`: Contains shared server logic and Go modules that can be reused across different Lambda functions. This is where you can place common utilities, middleware, data access layers, etc.

## Creating a New Lambda Function

To set up a new Lambda function:

1. **Create a Lambda Handler:**

   - Navigate to the `./api/lambdas` directory.
   - Create a new folder named after your Lambda function.
   - Inside this new folder, create a `handler.go` file that will serve as the entry point for your Lambda.

2. **Import Shared Code:**

   - Utilize the shared modules by importing necessary code from the `./api/shared` directory into your `handler.go` file.

3. **Update CDK Stack:**
   - Add the necessary infrastructure code to the `./infra-cdk/lib/data-scribe-backend-stack.ts` file to define the AWS resources required for your new Lambda function.

## Deployment

To deploy your Lambda functions along with the infrastructure to AWS, simply execute the following command:

```bash
npm run deploy
```

To build the lambdas into binaries simply:

```bash
npm run build
```

To hotswap the lambdas into the cloud (only updating Go code for iterating on an endpoint) simply:

```bash
npm run hotswap
```

## Design Philosophy

The architecture is crafted to ensure that each Lambda function acts as an independent microservice, containing only the code that it needs to perform its job. This results in faster start times and more efficient resource utilization, as unnecessary dependencies and bloat are eliminated. When API Gateway invokes a Lambda function, it starts up with the minimal set of binaries required for that specific endpoint, adhering to the principles of lean software and on-demand scalability.

## Contributing

When contributing to this repository, please ensure that any common logic that could be utilized across multiple Lambda functions is placed within the ./api/shared directory. This helps in maintaining a DRY codebase and simplifies the management of shared dependencies.
