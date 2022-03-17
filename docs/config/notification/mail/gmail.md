## Gmail

```yml
mail:
    smtpHostAddress: smtp.gmail.com # required
    smtpHostPort: 587 # required
    senderAddress: sender@gmail.com # required
    senderPassword: senderPassowrd # required
    receiverAddresses:
      - receiver@mail.com
```

### Sender password

Using the password of your Google account to send mails now requires an extra step, which depends on whether your account has the 2-step verification enabled.

If it's **enabled**, [create](https://myaccount.google.com/apppasswords) and use an app-specific password. Select the app *"Gmail"* and the device *"Other"*, the password given can then be placed as the value of `senderPassword`.

If it's **not enabled**, you will have to turn on [`less secure apps`](https://support.google.com/accounts/answer/6010255?hl=en) and use the Google's account password.
