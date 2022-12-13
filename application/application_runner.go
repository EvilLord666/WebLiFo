package application

type AppRunner interface {
	Start() (bool, error)
	Stop() (bool, error)
	Init() (bool, error)
}
