package db

import (
	"context"
	"source.golabs.io/engineering-platforms/lens/trace-app-golang/model"
)

func (p *Proxy) InsertBook(ctx context.Context, book model.Book) error {
	_, err := p.db.ModelContext(ctx, &book).Insert()
	return err
}

func (p *Proxy) ReadBook(ctx context.Context, title string) (model.Book, error) {
	var book model.Book
	err := p.db.ModelContext(ctx, &book).Where("title = ?", title).Select()
	return book, err
}
