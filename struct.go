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


// не сделал ли я всё сложнее?
type HotelSearch struct {
	Id      uint64 `json:"id"`
	TitleRu string `json:"title_ru,omitempty"`
	TitleEn string `json:"title_en,omitempty"`
	Title   string `json:"title"`
}

func (h *HotelSearch) prepare(isRussian bool) *HotelSearch {
	if isRussian {
		h.Title = h.TitleRu
	} else {
		h.Title = h.TitleEn
	}
	h.TitleRu = ``
	h.TitleEn = ``
	return h
}

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
