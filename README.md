# tgio

Forward input to a telegram chat

## Install

```
go get github.com/igolaizola/tgio/cmd/tgio
```

## Example

```
echo hello | tgio --token <my-bot-token> -chat <my-chat-id>
```

Or using a config file

```
echo hello | tgio --config tgio.conf
```

where `tgio.conf` content is:

```
token <my-bot-token>
chat <my-chat-id>
```

## How to get token and chat parameters

 - Talk to @BotFather to create a new bot with its token
 - Talk to @username_to_id_bot to obtain your chat ID or any other chat ID

## Use it in your code

```
import "github.com/igolaizola/tgio"

...

err := tgio.Forward(ctx, reader, token, chat)
```
