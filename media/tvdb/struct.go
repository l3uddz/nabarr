package tvdb

type lookupResponse struct {
	Data struct {
		Id              int           `json:"id,omitempty"`
		SeriesId        string        `json:"seriesId,omitempty"`
		SeriesName      string        `json:"seriesName,omitempty"`
		Aliases         []interface{} `json:"aliases,omitempty"`
		Season          string        `json:"season,omitempty"`
		Poster          string        `json:"poster,omitempty"`
		Banner          string        `json:"banner,omitempty"`
		Fanart          string        `json:"fanart,omitempty"`
		Status          string        `json:"status,omitempty"`
		FirstAired      string        `json:"firstAired,omitempty"`
		Network         string        `json:"network,omitempty"`
		NetworkId       string        `json:"networkId,omitempty"`
		Runtime         string        `json:"runtime,omitempty"`
		Language        string        `json:"language,omitempty"`
		Genre           []string      `json:"genre,omitempty"`
		Overview        string        `json:"overview,omitempty"`
		LastUpdated     int           `json:"lastUpdated,omitempty"`
		AirsDayOfWeek   string        `json:"airsDayOfWeek,omitempty"`
		AirsTime        string        `json:"airsTime,omitempty"`
		Rating          interface{}   `json:"rating,omitempty"`
		ImdbId          string        `json:"imdbId,omitempty"`
		Zap2ItId        string        `json:"zap2itId,omitempty"`
		Added           string        `json:"added,omitempty"`
		AddedBy         int           `json:"addedBy,omitempty"`
		SiteRating      float64       `json:"siteRating,omitempty"`
		SiteRatingCount int           `json:"siteRatingCount,omitempty"`
		Slug            string        `json:"slug,omitempty"`
	} `json:"data"`
}

type Item struct {
	Runtime         int      `json:"Runtime,omitempty"`
	Language        string   `json:"Language,omitempty"`
	Network         string   `json:"Network,omitempty"`
	Genre           []string `json:"Genre,omitempty"`
	AirsDayOfWeek   string   `json:"AirsDayOfWeek,omitempty"`
	SiteRating      float64  `json:"SiteRating,omitempty"`
	SiteRatingCount int      `json:"SiteRatingCount,omitempty"`
}
