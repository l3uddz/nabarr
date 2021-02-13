package omdb

type rating struct {
	Source string `json:"Source,omitempty"`
	Value  string `json:"Value,omitempty"`
}

type lookupResponse struct {
	Title      string   `json:"Title,omitempty"`
	Year       string   `json:"Year,omitempty"`
	Rated      string   `json:"Rated,omitempty"`
	Released   string   `json:"Released,omitempty"`
	Runtime    string   `json:"Runtime,omitempty"`
	Genre      string   `json:"Genre,omitempty"`
	Director   string   `json:"Director,omitempty"`
	Writer     string   `json:"Writer,omitempty"`
	Actors     string   `json:"Actors,omitempty"`
	Plot       string   `json:"Plot,omitempty"`
	Language   string   `json:"Language,omitempty"`
	Country    string   `json:"Country,omitempty"`
	Awards     string   `json:"Awards,omitempty"`
	Poster     string   `json:"Poster,omitempty"`
	Ratings    []rating `json:"Ratings,omitempty"`
	Metascore  string   `json:"Metascore,omitempty"`
	ImdbRating string   `json:"imdbRating,omitempty"`
	ImdbVotes  string   `json:"imdbVotes,omitempty"`
	ImdbID     string   `json:"imdbID,omitempty"`
	Type       string   `json:"Type,omitempty"`
	DVD        string   `json:"DVD,omitempty"`
	BoxOffice  string   `json:"BoxOffice,omitempty"`
	Production string   `json:"Production,omitempty"`
	Website    string   `json:"Website,omitempty"`
	Response   string   `json:"Response,omitempty"`
}

type Item struct {
	Metascore      int     `json:"metascore,omitempty"`
	RottenTomatoes int     `json:"rotten_tomatoes,omitempty"`
	ImdbRating     float64 `json:"imdb_rating,omitempty"`
}
