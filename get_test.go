package apiclient

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/gilcrest/errors"
	"github.com/gilcrest/servertoken"
	"github.com/gilcrest/srvr"
	"github.com/gilcrest/srvr/datastore"
	"github.com/rs/zerolog"
)

func TestViaServerToken(t *testing.T) {
	type args struct {
		ctx context.Context
		tx  *sql.Tx
	}

	srvr, err := srvr.NewServer(zerolog.DebugLevel)
	if err != nil {
		t.Errorf("Error from Newserver = %v", err)
	}

	token1 := servertoken.ServerToken("1234567")
	ctx := context.Background()
	ctx = token1.Add2Ctx(ctx)

	tx, err := srvr.DS.BeginTx(ctx, nil, datastore.AppDB)
	if err != nil {
		t.Errorf("Error with BeginTx = %v", err)
	}

	arg1 := args{ctx, tx}

	// Add test server token to context
	t2 := servertoken.ServerToken(os.Getenv("TEST_SERVER_TOKEN"))
	ctx2 := context.Background()
	ctx2 = t2.Add2Ctx(ctx2)

	// create a new client using ctx2
	client, err := setupTestClient(ctx2, srvr.Logger, tx)
	if err != nil {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				t.Logf("Could not roll back: %v\n", rollbackErr)
			}
			t.Errorf("Error from setupClient = %v", err)
		}
	}

	t3 := servertoken.ServerToken(client.ServerToken)
	ctx3 := context.Background()
	ctx3 = t3.Add2Ctx(ctx3)

	arg2 := args{ctx3, tx}

	tests := []struct {
		name    string
		args    args
		want    *Client
		wantErr bool
	}{
		{"Check Not Exists", arg1, nil, true},
		{"Check Exists", arg2, client, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ViaServerToken(tt.args.ctx, tt.args.tx)
			if (err != nil) != tt.wantErr {
				t.Errorf("ViaServerToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// Check non-dynamic fields are set properly
			if got != nil {
				if got.Name != tt.want.Name {
					t.Errorf("Name = %v, want %v", got.Name, tt.want.Name)
				}
			}
		})
	}

	err = Delete(ctx3, srvr.Logger, tx, client)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			t.Logf("Could not roll back: %v\n", rollbackErr)
		}
		t.Errorf("Error from client.delete = %v", err)
	}

	if err := tx.Commit(); err != nil {
		srvr.Logger.Fatal().Msgf("Could not commit: %v\n", err)
	}

}

func setupTestClient(ctx context.Context, log zerolog.Logger, tx *sql.Tx) (*Client, error) {
	const op errors.Op = "apiclient/ViaServerToken"

	client := new(Client)
	client.Name = "Test Client"
	client.HomeURL = "http://www.repomanfilm.com/"
	client.Description = "This is a fake client for testing only"
	client.RedirectURI = "http://www.repomanfilm.com/redirect"
	client.PrimaryUserID = "gilcrest"

	err := client.Finalize()
	if err != nil {
		return nil, errors.E(op, fmt.Errorf("Error = %v", err))
	}

	err = client.CreateClientDB(ctx, log, tx)
	if err != nil {
		return nil, errors.E(op, fmt.Errorf("Error = %v", err))
	}

	return client, nil

}
