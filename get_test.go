package apiclient

import (
	"context"
	"database/sql"
	"reflect"
	"testing"

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

	arg1 := args{ctx, db} //srvr.Logger, tx, os.Getenv("TEST_SERVER_TOKEN")}

	tests := []struct {
		name    string
		args    args
		want    *Client
		wantErr bool
	}{
		{"Check Client", arg, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ViaServerToken(tt.args.ctx, tt.args.db)
			if (err != nil) != tt.wantErr {
				t.Errorf("ViaServerToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ViaServerToken() = %v, want %v", got, tt.want)
			}
		})
	}
}
