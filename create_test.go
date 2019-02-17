package apiclient

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/gilcrest/servertoken"
	"github.com/gilcrest/srvr"
	"github.com/gilcrest/srvr/datastore"
	"github.com/rs/zerolog"
)

func TestClient_validate(t *testing.T) {
	type fields struct {
		Number             int
		ID                 string
		Name               string
		ServerToken        string
		HomeURL            string
		Description        string
		RedirectURI        string
		PrimaryUserID      string
		Secret             string
		CreateClientNumber int
		CreateTimestamp    time.Time
		UpdateClientNumber int
		UpdateTimestamp    time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				Number:             tt.fields.Number,
				ID:                 tt.fields.ID,
				Name:               tt.fields.Name,
				ServerToken:        tt.fields.ServerToken,
				HomeURL:            tt.fields.HomeURL,
				Description:        tt.fields.Description,
				RedirectURI:        tt.fields.RedirectURI,
				PrimaryUserID:      tt.fields.PrimaryUserID,
				Secret:             tt.fields.Secret,
				CreateClientNumber: tt.fields.CreateClientNumber,
				CreateTimestamp:    tt.fields.CreateTimestamp,
				UpdateClientNumber: tt.fields.UpdateClientNumber,
				UpdateTimestamp:    tt.fields.UpdateTimestamp,
			}
			if err := c.validate(); (err != nil) != tt.wantErr {
				t.Errorf("Client.validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func NewTestAPIClient(persist bool) (*Client, error) {

	client := new(Client)

	client.Name = "Mock Client"
	client.HomeURL = "www.testurl.com"
	client.Description = "Mock Description"
	client.RedirectURI = "Mock RedirectURI"
	client.PrimaryUserID = "gilcrest"

	err := client.Finalize()
	if err != nil {
		return nil, err
	}

	if persist {
		token := servertoken.ServerToken(os.Getenv("TEST_SERVER_TOKEN"))
		ctx := context.Background()
		ctx = token.Add2Ctx(ctx)

		srvr, err := srvr.NewServer(zerolog.DebugLevel)
		if err != nil {
			return nil, err
		}
		// get a new DB Tx
		tx, err := srvr.DS.BeginTx(ctx, nil, datastore.AppDB)
		if err != nil {
			return nil, err
		}

		// Call the CreateClientDB method of the Client object
		// to write to the db
		err = client.CreateClientDB(ctx, tx)
		if err != nil {
			fmt.Println(err)
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return nil, err
			}
		}

		if err := tx.Commit(); err != nil {
			return nil, err
		}
	}
	return client, nil
}
