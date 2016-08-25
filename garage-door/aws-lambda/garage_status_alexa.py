import boto3
import traceback
import sys

bucket = 'blakecaldwell.garage'
bucket_key = 'status'

def lambda_handler(event, context):
    try:
        if 'request' not in event or 'type' not in event['request']:
            return 'Invalid request'

        if event['request']['type'] not in ['LaunchRequest', 'IntentRequest']:
            return 'Invalid request'

        # Not using sessions for now
        sessionAttributes = {}

        if event['request']['type'] == "LaunchRequest":
            speechlet = onLaunch(event['request'])
            response = buildResponse(sessionAttributes, speechlet)
        elif event['request']['type'] == "IntentRequest":
            speechlet = onIntent(event['request'])
            response = buildResponse(sessionAttributes, speechlet)

        # Return a response for speech output
        return response

    except Exception, err:
        traceback.print_exc(file=sys.stdout)
        return 'an error occurred'

def status():
    s3 = boto3.client('s3')
    status_obj = s3.get_object(
        Bucket=bucket,
        Key=bucket_key)
    return str(status_obj['Body'].read())

# Called when the user launches the skill without specifying what they want.
def onLaunch(launchRequest):
    # Dispatch to your skill's launch.
    getWelcomeResponse()

# Called when the user specifies an intent for this skill.
def onIntent(intentRequest):
    intent = intentRequest['intent']
    intentName = intentRequest['intent']['name']

    # Dispatch to your skill's intent handlers
    if intentName == "StateIntent":
        return stateResponse(intent)
    elif intentName == "HelpIntent":
        return getWelcomeResponse()
    else:
        print "Invalid Intent (" + intentName + ")"
        raise

def getWelcomeResponse():
    sessionAttributes = {}
    cardTitle = "Welcome"
    speechOutput = "You can check your garage door's status by saying, ask my garage door if it's open."

    # If the user either does not reply to the welcome message or says something that is not
    # understood, they will be prompted again with this text.
    repromptText = "Ask if your garage door is open by saying, ask my garage door if it's open."
    shouldEndSession = True

    return (buildSpeechletResponse(cardTitle, speechOutput, repromptText, shouldEndSession))

def stateResponse(intent):
    # Ask my garage door if it's {open|closed}
    #     "intent": {
    #       "name": "StateIntent",
    #       "slots": {
    #         "doorstate": {
    #           "name": "doorstate",
    #           "value": "closed"
    #         }
    #       }
    #     }
    doorstate = status()

    if (intent['slots']['doorstate']['value'] == "open") or (intent['slots']['doorstate']['value'] == "up"):
        if doorstate == "open":
            speechOutput = "Yes, your garage door is open"
            cardTitle = "Yes, your garage door is open"
        elif doorstate == "closed":
            speechOutput = "No, your garage door is closed"
            cardTitle = "No, your garage door is closed"
        else:
            speechOutput = "Your garage door is " + doorstate
            cardTitle = "Your garage door is " + doorstate

    elif (intent['slots']['doorstate']['value'] == "closed") or (intent['slots']['doorstate']['value'] == "shut") or (intent['slots']['doorstate']['value'] == "down"):
        if doorstate == "closed":
            speechOutput = "Yes, your garage door is closed"
            cardTitle = "Yes, your garage door is closed"
        elif doorstate == "open":
            speechOutput = "No, your garage door is open"
            cardTitle = "No, your garage door is open"
        else:
            speechOutput = "Your garage door is " + doorstate
            cardTitle = "Your garage door is " + doorstate

    else:
        speechOutput = "I didn't understand that. You can say ask my garage door if it's open."
        cardTitle = "Try again"

    repromptText = "I didn't understand that. You can say ask the garage door if it's open"
    shouldEndSession = True

    return(buildSpeechletResponse(cardTitle, speechOutput, repromptText, shouldEndSession))

# --------------- Helpers that build all of the responses -----------------------

def buildSpeechletResponse(title, output, repromptText, shouldEndSession):
    return ({
        "outputSpeech": {
            "type": "PlainText",
            "text": output
        },
        "card": {
            "type": "Simple",
            "title": "Garage door - " + title,
            "content": "Garage door - " + output
        },
        "reprompt": {
            "outputSpeech": {
                "type": "PlainText",
                "text": repromptText
            }
        },
        "shouldEndSession": shouldEndSession
    })

def buildResponse(sessionAttributes, speechletResponse):
    return ({
        "version": "1.0",
        "sessionAttributes": sessionAttributes,
        "response": speechletResponse
    })
