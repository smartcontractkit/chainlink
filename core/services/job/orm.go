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
	Exec(sql string, values ...interface{}) *gorm.DB
	First(out interface{}, where ...interface{}) *gorm.DB
	Create(value interface{}) *gorm.DB
	Transaction(fn func(db *gorm.DB) error) error
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
