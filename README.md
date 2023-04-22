# Midjourney Discord Bot

## Installing
```go
go get -u github.com/hedeqiang/midjourney
```

## Configuration
```go
client, _ := NewClient(&ClientOptions{
    BotToken:           "botToken",
    AuthorizationToken: "authToken",
    GuildId:            "1098487760261755000",
    ChannelId:          "xxx",
    ApplicationId:      "936929561302675456",
})
```

## Usage
```go
if err := client.Connect(); err != nil {
    log.Fatal(err)
}

defer client.Disconnect()

client.OnMessage(func(m *discordgo.MessageCreate) {
    if m.Author.ID == client.dg.State.User.ID {
        return
    }
    if data, err := json.Marshal(m); err == nil {
        log.Println("received message:", string(data))
    }

    if len(m.Attachments) > 0 {
        log.Println(m.Attachments[0].URL)
        url := m.Attachments[0].URL
        customId, _ := ExtractDataFromURL(url)
        log.Println(customId)
    }
})

log.Println("Bot is now running. Press CTRL-C to exit.")

select {}
```

## GenerateImage
```go
client.GenerateImage("cat")
```

## Upscale
```go
client.Upscale(2, "1099257364063985717", "3c6d3742-28bb-4354-a735-c815d82c4184")
```

## Upscale
```go
client.Reset("1099257364063985717", "3c6d3742-28bb-4354-a735-c815d82c4184")
```