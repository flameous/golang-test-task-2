package task

type Hotel struct {
	Id             uint64          `json:"id"`
	TitleRu        string          `json:"title_ru"`
	TitleEn        string          `json:"title_en"`
	Description    string          `json:"description"`
	Location       Location        `json:"location"`
	Address        string          `json:"address"`
	Rooms          []Room          `json:"rooms"`
	Amenities      []string        `json:"amenities"`
	AvailableRooms []AvailableRoom `json:"available_rooms"`
}

type HotelSearch struct {
	Id      uint64 `json:"id"`
	TitleRu string `json:"title_ru,omitempty"`
	TitleEn string `json:"title_en,omitempty"`
	Title   string `json:"title"`
}

func (h *HotelSearch) prepare() *HotelSearch {
	h.Title = `rus: ` + h.TitleRu + `, eng: ` + h.TitleEn
	h.TitleRu = ``
	h.TitleEn = ``
	return h
}

type HotelSearchSlice []*HotelSearch

func (p HotelSearchSlice) Len() int           { return len(p) }
func (p HotelSearchSlice) Less(i, j int) bool { return p[i].Title < p[j].Title }
func (p HotelSearchSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

type Location struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type Room struct {
	Id         uint64   `json:"id"`
	Title      string   `json:"title"`
	Type       string   `json:"type"`
	Beds       int      `json:"beds"`
	Facilities []string `json:"facilities"`
	Capacity   int      `json:"capacity"`
}

type AvailableRoom struct {
	RoomId uint64 `json:"room_id"`
	Count  int    `json:"count"`
}
