package task

import (
	`log`
	`io/ioutil`
	`os`
	`github.com/olivere/elastic`
	`context`
	`encoding/json`
)

func InitElasticClient(path string) (context.Context, *elastic.Client) {
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

	dataFile, err := os.Open(`data.json`)
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
	//for _, v := range hotels {
	//	_, err = client.Index().Index(`hotels`).Type(`hotel`).Id(``).BodyJson(v).Do(ctx)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//}

	log.Println(`init elastic client - ok!`)
	return ctx, client
}
