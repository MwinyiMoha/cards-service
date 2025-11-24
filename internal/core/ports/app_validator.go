package ports

type AppValidator interface {
	Struct(s any) error
}
