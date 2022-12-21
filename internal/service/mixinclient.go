package service

import (
	"log"
	"sync"

	"github.com/lixvyang/betxin/internal/utils"

	"github.com/fox-one/mixin-sdk-go"
)

var (
	InitMixin   sync.Once
	MixinClient *mixin.Client
	err         error
)

func InitMixinClient() {
	store := &mixin.Keystore{
		ClientID:   utils.ClientId,
		SessionID:  utils.SessionId,
		PrivateKey: utils.PrivateKey,
		PinToken:   utils.PinToken,
	}

	MixinClient, err = mixin.NewFromKeystore(store)
	if err != nil {
		log.Fatal(err)
	}
}
