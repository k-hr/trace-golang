package db

import (
	"context"
	"source.golabs.io/engineering-platforms/lens/trace-app-golang/model"
)

func (p *Proxy) InsertBook(ctx context.Context, book model.Book) error {
	_, err := p.db.ModelContext(ctx, &book).Insert()
	return err
}

func (p *Proxy) ReadBook(ctx context.Context, title string, delay, runtimeErr bool) (model.Book, error) {
	var book model.Book
	var err error
	if runtimeErr {
		err = p.db.ModelContext(ctx, book).Where("title = ?", title).Select()
	} else {
		err = p.db.ModelContext(ctx, &book).Where("title = ?", title).Select()
	}
	if err != nil {
		return book, err
	}

	if delay {
		var delayModel struct {
			tableName struct{} `pg:",discard_unknown_columns"`
		}
		_, err = p.db.QueryContext(ctx, &delayModel, "SELECT pg_sleep(10);")
	}
	return book, err
}
