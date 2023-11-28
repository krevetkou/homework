package services

type MoviesRepository interface {
}

type MoviesService struct {
	Storage MoviesRepository
}

func NewMovieService(storage MoviesRepository) MoviesService {
	return MoviesService{
		Storage: storage,
	}
}
