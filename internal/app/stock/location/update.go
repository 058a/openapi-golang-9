package location

import (
	"openapi/internal/domain/stock/location"

	"github.com/google/uuid"
)

type updateRequest struct {
	Id   location.Id
	Name location.Name
}

func NewUpdateRequest(id uuid.UUID, name string) (updateRequest, error) {
	// validation
	validId, err := location.NewId(id)
	if err != nil {
		return updateRequest{}, err
	}

	validName, err := location.NewName(name)
	if err != nil {
		return updateRequest{}, err
	}

	// post processing
	return updateRequest{
		Id:   validId,
		Name: validName,
	}, nil
}

func Update(req updateRequest, r location.IRepository) error {
	id := location.Id(req.Id)
	a, err := r.Get(id)
	if err != nil {
		return err
	}

	a.Name = req.Name

	err = r.Save(a)
	if err != nil {
		return err
	}

	return nil
}
