# OwnProvider
An APNs Provider Server which based on JWT auth for iOS.  Build your own Apple Provider Server even you are not a developer.

# Env
go version go1.11 darwin/amd64

# Build & Install
```shell
go build -o ownprovider
```

```shell
/where/your/privider/install/ownprovider
```

# Push a Message
```shell
curl -X POST "http://127.0.0.1:9527/api/notify" -d 'topic=YourBundleId&token=YourDeviceToken&playload={"aps":{"alert":"Hello"}}'
```

# Building...