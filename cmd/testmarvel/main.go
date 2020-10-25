package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"marvel-chars/internal/utils"
	"net/http"
	"os"
	"strconv"
	"unsafe"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

var chars = make(map[string]interface{})

func init() {
	godotenv.Load()
}

func main() {
	ts := uuid.New().String()
	privateKey := os.Getenv("PRIVKEY")
	publicKey := os.Getenv("PUBKEY")
	log.Printf("privateKey: [%s]\n", privateKey)
	log.Printf("privatpublicKeyeKey: [%s]\n", publicKey)
	url := "http://gateway.marvel.com/v1/public/characters?"
	hash := utils.GetAPIKeyHash(ts, privateKey, publicKey)
	log.Printf("hash: [%s]\n", hash)

	paramStr := fmt.Sprintf("ts=%s&apikey=%s&hash=%s&offset=%s&limit=%s", ts, publicKey, hash, "0", "100")
	log.Printf("paramStr: [%s]\n", paramStr)
	log.Printf("geturl: [%s]\n", url+paramStr)
	rsp, err := http.Get(url + paramStr)
	if err != nil {
		log.Fatalf("failed http request: %s\n", err)
	}

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		log.Fatalf("failed to read body: %s\n", err)
	}
	rsp.Body.Close()
	/*
		saveTo, err := os.Create(os.ExpandEnv("resp.json"))
		if err != nil {
			log.Fatalln("Cannot create", "resp.json")
		}
		defer saveTo.Close()

		saveTo.Write(body)
	*/
	resp := make(map[string]interface{})

	err = json.Unmarshal(body, &resp)
	if err != nil {
		log.Fatalf("failed to unmarshal body: %s\n", err)
	}

	data := resp["data"].(map[string]interface{})

	// Initial markers
	offset := int(data["offset"].(float64)) // 0
	total := int(data["total"].(float64))   // 1493
	ctr := int(data["count"].(float64))     // 100
	limit := int(data["limit"].(float64))   // 100

	// initial set
	results := data["results"].([]interface{})
	for _, result := range results {
		r := result.(map[string]interface{})
		chars[strconv.Itoa(int(r["id"].(float64)))] = r
		chars[r["name"].(string)] = r
	}

	// rest of the set
	// set new starting offset
	offset += ctr
	max := total / limit
	log.Printf("max: %d\n", max)
	for i := 0; i < max; i++ {
		// call api
		paramStr = fmt.Sprintf("ts=%s&apikey=%s&hash=%s&offset=%s&limit=%s", ts, publicKey, hash, strconv.Itoa(offset), "100")
		rsp, err := http.Get(url + paramStr)
		if err != nil {
			log.Fatalf("failed http request: %s\n", err)
		}
		body, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			log.Fatalf("failed to read body: %s\n", err)
		}
		rsp.Body.Close()
		// store results
		err = json.Unmarshal(body, &resp)
		if err != nil {
			log.Fatalf("failed to unmarshal body: %s\n", err)
		}
		data := resp["data"].(map[string]interface{})
		results := data["results"].([]interface{})
		for _, result := range results {
			r := result.(map[string]interface{})
			chars[strconv.Itoa(int(r["id"].(float64)))] = r
			chars[r["name"].(string)] = r
		}
		// set new offset to current offset + count
		offset += int(data["count"].(float64))
		// update ctr to ctr + count
		ctr += int(data["count"].(float64))
	}

	log.Printf("offset: %d\n", offset)
	log.Printf("count: %d\n", ctr)
	log.Printf("total: %d\n", total)
	log.Printf("chars len: %d\n", len(chars))
	log.Printf("sizeof chars: %d\n", unsafe.Sizeof(chars))
	// get 1009144
	log.Printf("character %s: %+v\n", "1009144", chars["1009144"])

	// find ironman
	log.Printf("Iron Man: %+v\n", chars["Iron Man"])

	// delete map elements
	char := chars["1009368"].(map[string]interface{})
	name := char["name"].(string)
	delete(chars, "1009368")
	delete(chars, name)
	log.Printf("Iron Man: %+v\n", chars["1009368"])
	log.Printf("Iron Man: %+v\n", chars["Iron Man"])
}
