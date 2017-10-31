package task

import (
	"context"
	"encoding/json"
	"github.com/olivere/elastic"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"fmt"
	"strings"
	"reflect"
)

type ElasticClient struct {
	c   *elastic.Client
	ctx context.Context
}

func NewtElasticClient(path string) *ElasticClient {
	if path[len(path)-1] != '/' {
		path += "/"
	}
	mappingFile, err := os.Open(path + `mapping.json`)
	if err != nil {
		log.Fatal(err)
	}
	defer mappingFile.Close()
	b, err := ioutil.ReadAll(mappingFile)
	if err != nil {
		log.Fatal(err)
	}
	mapping := string(b)

	ctx := context.Background()
	client, err := elastic.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	ok, err := client.IndexExists(`hotels`).Do(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if !ok {
		_, err = client.CreateIndex(`hotels`).BodyString(mapping).Do(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}

	dataFile, err := os.Open(path + `data.json`)
	if err != nil {
		log.Fatal(err)
	}
	defer dataFile.Close()
	b, err = ioutil.ReadAll(dataFile)
	if err != nil {
		log.Fatal(err)
	}

	var hotels []*Hotel
	if err = json.Unmarshal(b, &hotels); err != nil {
		log.Fatal(err)
	}
	for _, v := range hotels {
		_, err = client.Index().
			Index(`hotels`).
			Type(`hotel`).
			Id(strconv.FormatUint(v.Id, 10)).
			BodyJson(v).
			Do(ctx)

		if err != nil {
			log.Fatal(err)
		}
	}

	log.Println(`init es client - ok!`)
	return &ElasticClient{ctx: ctx, c: client}
}

func (e *ElasticClient) getHotelById(id int) (*Hotel, error) {
	result, err := e.c.Get().
		Index(`hotels`).
		Id(strconv.Itoa(id)).
		Do(e.ctx)
	if err != nil {
		if elastic.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	if result.Source == nil {
		return nil, fmt.Errorf(`result's source is nil: %#v`, result)
	}

	h := new(Hotel)
	if err := json.Unmarshal(*result.Source, h); err != nil {
		return nil, err
	}
	return h, nil
}

func (e *ElasticClient) searchHotelsByTitle(title string) ([]*HotelSearch, error) {
	q := elastic.NewMultiMatchQuery(title, []string{`title_ru`, `title_en`}...)

	titleIsRussian := isRussian(title)
	var s elastic.Sorter
	if titleIsRussian {
		s = elastic.NewFieldSort(`title_ru.raw`)
	} else {
		s = elastic.NewFieldSort(`title_en.raw`)
	}

	result, err := e.c.Search(`hotels`).
		Query(q).
		SortBy(s).
		Do(e.ctx)

	if err != nil {
		return nil, fmt.Errorf(`elastic: search error: %v`, err)
	}
	return parseHotels(result, titleIsRussian)
}

func (e *ElasticClient) searchHotelsByGeo(geo, radius string) ([]*HotelSearch, error) {
	arr := strings.Split(geo, ",")
	if len(arr) != 2 {
		return nil, fmt.Errorf(`"geo" field format must be "float,float"`)
	}

	lat, err := strconv.ParseFloat(strings.Trim(arr[0], " "), 64)
	if err != nil {
		return nil, fmt.Errorf(`"geo" field - expect float64, found: "%e". err: %v"`, arr[0], err)
	}
	long, err := strconv.ParseFloat(strings.Trim(arr[1], " "), 64)
	if err != nil {
		return nil, fmt.Errorf(`"geo" field - expect float64, found: "%e". err: %v"`, arr[1], err)
	}

	q := elastic.NewGeoDistanceQuery(`location`)
	q.Distance(radius + `km`)
	q.Lat(lat)
	q.Lon(long)

	// fixme: достаются только 10 hit'ов
	result, err := e.c.Search(`hotels`).
		Query(q).
		SortBy(elastic.NewFieldSort(`title_ru.raw`)).
		Do(e.ctx)

	if err != nil {
		return nil, fmt.Errorf(`elastic: search error: %v`, err)
	}
	return parseHotels(result, true)
}

func parseHotels(r *elastic.SearchResult, isRussian bool) ([]*HotelSearch, error) {
	var hotels []*HotelSearch
	for _, v := range r.Each(reflect.TypeOf(HotelSearch{})) {
		if h, ok := v.(HotelSearch); ok {
			hotels = append(hotels, h.prepare(isRussian))
		} else {
			return nil, fmt.Errorf(`cannot cast data to type "Hotel" - data: %#v`, v)
		}
	}
	return hotels, nil
}

func isRussian(title string) bool {
	if title == `` {
		return true
	}
	r := []rune(title)[0]
	return 'А' <= r && r <= 'я'
}
