# OwnProvider
An APNs Provider Server which based on JWT auth for iOS.  Build your own Apple Provider Server even you are not a developer.

# Env
go version go1.15 darwin/amd64

# Build & Install for Linux server
```shell
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -trimpath -o ownprovider  main.go
```

```shell
/where/your/privider/install/ownprovider
```

# Push a Alert Message
```shell
curl -X POST "http://127.0.0.1:27953/ownprovider/inner/push" -d 'payload={"aps":{"alert":{"title":"Hello","body":"Baby"},"badge":1,"sound"     :"default","type":"6"}}' -d 'topic=com.example.app&env=sandbox&token=YourDeviceToken'

```

#Push a VoIP Message - Note: VoIP device token is different with alert device token
```
curl -X POST "http://127.0.0.1:27953/ownprovider/inner/push" -d 'payload={"aps":{"alert":{"title":"Hello","body":"Baby"},"badge":1,"sound"     :"default","type":"6"}}' -d 'topic=com.example.app&voip=voip&env=sandbox&token=YourVoIPDeviceToken'
```
