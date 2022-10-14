package main

import (
	"context"
	"log"
	"net/http"

	authpb "github.com/tirozhang/go-practise-demo/login/api/gen/v1"

	"google.golang.org/grpc"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	ctx := context.Background()
	ctx, cancelFunc := context.WithCancel(ctx)
	defer cancelFunc()
	mux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames:  true,
			UseEnumNumbers: true,
		}}),
	)
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := authpb.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, "localhost:50051", opts)
	if err != nil {
		log.Fatalf("failed to register gateway: %v", err)
	}
	err = http.ListenAndServe(":8081", mux)
	if err != nil {
		log.Fatalf("failed to listen and serve: %v", err)
	}
}
