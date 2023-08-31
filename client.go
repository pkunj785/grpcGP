package main

import (
	"context"
	"log"
	"math/rand"
	"time"

	"com.yyxx/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pbtime "google.golang.org/protobuf/types/known/timestamppb"
)

func main() {
	addr := "localhost:9091"
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	client := pb.NewAnomsClient(conn)
	req := pb.AnomRequest {
		Metrics: dummydata(),
	}

	resp, errn := client.Expose(context.Background(), &req)
	if errn != nil {
		log.Fatal(errn)
	}
	log.Printf("outliers at: %+v", resp.Indices)
}

func dummydata() []*pb.Metric {
	const size = 1000
	out := make([]*pb.Metric, size)
	t := time.Date(2023, 7, 13, 20, 15, 00, 00, time.UTC)
	for i := 0; i < size; i++ {
		m := pb.Metric {
			Time: Timestamp(t),
			Name: "CPU",
			Value: rand.Float64() * 40,
		}
		out[i] = &m
		t.Add(time.Second)
	}

	out[7].Value = 96.2
	out[113].Value = 92.1
	out[875].Value = 94.3
	return out
}


func Timestamp(t time.Time) *pbtime.Timestamp {
	return &pbtime.Timestamp{
		Seconds: t.Unix(),
		Nanos: int32(t.Nanosecond()),
	}
}