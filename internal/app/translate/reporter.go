package translate

type Reporter interface {
	Report(progress Progress)
}

type Progress struct {
	File  string
	Done  int
	Total int
}
