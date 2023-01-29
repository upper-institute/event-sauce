package postgres

import (
	"database/sql"
	"fmt"
	"net/url"
	"strings"

	_ "github.com/lib/pq"

	"github.com/upper-institute/flipbook/internal/drivers"
	"github.com/upper-institute/flipbook/internal/drivers/postgres/database"
	"github.com/upper-institute/flipbook/internal/helpers"
)

const (
	postgresUrl_flag = "postgres.url"
	batchSize_flag   = "postgres.batch.size"
)

type PostgresStoreAdapter struct {
	db *sql.DB
}

func (p *PostgresStoreAdapter) Bind(binder helpers.FlagBinder) {

	binder.BindString(postgresUrl_flag, "postgres://postgres@127.0.0.1:5432?dbname=mydb&sslmode=disable", "URL with parameters to connect to posgres (driver 'https://github.com/lib/pq')")
	binder.BindInt64(batchSize_flag, 50, "Default batch size to limit at Postgres select queries")

}

func (p *PostgresStoreAdapter) New(getter helpers.FlagGetter) (drivers.StoreDriver, error) {

	postgresUrl, err := url.Parse(getter.GetString(postgresUrl_flag))
	if err != nil {
		return nil, err
	}

	source := strings.Builder{}

	username := postgresUrl.User.Username()
	password, _ := postgresUrl.User.Password()

	if len(username) == 0 {
		return nil, fmt.Errorf("Missing postgres username parameter")
	}

	source.WriteString("user=")
	source.WriteString(username)

	if len(password) > 0 {
		source.WriteString(" password=")
		source.WriteString(password)
	}

	if len(postgresUrl.Port()) == 0 {
		postgresUrl.Host = fmt.Sprintf("%s:%d", postgresUrl.Host, 5432)
	}

	address := strings.Split(postgresUrl.Hostname(), ":")

	source.WriteString(" host=")
	source.WriteString(address[0])
	source.WriteString(" port=")
	source.WriteString(address[1])

	dbname := postgresUrl.Query().Get("dbname")
	if len(username) == 0 {
		return nil, fmt.Errorf("Missing postgres dbname parameter")
	}

	source.WriteString(" dbname=")
	source.WriteString(dbname)

	sslmode := postgresUrl.Query().Get("sslmode")
	if len(username) > 0 {
		source.WriteString(" sslmode=")
		source.WriteString(sslmode)
	}

	db, err := sql.Open("postgres", source.String())
	if err != nil {
		return nil, err
	}

	p.db = db

	return &postgresStore{
		db:               p.db,
		queries:          database.New(p.db),
		defaultBatchSize: int32(getter.GetInt64(batchSize_flag)),
	}, nil
}

func (p *PostgresStoreAdapter) Destroy(store drivers.StoreDriver) error {

	if p.db != nil {

		if err := p.db.Close(); err != nil {
			return err
		}

	}
	return nil
}
