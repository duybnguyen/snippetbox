package mysql

import (
	"database/sql"

	"github.com/duybnguyen/snippetbox/pkg/models"
)

// DB.Query() is used for SELECT queries which return multiple rows. DB.
// QueryRow() is used for SELECT queries which return a single row.
// DB.Exec() is used for statements which don’t return rows (like INSERT and DELETE).

// wraps a sql.DB connection pool
// By creating a custom SnippetModel type and implementing methods on it we’ve been able to make our model a single, neatly encapsulated object, which we can easily initialize and then pass to our handlers as a dependency.
type SnippetModel struct {
	DB *sql.DB
}

// insert a new snippet into db
func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	query := `INSERT INTO snippets (title, content, created, expires)
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := m.DB.Exec(query, title, content, expires)
	if err != nil {
		return 0, err
	}

	// get the ID of the newly inserted record
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	// The ID returned has the type int64, so we convert it to an int type
	// before returning.
	return int(id), nil
}

// select a snippet from db using its id
func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	query := `SELECT id, title, content, created, expires FROM snippets
    WHERE expires > UTC_TIMESTAMP() AND id = ?`

	row := m.DB.QueryRow(query, id)

	s := &models.Snippet{}

	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}
	return s, nil
}

// select latest 10 snippets from db
func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	query := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() ORDER BY created DESC LIMIT 10`

	//This returns a sql.Rows resultset containing the result of our query
	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	// Closing a resul tset iscritical here. As long as a resultset is open it will keep the underlying database connection open... so if something goes wrong in this method and the resultset isn’t closed, it can rapidly lead to all the connections in your pool being used up.

	defer rows.Close()

	snippets := []*models.Snippet{}

	for rows.Next() {
		s := &models.Snippet{}
		// copy the values from each field in the row to the
		// new Snippet object that we created
		// the arguments to row.Scan
		// must be pointers to the place you want to copy the data into
		// number of arguments must be exactly the same as the number of
		// columns returned by your statement.
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, s)
	}

	// When the rows.Next() loop has finished we call rows.Err() to retrieve an
	// error that was encountered during the iteration.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil

}
