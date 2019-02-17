package apiclient

import (
	"testing"
	"time"
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
