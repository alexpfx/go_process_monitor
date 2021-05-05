package psmon

type Source interface {
	Generate() (chan string, error)
}