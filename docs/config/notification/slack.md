## Slack

```yml
slack:
	botAccessToken: botAccessToken
	channelIDs:
		- channel1
		- channel2
```

### Bot Access Token

In order to send messages to a Slack channel, we need to create an application and use its bot access token following these steps:

1. Go to your [Slack apps](https://api.slack.com/apps).
2. Click `Create New App`.
3. Select the workspace in which the app will be.
4. Under `Features > OAuth & Permissions > Scopes > Bot Token Scopes` use the `Add an OAuth Scope` button to give the app these permissions:
	- chat:write
	- chat:write.customize
	- chat:write.public
5. On the same page, after the changes were applied, go to `OAuth Tokens for Your Workspace` and copy the `Bot User OAuth Token`. This is the token Comeet will use to send message to any channel desired.

### Channel IDs

If you use Slack in a browser, the channel ID appears in the URL. If you are using the desktop application, simply right click the channel, select *"Copy link"* and extract the ID from it. 