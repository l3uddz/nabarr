package sonarr

type systemStatus struct {
	Version string
}

type qualityProfile struct {
	Name string
	Id   int
}

type lookupRequest struct {
	Id         int    `json:"id,omitempty"`
	Title      string `json:"title"`
	TitleSlug  string `json:"titleSlug"`
	Year       int    `json:"year,omitempty"`
	TvdbId     int    `json:"tvdbId"`
	SeriesType string `json:"seriesType"`
}

type addRequest struct {
	Title            string     `json:"title"`
	TitleSlug        string     `json:"titleSlug"`
	Year             int        `json:"year"`
	QualityProfileId int        `json:"qualityProfileId"`
	Images           []string   `json:"images"`
	Tags             []string   `json:"tags"`
	Monitored        bool       `json:"monitored"`
	RootFolderPath   string     `json:"rootFolderPath"`
	AddOptions       addOptions `json:"addOptions"`
	Seasons          []string   `json:"seasons"`
	SeriesType       string     `json:"seriesType"`
	SeasonFolder     bool       `json:"seasonFolder"`
	TvdbId           int        `json:"tvdbId"`
}

type addOptions struct {
	SearchForMissingEpisodes   bool `json:"searchForMissingEpisodes"`
	IgnoreEpisodesWithFiles    bool `json:"ignoreEpisodesWithFiles"`
	IgnoreEpisodesWithoutFiles bool `json:"ignoreEpisodesWithoutFiles"`
}
