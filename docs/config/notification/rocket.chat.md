## Rocket.chat

```yml
rocketchat:
  serverURL: server
  userID: 123456789
  token: token
  roomIDs:
    - "roomID1"
    - "#channel"
    - "@user"
```

### Token and userID

To generate a new token:

1. Click on your profile and go to `My account > Personal Access Tokens`.
2. Enter the name of the token, set to `Ignore (Two Factor Authentication)` and click `Add`.
3. Copy the token and the user ID.

> Visit Rocket.chat's official docs [post](https://docs.rocket.chat/guides/user-guides/user-panel/managing-your-account/personal-access-token) for detailed information about personal access tokens.

### Room IDs

The room ID could be the ID of a direct message, the name of a channel or a username. Channel names should be preceded by a **\#** (i.e. #general) and usernames by **@** (i.e. @my-user).

> Direct messages IDs can be found on the web browser when visiting Rocket.chat.