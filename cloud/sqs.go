package cloud

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

var sqsc *sqs.SQS

func init() {
	sqsc = sqs.New(session.New())
}

// SqsSendMessage sends the message to the queue and returns the
// message id and an error if anything goes wrong. If the len of the payload
// is zero, this function panics.
func SqsSendMessage(ctx aws.Context, queue, payload string) (string, error) {

	if len(payload) == 0 {
		panic(fmt.Sprintf("Empty msg payload to SqsSendMessage for queue: %s", queue))
	}

	res, err := sqsc.SendMessageWithContext(
		ctx,
		&sqs.SendMessageInput{
			MessageBody: aws.String(payload),
			QueueUrl:    aws.String(queue),
		},
	)

	if err != nil {
		return "", fmt.Errorf("Failed to send message to queue: %s", queue)
	}

	return aws.StringValue(res.MessageId), nil
}
