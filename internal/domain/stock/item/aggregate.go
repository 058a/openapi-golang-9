package item

type (
	Aggregate struct {
		Id      Id
		Name    Name
		deleted bool
	}
)

func NewAggregate(id Id, name Name) *Aggregate {
	return &Aggregate{
		Id:      id,
		Name:    name,
		deleted: false,
	}
}

func RestoreAggregate(id Id, name Name, deleted bool) *Aggregate {
	return &Aggregate{
		Id:      id,
		Name:    name,
		deleted: deleted,
	}
}

func (a Aggregate) IsDeleted() bool {
	return a.deleted
}

func (a *Aggregate) Delete() {
	a.deleted = true
}
