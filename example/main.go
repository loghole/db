package main

import (
	"context"
	"database/sql"
	"log"
	"runtime"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	outDirver "github.com/loghole/db/driver"
)

var (
	testHook   = new(Hook)
	dsn        = "postgresql://root@localhost:29999/defaultdb?sslmode=disable"
	driverName = "test-pq"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Ltime)

	sql.Register(driverName, outDirver.Wrap(&pq.Driver{}, testHook))

	sqlxSQL()

	//defaultSQL()

	log.Println("ok")
}

type TestToken struct {
	Token string `db:"token"`
}

func sqlxSQL() {
	dbSTD, err := sql.Open(driverName, dsn)
	if err != nil {
		log.Fatal(err)
	}

	db := sqlx.NewDb(dbSTD, "postgres")

	rows, err := db.QueryxContext(context.Background(), `SELECT token FROM tokens limit 2`)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var id string
		err := rows.Scan(&id)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("id", id)
	}

	_, err = db.ExecContext(context.Background(), `INSERT INTO tokens(token) VALUES($1)`, strconv.FormatInt(time.Now().Unix(), 32))
	if err != nil {
		log.Panicln(err)
	}

	tt := TestToken{Token: strconv.FormatInt(time.Now().Unix(), 30)}

	if _, err := db.NamedExecContext(context.Background(), `INSERT INTO tokens(token) VALUES(:token)`, tt); err != nil {
		log.Panicln(err)
	}
}

func defaultSQL() {
	db, err := sql.Open(driverName, dsn)
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.QueryContext(context.Background(), `SELECT id FROM tokens`)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			log.Fatal(err)
		}

		log.Println("id", id)
	}
}

type Hook struct {
}

func (h *Hook) Before(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
	log.Println("before hook", query)

	log.Println(getFrame(7))

	return ctx, nil
}

func (h *Hook) After(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
	log.Println("after hook", query)

	return ctx, nil
}

func getFrame(skipFrames int) runtime.Frame {
	// We need the frame at index skipFrames+2, since we never want runtime.Callers and getFrame
	targetFrameIndex := skipFrames + 2

	// Set size to targetFrameIndex+2 to ensure we have room for one more caller than we need
	programCounters := make([]uintptr, targetFrameIndex+2)
	n := runtime.Callers(0, programCounters)

	frame := runtime.Frame{Function: "unknown"}
	if n > 0 {
		frames := runtime.CallersFrames(programCounters[:n])
		for more, frameIndex := true, 0; more && frameIndex <= targetFrameIndex; frameIndex++ {
			var frameCandidate runtime.Frame
			frameCandidate, more = frames.Next()
			if frameIndex == targetFrameIndex {
				frame = frameCandidate
			}
		}
	}

	return frame
}
