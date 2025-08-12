# Serverless Audit Log Ingestion Pipeline

## Overview

<p align="center">
  <img src="./infra.svg" alt="Required User Permissions" />
</p>

This project leverages AWS Serverless Application Model (SAM) to build a serverless audit log ingestion pipeline using AWS services like API Gateway, Lambda, Kinesis Data Firehose, S3, Glue, and Athena. It enables reliable, scalable, and cost-effective collection, transformation, and querying of audit logs.

For a detailed explanation of the implementation, check out my blog post: [Building a Serverless Audit Log Pipeline with AWS SAM — Easy, Fast & Fun](https://jfolgado.com/posts/audittrail/).

Features
- **Serverless Deployment**: Easily deploy the entire pipeline using AWS SAM infrastructure-as-code.
- **Dynamic Partitioning**: Firehose delivers logs to S3 with dynamic partitions for efficient querying.
- **Real-Time Transformation**: Lambda function enriches and partitions logs on the fly.
- **Schema Management**: AWS Glue Crawler automates schema detection for Athena queries.

## Prerequisites
Make sure you have these installed and configured:
- AWS CLI
- AWS SAM CLI
- Go 1.22+ (or relevant version used for the Lambda function)

## Deployment Steps

```
git clone git@github.com:Flgado/Aws-audit-data-pipeline.git

cd /Aws-audit-data-pipeline

sam build -u

sam deploy -g
``` 
Follow the prompts to configure your AWS credentials, region, and stack settings.

## Usage
Once deployed, your pipeline will:
- Accept audit log JSON via API Gateway POST requests.
- Send logs through Lambda and Kinesis Data Firehose to S3.
- Organize logs into dynamic partitions by audit type and date.
- Use AWS Glue to catalog data schemas automatically.
- Enable Athena queries to filter and analyze audit logs efficiently.

You can invoke the API Gateway endpoint or send audit logs directly to Firehose (adjust as needed).