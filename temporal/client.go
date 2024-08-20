package temporal

import "go.temporal.io/sdk/client"

var (
	temporalClient client.Client
)

func InitClient() {
	var err error
	temporalClient, err = client.NewLazyClient(client.Options{})
	if err != nil {
		panic(err)
	}
}

func GetClient() client.Client {
	return temporalClient
}
