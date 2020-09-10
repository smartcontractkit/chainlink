package job

//go:generate mockery --name ORM --output ./mocks/ --case=underscore

type ORM interface {
	UnclaimedJobs() ([]JobSpec, error)
	CreateSpec(spec Spec) error
}

type orm struct {
	db database
}

type database interface {
	Create(value interface{}) *gorm.DB
}

func NewORM(o database) *orm {
	return &orm{o}
}

func (o *orm) UnclaimedJobs() ([]Spec, error) {
	panic("unimplemented")
}

func (o *orm) CreateSpec(spec Spec) error {
	return o.db.Create(spec)
}
