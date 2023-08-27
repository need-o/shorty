package memory

type Storage struct {
	Shorty *ShortyStorage
}

func NewStorage() *Storage {
	return &Storage{
		Shorty: NewShortyStorage(),
	}
}
