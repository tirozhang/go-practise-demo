protoc --go_out=.  --go-grpc_out=.   trip.proto
protoc -I . --grpc-gateway_out . --grpc-gateway_opt grpc_api_configuration=trip.yaml trip.proto


