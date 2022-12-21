package service

import (
	"context"

	"github.com/lixvyang/betxin/model"

	"github.com/fox-one/mixin-sdk-go"
)

func CreateMessage(ctx context.Context, client *mixin.Client, msg *model.MixinMessage) {
	msg.ConversationId = mixin.UniqueConversationID(client.ClientID, msg.UserId)
	input := &mixin.MessageRequest{
		ConversationID: msg.ConversationId,
		RecipientID:    msg.UserId,
		MessageID:      mixin.RandomTraceID(),
		Category:       mixin.MessageCategoryPlainText,
		Data:           msg.Content,
	}
	_ = client.SendMessage(ctx, input)
}
