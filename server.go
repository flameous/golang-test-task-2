package task

import (
	"fmt"
	"reflect"
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic"
	"context"
	"net/http"
	"strconv"
)

type Server struct {
	c   *elastic.Client
	ctx context.Context
}

func NewServer(ctx context.Context, c *elastic.Client) *Server {
	return &Server{c, ctx}
}

func (s *Server) searchHotels(c *gin.Context) {
	var err error

	var title, geo, radius string
	title, ok := c.GetQuery(`title`)
	if !ok {
		geo, ok = c.GetQuery(`geo`)
		radius = c.DefaultQuery(`radius`, `50km`)
		if !ok {
			c.String(http.StatusBadRequest, `missing field "title" or "geo"`)
			return
		}
	}

	ids, err := s.getHotelsByTitle(title, geo, radius)
	if len(ids) == 0 {
		c.String(http.StatusOK, `not found`)
		return
	}

	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, ids)
}

func (s *Server) getHotel(c *gin.Context) {
	id, err := strconv.Atoi(c.Param(`id`))
	if err != nil {
		c.String(http.StatusBadRequest, `invalid id: error - `+err.Error())
		return
	}
	h, err := s.getHotelsById(id)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	if h == nil {
		c.String(http.StatusOK, `not found`)
		return
	}
	c.JSON(http.StatusOK, h)
}

func (s *Server) getHotelsByTitle(title, geo, radius string) ([]*Hotel, error) {
	mmq := elastic.NewMultiMatchQuery(title, []string{`title.ru`, `title.en`}...)

	res, err := s.c.Search(`hotels`).Query(mmq).Do(s.ctx)
	if err != nil {
		return nil, fmt.Errorf(`elastic: multi-match error: %v`, err)
	}

	for _, v := range res.Each(reflect.TypeOf(`todo`)) {
		panic(v)
	}
	return nil, nil
}

func (s *Server) getHotelsByGeo(geo, radius string) ([]*Hotel, error) {
	return nil, nil
}

func (s *Server) getHotelsById(id int) (*Hotel, error) {
	return nil, nil
}

func (s *Server) Run() {
	router := gin.Default()
	router.GET(`/searchHotels`, s.searchHotels)
	router.GET(`/hotels/:id`, s.getHotel)
	router.Run(":8080")
}
