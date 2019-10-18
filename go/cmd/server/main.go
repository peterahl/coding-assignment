package main

import (
	"fmt"
	"golang.org/x/net/context"
	"log"
	"net/http"

	"github.com/golang/protobuf/proto"
	lift "github.com/liftbridge-io/go-liftbridge"
	liftmodels "github.com/liftbridge-io/liftbridge-grpc/go"
	"github.com/peterahl/coding-assignment/go/pkg/memstore"
	"github.com/peterahl/coding-assignment/go/pkg/models"
)

func main() {

	db := &memstore.Store{
		Messages: make(map[uint64]models.Message),
	}

	var (
		subject    = "foo"
		streamName = "foo-stream"
	)

	addrs := []string{"localhost:9292"}
	// addrs := []string{"liftbridge:9292"}
	client, err := lift.Connect(addrs)
	if err != nil {
		panic(err)
	}

	defer client.Close()

	if err := client.CreateStream(context.Background(), subject, streamName); err != nil {
		if err != lift.ErrStreamExists {
			panic(err)
		}
	}

	ctx := context.Background()

	r := newRouter(db, client)

	liftHandler := newLiftHandler(db)

	if err := client.Subscribe(ctx, streamName, liftHandler, lift.StartAtEarliestReceived()); err != nil {
		panic(err)
	}

	fmt.Println("Starting server")

	log.Fatal(http.ListenAndServe(":3000", r))

	<-ctx.Done()
}

func newLiftHandler(db dataStore) func(msg *liftmodels.Message, err error) {
	return func(msg *liftmodels.Message, err error) {
		if err != nil {
			panic(err)
		}
		var pbMessage models.Message
		proto.Unmarshal(msg.Value, &pbMessage)
		db.AddCommand(pbMessage)
		switch pbMessage.GetCmd() {
		case "update":
			db.UpdateMessage(pbMessage)
		case "create":
			db.NewMessage(pbMessage)
		case "delete":
			db.DeleteMessage(pbMessage)
		default:
			fmt.Println(msg.Offset, pbMessage.GetText())
		}
	}
}
