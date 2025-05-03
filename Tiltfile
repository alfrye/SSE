local_resource(
  'go-compile',
  'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/authorize ./server-events/cmd/main.go',
  deps=['./main.go'],
  #resource_deps = ['deploy'],
  labels=["SSE"]
  
  )
