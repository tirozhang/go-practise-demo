package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"

	trippb "github.com/tirozhang/go-practise-demo/grpc-gateway/trip/gen/v1"
)

type Server struct {
	trippb.UnimplementedTripServiceServer
}

func (*Server) GetTrip(ctx context.Context, req *trippb.GetTripReq) (*trippb.GetTripResp, error) {
	return &trippb.GetTripResp{
		Trip: &trippb.Trip{
			Start:       "abc",
			End:         "def",
			DurationSec: 3600,
			FeeCent:     10000,
			StartPos: &trippb.Location{
				Latitude:  100,
				Longitude: 200,
			},
			Status: trippb.TripStatus_TS_NOT_SPEC,
		},
	}, nil
}

func main() {
	// Write your code here
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	s := grpc.NewServer()
	trippb.RegisterTripServiceServer(s, &Server{})
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
