package radarr

type systemStatus struct {
	Version string
}

type qualityProfile struct {
	Name string
	Id   int
}

type lookupRequest struct {
	Id        int    `json:"id,omitempty"`
	Title     string `json:"title"`
	TitleSlug string `json:"titleSlug"`
	Year      int    `json:"year,omitempty"`
	ImdbId    string `json:"imdbId"`
	TmdbId    int    `json:"tmdbId"`
}

type addRequest struct {
	Title               string     `json:"title"`
	TitleSlug           string     `json:"titleSlug"`
	Year                int        `json:"year"`
	QualityProfileId    int        `json:"qualityProfileId"`
	Images              []string   `json:"images"`
	Monitored           bool       `json:"monitored"`
	RootFolderPath      string     `json:"rootFolderPath"`
	MinimumAvailability string     `json:"minimumAvailability"`
	AddOptions          addOptions `json:"addOptions"`
	TmdbId              int        `json:"tmdbId,omitempty"`
	ImdbId              string     `json:"imdbId,omitempty"`
}

type addOptions struct {
	SearchForMovie             bool `json:"searchForMovie"`
	IgnoreEpisodesWithFiles    bool `json:"ignoreEpisodesWithFiles"`
	IgnoreEpisodesWithoutFiles bool `json:"ignoreEpisodesWithoutFiles"`
}
