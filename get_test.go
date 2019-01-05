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
		db  *sql.DB
	}

	srvr, err := srvr.NewServer(zerolog.DebugLevel)
	if err != nil {
		t.Errorf("Error from Newserver = %v", err)
	}
	db, err := srvr.DS.DB(datastore.AppDB)
	if err != nil {
		t.Errorf("Error getting DB = %v", err)
	}

	token1 := servertoken.ServerToken("1234567")
	ctx := context.Background()
	ctx = token1.Add2Ctx(ctx)

	arg1 := args{ctx, db}

	t2 := servertoken.ServerToken(os.Getenv("TEST_SERVER_TOKEN"))
	ctx2 := context.Background()
	ctx2 = t2.Add2Ctx(ctx2)

	client, err := setupClient(ctx2)
	if err != nil {
		t.Errorf("Error from setupClient = %v", err)
	}

	t3 := servertoken.ServerToken(client.ServerToken)
	ctx3 := context.Background()
	ctx3 = t3.Add2Ctx(ctx3)

	arg2 := args{ctx3, db}

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
			got, err := ViaServerToken(tt.args.ctx, tt.args.db)
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
}

func setupClient(ctx context.Context) (*Client, error) {
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

	srvr, err := srvr.NewServer(zerolog.DebugLevel)
	if err != nil {
		return nil, errors.E(op, fmt.Errorf("Error = %v", err))
	}

	tx, err := srvr.DS.BeginTx(ctx, nil, datastore.AppDB)
	if err != nil {
		return nil, errors.E(op, fmt.Errorf("Error = %v", err))
	}

	tx, err = client.CreateClientDB(ctx, tx)
	if err != nil {
		return nil, errors.E(op, fmt.Errorf("Error = %v", err))
	}

	if !client.DMLTime.IsZero() {
		err := tx.Commit()
		if err != nil {
			return nil, errors.E(op, fmt.Errorf("Error = %v", err))
		}
	} else {
		err = tx.Rollback()
		if err != nil {
			return nil, errors.E(op, fmt.Errorf("Error = %v", err))
		}
	}

	return client, nil

}
