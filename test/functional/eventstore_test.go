package functional_test

import (
	"context"
	"crypto/tls"
	"flag"
	"testing"
	"time"

	"github.com/google/uuid"
	flipbookv1 "github.com/upper-institute/flipbook/proto/api/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	serverAddr = flag.String("server-address", "127.0.0.1:6333", "Flipbook gRPC server address")
)

func getEventStoreClient(t *testing.T) flipbookv1.EventStoreClient {

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}

	creds := credentials.NewTLS(tlsConfig)

	conn, err := grpc.Dial(*serverAddr, grpc.WithTransportCredentials(creds))

	if err != nil {
		t.Fatal(err)
	}

	return flipbookv1.NewEventStoreClient(conn)

}

func generateEvents(count int, startSortKey int, sortKeyType flipbookv1.SortingKeyType) []*flipbookv1.Event {

	evs := []*flipbookv1.Event{}

	partKey := uuid.New()

	ts := time.Now()

	for i := 0; i < count; i++ {

		tsPb := timestamppb.New(ts.Add(time.Duration(i) * time.Hour))

		evPayload, _ := anypb.New(tsPb)

		evs = append(evs, &flipbookv1.Event{
			PartitionKey:   partKey.String(),
			SortingKey:     int64(startSortKey + i),
			SortingKeyType: sortKeyType,
			EventPayload:   evPayload,
		})
	}

	return evs

}

func TestEventStoreIncreasingSequence(t *testing.T) {

	evc := getEventStoreClient(t)

	ctx := context.Background()

	evs_1 := generateEvents(5, 0, flipbookv1.SortingKeyType_SORTING_KEY_INCREASING_SEQUENCE)
	evs_2 := generateEvents(10, 0, flipbookv1.SortingKeyType_SORTING_KEY_INCREASING_SEQUENCE)

	evs := append([]*flipbookv1.Event{}, evs_1...)
	evs = append(evs, evs_2...)

	_, err := evc.Append(ctx, &flipbookv1.Event_AppendRequest{
		Events: evs,
	})

	if err != nil {
		t.Fatal(err)
	}

}
