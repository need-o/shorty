package memory

type Storage struct {
	Shortenings *ShorteningStorage
}

func NewStorage() *Storage {
	return &Storage{
		Shortenings: NewShorteningStorage(),
	}
}
