import boto3

bucket = 'blakecaldwell.garage'
bucket_key = 'status'

def lambda_handler(event, context):
    # make sure event looks good
    if 'status' not in event:
        return 'bad request'

    # validate status
    if event['status'] not in ['open', 'closed']:
        return 'bad request'

    # store the object
    s3 = boto3.client('s3')
    s3.put_object(
        Bucket=bucket,
        Key=bucket_key,
        Body=event['status'])

    return 'success'
