package task

//func init() {
//	log.SetFlags(log.Lshortfile | log.Lmicroseconds)
//	ctx := context.Background()
//	client, err := elastic.NewClient()
//	if err != nil {
//		log.Fatal(err)
//	}
//	s = Server{client, ctx}
//
//	mappingFile, err := os.Open(`mapping.json`)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer mappingFile.Close()
//	b, err := ioutil.ReadAll(mappingFile)
//	if err != nil {
//		log.Fatal(err)
//	}
//	mapping := string(b)
//
//	ok, err := s.c.IndexExists(`hotels`).Do(s.ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	if !ok {
//		_, err = client.CreateIndex("hotels").BodyString(mapping).Do(s.ctx)
//		if err != nil {
//			log.Fatal(err)
//		}
//	}
//
//	dataFile, err := os.Open(`data.json`)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer dataFile.Close()
//	b, err = ioutil.ReadAll(dataFile)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	var hotels []*Hotel
//	if err = json.Unmarshal(b, &hotels); err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(hotels[0])
//	os.Exit(0)
//	for _, v := range hotels {
//		_, err = client.Index().Index(`hotels`).Type(`hotel`).Id(`321`).BodyJson(v).Do(ctx)
//		if err != nil {
//			log.Fatal(err)
//		}
//	}
//	log.Println(`init ok!`)
//}

type Hotel struct {
	Id             uint64          `json:"id"`
	Title          string          `json:"title"`
	Description    string          `json:"description"`
	Location       Location        `json:"location"`
	Address        string          `json:"address"`
	Rooms          []Room          `json:"rooms"`
	Amenities      []string        `json:"amenities"`
	AvailableRooms []AvailableRoom `json:"available_rooms"`
}

type Location struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}

type Room struct {
	Id       uint64 `json:"id"`
	Title    string `json:"title"`
	Type     string `json:"type"`
	Beds     int    `json:"beds"`
	Capacity int    `json:"capacity"`
}

type AvailableRoom struct {
	Id       uint64 `json:"id"`
	Capacity int    `json:"capacity"`
}
