package keys

//go:generate mockery --name Keystore --output ./mocks/ --case=underscore --filename keystore.go

type Keystore interface {
	Get(id string) (Key, error)
}
