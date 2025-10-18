package document

import "context"

//go:generate mockgen -destination=../mocks/mock$GOPACKAGE/$GOFILE -package=mock$GOPACKAGE . Repository

type Document struct {
	ID int64
}

type Repository interface {
	FindDocuments(ctx context.Context, signerID int64, status Status) ([]Document, error)
}
