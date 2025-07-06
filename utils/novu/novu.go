package novu

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"users-service/constants"
	"users-service/utils/configs"

	novu "github.com/lonesta/go-novu/lib"
)

func SendNovuTrigger(triggerName string, payload map[string]interface{}, emailTo map[string]interface{}) {

	NovuConfig, err := configs.Get(constants.NovuConfig)
	if err != nil {
		fmt.Printf("err getting congis: %v\n", err)
	}
	subscriberID := NovuConfig.GetString(constants.SubscriberID)
	apiKey := NovuConfig.GetString(constants.ApiKey)
	eventId := triggerName

	ctx := context.Background()
	emailTo["subscriberId"] = subscriberID
	overrides := novu.Overrides{
		Fcm: novu.Fcm{Type: novu.TypeData},
	}
	data := novu.ITriggerPayloadOptions{To: emailTo, Payload: payload, Overrides: overrides}
	novuClient := novu.NewAPIClient(apiKey, &novu.Config{BackendURL: NovuConfig.GetString(constants.BackendUrl)})

	resp, err := novuClient.EventApi.Trigger(ctx, eventId, data)
	if err != nil {
		log.Fatal("novu error", err.Error())
		return
	}

	fmt.Println(resp)
}

func SetNovuCreds(tokens []string) {
	NovuConfig, err := configs.Get(constants.NovuConfig)
	if err != nil {
		fmt.Printf("err getting congis: %v\n", err)
	}
	apiKey := NovuConfig.GetString(constants.ApiKey)
	subscriberID := NovuConfig.GetString(constants.SubscriberID)
	ctx := context.Background()
	novuClient := novu.NewAPIClient(apiKey, &novu.Config{BackendURL: NovuConfig.GetString(constants.BackendUrl)})

	resp, err := novuClient.SubscriberApi.SetCredentials(ctx, subscriberID, novu.PushProviderFCM, novu.ChannelCredentials{DeviceTokens: tokens})

	fmt.Printf("err: %v\n", err)
	jsonData, err := json.Marshal(resp.Data)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("resp: %v\n", string(jsonData))
}
