package document

import "context"

//go:generate mockgen -destination=../mocks/mock$GOPACKAGE/$GOFILE -package=mock$GOPACKAGE . Service

type SignRequest struct {
	SignerID int64
}

type SignResponse struct {
	SignatureID int64
}

type Service interface {
	Sign(ctx context.Context, req *SignRequest) (*SignResponse, error)
}
