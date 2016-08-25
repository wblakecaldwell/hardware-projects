import boto3
import traceback
import sys

bucket = 'blakecaldwell.garage'
bucket_key = 'status'

def lambda_handler(event, context):
    try:
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

    except Exception, err:
        traceback.print_exc(file=sys.stdout)
        return 'an error occurred'
