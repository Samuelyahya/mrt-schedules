package station

type Station struct {
	ID string `json:"nid"`
	Name string `json:"title"`
}

type StationResponse struct {
	ID string `json:"id"`
	Name string `json:"name"`
}