package location

type IRepository interface {
	Save(a *Aggregate) error
	Get(id Id) (*Aggregate, error)
	Find(id Id) (bool, error)
}
