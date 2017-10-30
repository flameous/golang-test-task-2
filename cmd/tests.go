package main

import (
	"net/http"
	"net/url"
	"log"
	"io/ioutil"
	"fmt"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Lmicroseconds)
	host := `localhost:8080`
	searchByTitle(host, `альдего`)
	searchByTitle(host, `Aldego`)
	searchByTitle(host, `альде`)
	searchByTitle(host, `такого отеля нет!`)

	searchByGeo(host, `53.6,58.5`, `100`)
	searchByGeo(host, `53.6,58.5`, `1`)
	searchByGeo(host, `-20,70`, `1000`)

	getHotelById(host, "1")
}

func searchByTitle(host, title string) {
	v := url.Values{}
	v.Set(`title`, title)
	link := url.URL{
		Scheme:   `http`,
		Host:     host,
		Path:     `searchHotels`,
		RawQuery: v.Encode(),
	}

	resp, err := http.Get(link.String())
	if err != nil {
		log.Println(err)
		return
	}
	if resp != nil {
		defer resp.Body.Close()
	}
	b, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf(`search by title: "%s" -- data: %s`+"\n", title, b)
}

func searchByGeo(host, geo, radius string) {
	v := url.Values{}
	v.Set(`geo`, geo)
	if radius != `` {
		v.Set(`radius`, radius)
	}
	link := url.URL{
		Scheme:   `http`,
		Host:     host,
		Path:     `searchHotels`,
		RawQuery: v.Encode(),
	}

	resp, err := http.Get(link.String())
	if err != nil {
		log.Println(err)
		return
	}
	if resp != nil {
		defer resp.Body.Close()
	}
	b, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf(`search by geo: "%s" (radius: %s km) -- data: %s`+"\n", geo, radius, b)
}

func getHotelById(host, id string) {
	resp, err := http.Get( `http://`+ host + `/hotels/` + id)
	if err != nil {
		log.Println(err)
		return
	}
	if resp != nil {
		defer resp.Body.Close()
	}
	b, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf(`get hotel by id (%s) -- data: %s`+"\n", id, b)
}
