package task

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Server struct {
	ElasticClient *ElasticClient
}

func NewServer(e *ElasticClient) *Server {
	return &Server{e}
}

func (s *Server) searchHotels(c *gin.Context) {
	var (
		title, geo, radius string
		hotels             HotelSearchSlice
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
		hotels, errSearch = s.ElasticClient.searchHotelsByTitle(title)
	} else {
		hotels, errSearch = s.ElasticClient.searchHotelsByGeo(geo, radius)
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

func (s *Server) getHotel(c *gin.Context) {
	id, err := strconv.Atoi(c.Param(`id`))
	if err != nil {
		c.String(http.StatusBadRequest, `invalid id: error - `+err.Error())
		return
	}
	h, err := s.ElasticClient.getHotelById(id)
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

func (s *Server) Run() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.GET(`/searchHotels`, s.searchHotels)
	router.GET(`/hotels/:id`, s.getHotel)
	router.Run(":8080")
}
