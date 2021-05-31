# tgio

Forward input to a telegram chat

## Install

```
go install github.com/igolaizola/tgio/cmd/tgio@latest
```

## Example

```
ping google.com | tgio --token <my-bot-token> -chat <my-chat-id>
```

## How to get token and chat parameters

 - Talk to @BotFather to create a new bot with its token
 - Talk to @username_to_id_bot to obtain your chat ID or any other chat ID

## Use it in your code

```
import "github.com/igolaizola/tgio"

...

err := tgio.Forward(reader, token, chat)
```
