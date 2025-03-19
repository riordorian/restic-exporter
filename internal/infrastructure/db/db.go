package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"reflect"
)

type Db struct {
	url string
	dbx *sqlx.DB
}

const driverName = "pgx"

func GetContextDb(ctx context.Context, dsn string) (context.Context, error) {
	dbx := &Db{
		url: dsn,
	}
	if err := dbx.Connect(ctx); err != nil {
		return nil, err
	}

	ctx = context.WithValue(ctx, "ctxDb", dbx)

	return ctx, nil
}

func GetDb(ctx context.Context) (*Db, error) {
	db, ok := ctx.Value("ctxDb").(*Db)

	// TODO: Type assertion?
	if ok && db.dbx != nil && reflect.TypeOf(db.dbx).String() == "*sqlx.DB" {
		return db, nil
	}

	return nil, errors.New("no database connection in context")
}

func (d *Db) Connect(ctx context.Context) error {
	fmt.Println("DB Start")
	dbx, err := sqlx.ConnectContext(ctx, driverName, d.url)

	if err != nil {
		return err
	}

	d.dbx = dbx
	return nil
}

func (d *Db) Close() error {
	return d.dbx.Close()
}

// TODO: Test for error case in MakeTransaction. It should not insert any data to the database
func (d *Db) MakeTransaction(ctx context.Context, tFunc func(ctx context.Context) error) error {
	// begin transaction
	tx, err := d.dbx.Beginx()
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	// run callback
	err = tFunc(addCtxTransact(ctx, tx))
	if err != nil {
		// if error, rollback
		if errRbck := tx.Rollback(); errRbck != nil {
			if errs := errors.Join(err, errRbck); errs != nil {
				panic("Errors during transaction")
			}
		}
		return err
	}

	// if no error, commit
	if errCmt := tx.Commit(); errCmt != nil {
		return errCmt
	}
	return nil
}

func (d *Db) Model(ctx context.Context) QuerierInterface {
	qa := &QueriesAdapter{}
	if tx := extractCtxTransact(ctx); tx != nil {
		return qa.adapt(d.dbx, tx)
	}

	return qa.adapt(d.dbx, nil)
}

// QueriesAdapter realization
type QueriesAdapter struct{}

type QuerierInterface interface {
	Queryx(query string, args ...interface{}) (*sqlx.Rows, error)
	QueryRowx(query string, args ...interface{}) *sqlx.Row
	Select(dest interface{}, query string, args ...interface{}) error
	Get(dest interface{}, query string, args ...interface{}) error
	NamedExec(query string, arg interface{}) (sql.Result, error)
	NamedQuery(query string, arg interface{}) (*sqlx.Rows, error)
}

// adapt is a method of QueriesAdapter that returns a QuerierInterface. It takes in
// a dbx pointer and a tx pointer as parameters and returns a QuerierInterface.
func (q QueriesAdapter) adapt(dbx *sqlx.DB, tx *sqlx.Tx) QuerierInterface {
	if tx != nil {
		return tx
	}

	return dbx
}

func (d *Db) Queryx(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error) {
	return d.Model(ctx).Queryx(query, args)
}
func (d *Db) QueryRowx(ctx context.Context, query string, args ...interface{}) *sqlx.Row {
	return d.Model(ctx).QueryRowx(query, args)
}
func (d *Db) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return d.Model(ctx).Select(dest, query, args)
}
func (d *Db) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return d.Model(ctx).Get(dest, query, args)
}
func (d *Db) NamedExec(ctx context.Context, query string, arg interface{}) (sql.Result, error) {
	return d.Model(ctx).NamedExec(query, arg)
}
func (d *Db) NamedQuery(ctx context.Context, query string, arg interface{}) (*sqlx.Rows, error) {
	return d.Model(ctx).NamedQuery(query, arg)
}
func (d *Db) PrepareNamed(query string) (*sqlx.NamedStmt, error) {
	return d.dbx.PrepareNamed(query)
}
