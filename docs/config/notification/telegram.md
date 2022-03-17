## Telegram

```yml
telegram:
    botAPIToken: botAPIToken # required
    chatIDs: 
      - chatID1
      - chatID2
```

### Bot API token

Creating a bot API token in telegram is pretty simple:

1. Search for the user [@Botfather](https://t.me/Botfather).
2. Type the command *"/newbot"* and send it.
3. Choose a name and username for the bot.
4. That's it, the next message should contain the API used to send messages.

### Chat IDs

Once you invite the bot to your group, it will output a JSON object which contains the chat ID at **message.chat.id**.

If it didn't work, send a message that contains its username and make a request to get the updates from the bot

```
curl https://api.telegram.org/bot<botAPIToken>/getUpdates
```