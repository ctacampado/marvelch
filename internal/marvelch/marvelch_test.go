package marvelch

import (
	"encoding/json"
	"fmt"
	"marvel-chars/internal/kvstore"
	"net/http"
	"strconv"
	"testing"

	"github.com/joho/godotenv"
)

func TestCallMarvelAPI(t *testing.T) {
	tc := struct {
		code int
	}{
		code: http.StatusOK,
	}

	godotenv.Load("../../.env")
	p, err := newMarvelAPIReqParams()
	if err != nil {
		t.Errorf(err.Error())
	}

	b, err := callMarvelAPI(p.ts,
		p.pbkey,
		p.hash,
		strconv.Itoa(p.offset),
		strconv.Itoa(p.limit),
		p.url)
	if err != nil {
		t.Errorf(err.Error())
	}

	rsp := make(map[string]interface{})
	err = json.Unmarshal(b, &rsp)
	if err != nil {
		t.Errorf(err.Error())
	}

	want := tc.code
	got := int(rsp["code"].(float64))
	if want != got {
		t.Errorf("want %+v got %+v | fail\n", want, got)
	}
}

func TestSaveToCache(t *testing.T) {
	godotenv.Load("../../.env")

	fmt.Printf("creating cache...\n")
	c := kvstore.New()

	fmt.Printf("seting cache...\n")
	params := &MarvelAPIReqParams{}
	err := setSvcCache(c, params)
	if err != nil {
		t.Errorf(err.Error())
	}

	fmt.Printf("fetching from cache...\n")
	data := c.Get("1009368")
	if data == nil {
		t.Errorf("data is empty")
	}

	char := data.(map[string]interface{})
	got := char["name"]
	want := "Iron Man"

	if want != got {
		t.Errorf("want %s got %s | fail\ndata %+v\n char %+v", want, got, data, char)
	}
}
