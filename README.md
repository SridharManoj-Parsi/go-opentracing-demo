# go-opentracing-demo


start all the 5 services :
    go run ./cmd/istio/main.go
    go run ./cmd/fabric/main.go
    go run ./cmd/epg/main.go
    go run ./cmd/monetisation/main.go
    go run ./cmd/ads/main.go


start the traces collector:
    docker run --rm --name jaeger \
  -e COLLECTOR_ZIPKIN_HOST_PORT=:9411 \
  -p 6831:6831/udp \
  -p 6832:6832/udp \
  -p 5778:5778 \
  -p 16686:16686 \
  -p 4317:4317 \
  -p 4318:4318 \
  -p 14250:14250 \
  -p 14268:14268 \
  -p 14269:14269 \
  -p 9411:9411 \
  jaegertracing/all-in-one:1.49


call the istio gateway to get the flow going : POST  http://127.0.0.1:3000/istio