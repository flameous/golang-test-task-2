# Test task for Golang developer vacancy

### run:
```bash
pip3 install -r parser/requirements.txt
python3 parser/parse_hotel_data.py _link_
go run cmd/main.go --path=`pwd`
go run cmd/tests.go
```

### implemented:
1. fuzzy search.

    params: _title_ (string)

    `http://localhost:8080/searchHotels?title=FooBar`


2. search by geo.

    params: _geo_ (float,float), _radius_ (float)

    `http://localhost:8080/searchHotels?geo=55.34,-4.32&radius=100`


3. get info by id.

    params: _/id_ (int) - in path

    `http://localhost:8080/hotels/1`
