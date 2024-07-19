# tgio

Forward input to a telegram chat

## 📦 Installation

You can use the golang binary to install tgio:

```
go install github.com/igolaizola/tgio/cmd/tgio@latest
```

Or you can download the binary from the [releases](https://github.com/igolaizola/tgio/releases)

## 🕹️ Usage 

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

### Use it in your code

```
import "github.com/igolaizola/tgio"

...

err := tgio.Forward(ctx, reader, token, chat)
```

## 🛠️ How to get token and chat parameters

 - Talk to @BotFather to create a new bot with its token
 - Talk to @username_to_id_bot to obtain your chat ID or any other chat ID

