package models

import "github.com/globalsign/mgo/bson"

type Movies struct {
	Id          bson.ObjectId `bson:"_id" json:"id"`
	Name        string        `bson:"name" json:"name"`
	CoverImage  string        `bson:"cover_image" json:"cover_image"`
	Description string        `bson:"description" json:"description"`
}

const (
	db         = "Movies"
	collection = "MovieModel"
)

func (m *Movies) InsertMovie(movie Movies) error {
	return Insert(db, collection, movie)
}

func (m *Movies) FindAllMovies() ([]Movies, error) {
	var result []Movies
	err := FindAll(db, collection, nil, nil, &result)
	return result, err
}

func (m *Movies) FindMovieById(id string) (Movies, error) {
	var result Movies
	err := FindOne(db, collection, bson.M{"_id": bson.ObjectIdHex(id)}, nil, &result)
	return result, err
}

func (m *Movies) UpdateMovie(movie Movies) error {
	return Update(db, collection, bson.M{"_id": movie.Id}, movie)
}

func (m *Movies) RemoveMovie(id string) error {
	return Remove(db, collection, bson.M{"_id": bson.ObjectIdHex(id)})
}
