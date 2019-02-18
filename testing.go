package apiclient

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/gilcrest/servertoken"
	"github.com/gilcrest/srvr"
	"github.com/gilcrest/srvr/datastore"
	"github.com/rs/zerolog"
)

// NewTestAPIClient creates a client to be used in testing only
func NewTestAPIClient(t *testing.T, persist bool) *Client {

	client := new(Client)

	client.Name = "Test Client"
	client.HomeURL = "www.testurl.com"
	client.Description = "Test Description"
	client.RedirectURI = "Test RedirectURI"
	client.PrimaryUserID = "gilcrest"

	err := client.Finalize()
	if err != nil {
		t.Fatalf("Client err: %s", err)
	}

	if persist {
		token := servertoken.ServerToken(os.Getenv("TEST_SERVER_TOKEN"))
		ctx := context.Background()
		ctx = token.Add2Ctx(ctx)

		srvr, err := srvr.NewServer(zerolog.DebugLevel)
		if err != nil {
			t.Fatalf("Client err: %s", err)
		}
		// get a new DB Tx
		tx, err := srvr.DS.BeginTx(ctx, nil, datastore.AppDB)
		if err != nil {
			t.Fatalf("Client err: %s", err)
		}

		// Call the CreateClientDB method of the Client object
		// to write to the db
		err = client.CreateClientDB(ctx, tx)
		if err != nil {
			fmt.Println(err)
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				t.Fatalf("Client err: %s", err)
			}
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("Client err: %s", err)
		}
	}
	return client
}
