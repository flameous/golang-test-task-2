package task

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic"
	"context"
	"net/http"
	"strconv"
	"strings"
	"reflect"
	"encoding/json"
)

type Server struct {
	c   *elastic.Client
	ctx context.Context
}

func NewServer(ctx context.Context, c *elastic.Client) *Server {
	return &Server{c, ctx}
}

func (s *Server) searchHotels(c *gin.Context) {
	var (
		title, geo, radius string
		hotels             []*HotelSearch
		errSearch          error
	)

	title, ok := c.GetQuery(`title`)
	if !ok {
		if geo, ok = c.GetQuery(`geo`); !ok {
			c.String(http.StatusBadRequest, `missing field "title" or "geo"`)
			return
		}
		if radius, ok = c.GetQuery(`radius`); !ok {
			c.String(http.StatusBadRequest, `missing field "radius"`)
			return
		}
		if _, err := strconv.Atoi(radius); err != nil {
			c.String(http.StatusBadRequest,
				fmt.Sprintf(`"radius" field - expect int, found: "%s". err: %v"`, radius, err))
			return
		}
	}

	// todo: сортировка по алфавиту
	if geo == `` {
		hotels, errSearch = s.searchHotelsByTitle(title)
	} else {
		hotels, errSearch = s.searchHotelsByGeo(geo, radius)
	}

	if errSearch != nil {
		c.String(http.StatusInternalServerError, errSearch.Error())
		return
	}
	if len(hotels) == 0 {
		c.String(http.StatusOK, `ничего  не  найдено`)
		return
	}
	c.IndentedJSON(http.StatusOK, hotels)
}

func (s *Server) searchHotelsByTitle(title string) ([]*HotelSearch, error) {
	q := elastic.NewMultiMatchQuery(title, []string{`title_ru`, `title_en`}...)
	q.Fuzziness(`5`)
	result, err := s.c.Search(`hotels`).
		Query(q).
		Do(s.ctx)

	if err != nil {
		return nil, fmt.Errorf(`elastic: multi-match error: %v`, err)
	}
	return parseHotels(result)
}

func (s *Server) searchHotelsByGeo(geo, radius string) ([]*HotelSearch, error) {
	arr := strings.Split(geo, ",")
	if len(arr) != 2 {
		return nil, fmt.Errorf(`"geo" field format must be "-50, 60"`)
	}

	lat, err := strconv.ParseFloat(arr[0], 64)
	if err != nil {
		return nil, fmt.Errorf(`"geo" field - expect float64, found: "%s". err: %v"`, arr[0], err)
	}
	long, err := strconv.ParseFloat(arr[1], 64)
	if err != nil {
		return nil, fmt.Errorf(`"geo" field - expect float64, found: "%s". err: %v"`, arr[1], err)
	}

	q := elastic.NewGeoDistanceQuery(`location`)
	q.Distance(radius + `km`)
	q.Lat(lat)
	q.Lon(long)

	result, err := s.c.Search(`hotels`).
		Query(q).
		Do(s.ctx)

	if err != nil {
		return nil, err
	}
	return parseHotels(result)
}

func parseHotels(r *elastic.SearchResult) ([]*HotelSearch, error) {
	var hotels []*HotelSearch
	for _, v := range r.Each(reflect.TypeOf(HotelSearch{})) {
		if h, ok := v.(HotelSearch); ok {
			hotels = append(hotels, &h)
		} else {
			return nil, fmt.Errorf(`cannot cast data to type "Hotel" - data: %#v`, v)
		}
	}
	return hotels, nil
}

func (s *Server) getHotel(c *gin.Context) {
	id, err := strconv.Atoi(c.Param(`id`))
	if err != nil {
		c.String(http.StatusBadRequest, `invalid id: error - `+err.Error())
		return
	}
	h, err := s.getHotelById(id)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	if h == nil {
		c.String(http.StatusOK, `ничего  не  найдено`)
		return
	}
	c.IndentedJSON(http.StatusOK, h)
}

func (s *Server) getHotelById(id int) (*Hotel, error) {
	result, err := s.c.Get().
		Index(`hotels`).
		Id(strconv.Itoa(id)).
		Do(s.ctx)
	if err != nil {
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

func (s *Server) Run() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.GET(`/searchHotels`, s.searchHotels)
	router.GET(`/hotels/:id`, s.getHotel)
	router.Run(":8080")
}
