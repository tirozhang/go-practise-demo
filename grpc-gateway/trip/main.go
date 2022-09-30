package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"google.golang.org/protobuf/encoding/protojson"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"google.golang.org/grpc"

	"google.golang.org/protobuf/proto"

	tripPB "github.com/tirozhang/go-practise-demo/grpc-gateway/trip/gen/v1"
)

func main() {
	// Write your code here
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	t1 := tripPB.Trip{
		Start:       "abc",
		End:         "def",
		DurationSec: 3600,
		FeeCent:     10000,
		StartPos: &tripPB.Location{
			Latitude:  100,
			Longitude: 200,
		},
		Status: tripPB.TripStatus_TS_NOT_SPEC,
	}
	fmt.Println(&t1)
	fmt.Println(t1.String())
	fmt.Println(t1)

	b, _ := proto.Marshal(&t1)
	fmt.Printf("%X\n", b)

	var trip2 tripPB.Trip
	err := proto.Unmarshal(b, &trip2)
	if err != nil {
		panic(err)
	}
	fmt.Println(&trip2)

	b, err = json.Marshal(&trip2)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", b)
	go startGrpcClient()
	startGrpcGateway()
}

func startGrpcClient() {
	dial, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		return
	}
	defer func(dial *grpc.ClientConn) {
		err := dial.Close()
		if err != nil {
			log.Fatalf("close dial failed: %v", err)
		}
	}(dial)

	client := tripPB.NewTripServiceClient(dial)
	trip, err := client.GetTrip(context.Background(), &tripPB.GetTripReq{})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	fmt.Println(trip)
}

func startGrpcGateway() {
	// Write your code here
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	mux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames:  true,
			UseEnumNumbers: true,
		}}),
	)
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := tripPB.RegisterTripServiceHandlerFromEndpoint(ctx, mux, "localhost:50051", opts)
	if err != nil {
		log.Fatalf("failed to register gateway: %v", err)
	}
	err = http.ListenAndServe(":8081", mux)
	if err != nil {
		log.Fatalf("failed to listen and serve: %v", err)
	}
}
