package item

import (
	"errors"

	"github.com/google/uuid"

	"openapi/internal/domain/stock/item"
)

type (
	createRequest struct {
		Id   item.Id
		Name item.Name
	}
	createResponse struct {
		Id uuid.UUID
	}
)

var (
	ErrValidation = errors.New("NewCreateRequest: validation error")
)

func NewCreateRequest(id uuid.UUID, name string) (createRequest, error) {
	// validation
	validId, err := item.NewId(id)
	if err != nil {
		return createRequest{}, errors.Join(ErrValidation, err)
	}

	validName, err := item.NewName(name)
	if err != nil {
		return createRequest{}, errors.Join(ErrValidation, err)
	}

	// post processing	
	return createRequest{
		Id:   validId,
		Name: validName,
	}, nil
}

func NewCreateResponse(id item.Id) createResponse {
	return createResponse{
		Id: id.UUID(),
	}
}

func Create(req createRequest, r item.IRepository) (createResponse, error) {
	// Preprocessing
	a := item.NewAggregate(req.Id, req.Name)

	// Main
	if err := r.Save(a); err != nil {
		return createResponse{}, err
	}

	// Postprocessing
	res := NewCreateResponse(a.Id)
	return res, nil
}
