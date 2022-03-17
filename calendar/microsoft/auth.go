package microsoft

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	azure "github.com/microsoft/kiota/authentication/go/azure"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
)

const redirectURL = "https://login.microsoftonline.com/common/oauth2/authorize"

func getClient(account *Account) (*msgraphsdk.GraphServiceClient, error) {
	cred, err := azidentity.NewAuthorizationCodeCredential(
		account.TenantID,
		account.ClientID,
		"", // TODO: get authorization code like in Google service
		redirectURL,
		nil,
	)
	if err != nil {
		return nil, err
	}

	// scopes := []string{"User.Read", "Calendars.Read", "Calendars.Read.Shared"}
	auth, err := azure.NewAzureIdentityAuthenticationProviderWithScopes(cred, nil)
	if err != nil {
		return nil, err
	}

	adapter, err := msgraphsdk.NewGraphRequestAdapter(auth)
	if err != nil {
		return nil, err
	}

	client := msgraphsdk.NewGraphServiceClient(adapter)

	return client, nil
}
