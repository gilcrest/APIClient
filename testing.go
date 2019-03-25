package apiclient

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/gilcrest/env"
	"github.com/gilcrest/env/datastore"
	"github.com/gilcrest/servertoken"
	"github.com/rs/zerolog"
)

// TestAPIClientHelper creates a client to be used in testing only
func TestAPIClientHelper(t *testing.T) *Client {
	t.Helper()

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

	token := servertoken.ServerToken(os.Getenv("TEST_SERVER_TOKEN"))
	ctx := context.Background()
	ctx = token.Add2Ctx(ctx)

	env, err := env.NewEnv(env.Dev, zerolog.DebugLevel)
	if err != nil {
		t.Fatalf("Client err: %s", err)
	}
	// get a new DB Tx
	tx, err := env.DS.BeginTx(ctx, nil, datastore.AppDB)
	if err != nil {
		t.Fatalf("Client err: %s", err)
	}

	// Call the CreateClientDB method of the Client object
	// to write to the db
	err = client.CreateClientDB(ctx, env.Logger, tx)
	if err != nil {
		fmt.Println(err)
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			t.Fatalf("Client err: %s", err)
		}
	}

	if err := tx.Commit(); err != nil {
		t.Fatalf("Client err: %s", err)
	}
	return client
}

// TestDeleteAPIClientHelper creates a client to be used in testing only
func TestDeleteAPIClientHelper(t *testing.T, c Client) {
	t.Helper()

}
