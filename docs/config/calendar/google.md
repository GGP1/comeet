## Google calendar

```yml
enabled: true
accounts:
	-
	clientID: clientID # required
	clientSecret: clientSecret # required
	tokenPath: /path/to/tokenFile
	calendarID: primary
```

On the first execution you will be prompted for a verification that is obtained by following the URL provided and authenticating with the account for the calendar you use. If you are using multiple accounts, you will be prompted many times.

This will create a token file to avoid repeating this process. Note that anyone with access to it can retrieve your calendar information (your Google account isn't at risk though).

### Client credentials

In order to obtain the **clientID** and the **clientSecret**, we will have to create a project inside the Google Cloud Plaform, for it, follow these steps:

1. Go to [Google Cloud Console - New project](https://console.cloud.google.com/projectcreate) (sign-in with a Google account if you didn't already).
2. Go to the `APIs & Services` dashboard and click `Enable APIs and Services`.
3. Look for and enable the Google Calendar API.
4. Go to the `Credentials` dashboard, click the `Create credentials` button and select *"OAuth client ID"*.
5. Select "Desktop app" as the application type and create.
6. Use the credentials given in Comeet's configuration file.
