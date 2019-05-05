package apiclient

import (
	"context"
	"database/sql"
	"time"

	"github.com/gilcrest/errors"
	"github.com/gilcrest/rand"
	"github.com/rs/zerolog"
)

// Client is used for the client service and response
type Client struct {
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

// Validate method validates Client input data
func (c *Client) validate() error {
	const op errors.Op = "apiclient/validate"

	if len(c.ID) > 0 {
		return errors.E(op, errors.InputUnwanted("ID"))
	}
	if len(c.Name) == 0 {
		return errors.E(op, errors.MissingField("Name"))
	}
	if len(c.Description) == 0 {
		return errors.E(op, errors.MissingField("Description"))
	}
	if len(c.PrimaryUserID) == 0 {
		return errors.E(op, errors.MissingField("PrimaryUserID"))
	}
	if len(c.Secret) != 0 {
		return errors.E(op, errors.InputUnwanted("Secret"))
	}
	if len(c.ServerToken) != 0 {
		return errors.E(op, errors.InputUnwanted("ServerToken"))
	}
	if c.CreateClientNumber != 0 {
		return errors.E(op, errors.InputUnwanted("CreateClientNumber"))
	}
	if c.CreateTimestamp.IsZero() != true {
		return errors.E(op, errors.InputUnwanted("CreateTimestamp"))
	}
	if c.UpdateClientNumber != 0 {
		return errors.E(op, errors.InputUnwanted("UpdateClientNumber"))
	}
	if c.UpdateTimestamp.IsZero() != true {
		return errors.E(op, errors.InputUnwanted("UpdateTimestamp"))
	}

	return nil
}

// Finalize validates user input and generates token info
func (c *Client) Finalize() error {
	const op errors.Op = "apiclient/Finalize"

	// Validate that all user input is acceptable
	err := c.validate()
	if err != nil {
		return errors.E(op, err)
	}

	err = c.generateID()
	if err != nil {
		return errors.E(op, err)
	}

	err = c.issueSecretToken()
	if err != nil {
		return errors.E(op, err)
	}

	err = c.issueServerToken()
	if err != nil {
		return errors.E(op, err)
	}

	return nil
}

func (c *Client) generateID() error {
	const op errors.Op = "apiclient/issueSecretToken"

	// Generate a Client ID
	id, err := rand.CryptoString(24)
	if err != nil {
		return errors.E(op, err)
	}

	c.ID = id

	return nil
}

func (c *Client) issueSecretToken() error {
	const op errors.Op = "apiclient/issueSecretToken"

	// Generate a Client Secret
	id, err := rand.CryptoString(30)
	if err != nil {
		return errors.E(op, err)
	}

	c.Secret = id

	return nil
}

func (c *Client) issueServerToken() error {
	const op errors.Op = "apiclient/issueServerToken"

	// Generate a Client Secret
	id, err := rand.CryptoString(30)
	if err != nil {
		return errors.E(op, err)
	}

	c.ServerToken = id

	return nil
}

// grant_types     VARCHAR(80),
// scope           VARCHAR(4000),

// CreateClientDB creates a client/app in the database
func (c *Client) CreateClientDB(ctx context.Context, log zerolog.Logger, tx *sql.Tx) error {
	const op errors.Op = "apiclient/CreateClientDB"

	// Get the API client that is creating the new API client :)
	createClient, err := FromCtx(ctx)
	if err != nil {
		return errors.E(op, err)
	}

	// Prepare the sql statement using bind variables
	stmt, err := tx.PrepareContext(ctx, `
	select o_client_num, 
	       o_create_client_num, o_create_timestamp,
           o_update_client_num, o_update_timestamp
	  from auth.create_client (
			p_client_id => $1,
			p_client_name => $2,
			p_server_token => $3,
			p_homepage_url => $4,
			p_app_description => $5,
			p_redirect_uri => $6,
			p_client_secret => $7,
			p_primary_username => $8,
			p_create_client_num => $9)`)

	if err != nil {
		return errors.E(op, err)
	}
	defer stmt.Close()

	// Execute stored function that returns the create_date timestamp,
	// hence the use of QueryContext instead of Exec
	rows, err := stmt.QueryContext(ctx,
		c.ID,                //$1
		c.Name,              //$2
		c.ServerToken,       //$3
		c.HomeURL,           //$4
		c.Description,       //$5
		c.RedirectURI,       //$6
		c.Secret,            //$7
		c.PrimaryUserID,     //$8
		createClient.Number) //$9

	if err != nil {
		return errors.E(op, err)
	}
	defer rows.Close()

	// Iterate through the returned record(s)
	for rows.Next() {
		if err := rows.Scan(&c.Number, &c.CreateClientNumber, &c.CreateTimestamp,
			&c.UpdateClientNumber, &c.UpdateTimestamp); err != nil {
			return errors.E(op, err)
		}
	}

	if err := rows.Err(); err != nil {
		return errors.E(op, err)
	}

	return nil

}
