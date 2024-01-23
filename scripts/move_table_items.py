import boto3

dynamodb = boto3.resource('dynamodb')

source_table = dynamodb.Table('DataScribeBackendStack-ReportTable270236C0-IXUO3HMP50Z0')
target_table = dynamodb.Table('DynamoDBStack-ReportTable270236C0-1FGP28WTJ6YO')

# Scan the source table
response = source_table.scan()

for item in response['Items']:
    # Write each item to the target table
    target_table.put_item(Item=item)

print("Move Complete")