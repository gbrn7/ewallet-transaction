package external

import (
	"context"
	"ewallet-transaction/constants"
	"ewallet-transaction/external/proto/notification"
	"ewallet-transaction/helpers"
	"fmt"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func (*External) SendNotification(ctx context.Context, recipient string, templateName string, placeHolder map[string]string) error {

	conn, err := grpc.NewClient(helpers.GetEnv("NOTIFICATION_GRPC_HOST", "7003"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return errors.Wrap(err, "failed to dial grpc")
	}

	defer conn.Close()

	client := notification.NewNotificationServiceClient(conn)
	request := &notification.SendNotificationRequest{
		Recipient:    recipient,
		TemplateName: templateName,
		Placeholders: placeHolder,
	}

	resp, err := client.SendNotification(ctx, request)
	if err != nil {
		return err
	}

	if resp.Message != constants.SuccessMessage {
		return fmt.Errorf("get response error from notification: %s", resp.Message)
	}

	return nil
}
