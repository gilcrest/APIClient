package apiclient

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gilcrest/errors"
)

// delete is a physical delete from the database.
// This is only to be used for tests. Use logicalDelete
// for normal deletes
func (c *Client) delete(ctx context.Context, tx *sql.Tx) error {
	const op errors.Op = "apiclient/delete"

	sql := `DELETE FROM auth.client
	         WHERE client_num = $1;`
	result, err := tx.Exec(sql, c.Number)
	if err != nil {
		return errors.E(op, err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		return errors.E(op, err)
	}

	if count != 1 {
		return errors.E(op, fmt.Sprintf("Rows Deleted should be 1, but was %d", count))
	}

	return nil
}

// logicalDelete is a logical delete of the record in the database.
// This method will "end date" the record of the client, rendering
// it inactive/unusable
func (c *Client) logicalDelete() error {
	//TODO
	return nil
}
