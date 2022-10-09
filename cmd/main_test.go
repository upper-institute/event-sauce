package main_test

import (
	"context"
	"flag"
	"testing"
	"time"

	"github.com/google/uuid"
	apiv1 "github.com/upper-institute/event-sauce/pkg/api/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	listenAddr string
	snapshot   = &apiv1.Snapshot{
		Id:               uuid.New().String(),
		Version:          2,
		NaturalTimestamp: timestamppb.Now(),
		Label:            "paprika_label",
	}
	durationPayload = durationpb.New(5 * time.Second)
)

func init() {
	flag.StringVar(&listenAddr, "listenAddr", "192.168.0.101:6336", "Paprika service server dial address")
}

func getSnapshotStoreClient(t *testing.T) (context.Context, context.CancelFunc, apiv1.SnapshotStoreClient) {

	conn, err := grpc.Dial(listenAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		t.Fatal(err)
	}

	snapshotClient := apiv1.NewSnapshotStoreClient(conn)

	ctx := context.Background()

	deadline, ok := t.Deadline()

	var cancel context.CancelFunc

	if ok {

		ctx, cancel = context.WithDeadline(ctx, deadline)
	}

	return ctx, cancel, snapshotClient

}

func TestRepeatedTransaction(t *testing.T) {
}

func TestSet(t *testing.T) {

	ctx, cancel, snapshotClient := getSnapshotStoreClient(t)

	defer cancel()

	payload, err := anypb.New(durationPayload)

	if err != nil {
		t.Fatal(err)
	}

	snapshot.Payload = payload

	_, err = snapshotClient.Set(ctx, snapshot)

	if err != nil {
		t.Error(err)
	}

}

func TestGet(t *testing.T) {

	ctx, cancel, snapshotClient := getSnapshotStoreClient(t)

	defer cancel()

	resSnapshot, err := snapshotClient.Get(ctx, &apiv1.Snapshot_GetRequest{
		Id: snapshot.Id,
	})

	if err != nil {
		t.Fatal(err)
	}

	if resSnapshot.Label != snapshot.Label {
		t.Errorf("Labels don't match, received '%s' expected '%s'", resSnapshot.Label, snapshot.Label)
	}

	if !resSnapshot.NaturalTimestamp.AsTime().Equal(snapshot.NaturalTimestamp.AsTime()) {
		t.Errorf("NaturalTimestamp don't match, received '%s' expected '%s'", resSnapshot.NaturalTimestamp.String(), snapshot.NaturalTimestamp.String())
	}

	if !resSnapshot.StoreTimestamp.AsTime().Equal(snapshot.StoreTimestamp.AsTime()) {
		t.Errorf("StoreTimestamp don't match, received '%s' expected '%s'", resSnapshot.StoreTimestamp.String(), snapshot.StoreTimestamp.String())
	}

	msg, err := resSnapshot.Payload.UnmarshalNew()

	if err != nil {
		t.Fatal(err)
	}

	switch msg := msg.(type) {

	case *durationpb.Duration:

		if msg.AsDuration() != durationPayload.AsDuration() {
			t.Errorf("Payload don't match, received '%v' expected '%v'", msg, durationPayload)

		}

	default:

		t.Errorf("Invalid returned payload: %v", msg)

	}

}
