package marvelch

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"marvel-chars/internal/service"
	"marvel-chars/internal/utils"
	"net/http"
	"os"
	"strconv"
	"unsafe"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/google/uuid"
)

// MarvelAPIRsp response data to /characters/{id}
type MarvelAPIRsp struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Desc string `json:"description"`
}

// MarvelAPIReqParams is for storing parameters needed for
// the marvel api call
type MarvelAPIReqParams struct {
	ts     string
	pbkey  string
	hash   string
	offset int
	limit  int
	url    string
}

// MarvelAPISvcStats contains service stats
// total - number of characters available from API call
// ctr   - number of characters fetched
// sz    - size in bytes of the entire character set
type MarvelAPISvcStats struct {
	total int
	ctr   int
	sz    uintptr
}

// MarvelCharSvc service
type MarvelCharSvc struct {
	Svc    service.Service
	Params *MarvelAPIReqParams
	Stats  *MarvelAPISvcStats
}

// Start the marvel characters service
func (m MarvelCharSvc) Start(msg string) error {
	b := service.Builder{}
	m.Svc = b.
		LoadEnv().
		Cache(m.setCache).
		Router(m.setRoutes).
		Build()
	return m.Svc.Start(msg)
}

// SetRoutes is the Router init function
// that uses go-chi router
func (m *MarvelCharSvc) setRoutes(smux service.Mux) {
	mux := smux.(*chi.Mux)
	mux.Use(middleware.Logger)
	mux.Route("/characters", func(r chi.Router) {
		r.Get("/", m.handleGetCharacterIDs)
		r.Get("/{id}", m.handleGetCharacterByID)
	})
}

func (m *MarvelCharSvc) handleGetCharacterIDs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ids, err := GetCharacterIDs(m.Svc.SvcCache.Cache)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	ret, err := json.Marshal(ids)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(ret)
}

func (m *MarvelCharSvc) handleGetCharacterByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := chi.URLParam(r, "id")
	log.Printf("%s", id)
	rsp := MarvelAPIRsp{}
	rsp = GetCharacterByID(m.Svc.SvcCache.Cache, id)
	ret, err := json.Marshal(&rsp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(ret)
}

func (m *MarvelCharSvc) setCache(c service.Cache) error {
	log.Printf("setting cache...\n")
	err := setSvcCache(c, m.Params)
	if err != nil {
		return err
	}
	return nil
}

func setSvcCache(c service.Cache, params *MarvelAPIReqParams) error {
	var err error
	log.Printf("creating API params...")
	params, err = newMarvelAPIReqParams()
	if err != nil {
		return fmt.Errorf("newMarvelAPIReqParams: %w", err)
	}
	log.Printf("done!\n")
	resp := make(map[string]interface{})
	stats := &MarvelAPISvcStats{}

	log.Printf("starting to fetch...\n")
	for {
		log.Printf("fetching...\n")
		if err = fetchDataFromMarvelAPI(*params, resp, stats); err != nil {
			return fmt.Errorf("fetchDataFromMarvelAPI: %w", err)
		}

		i, err := saveMarvelAPIResults(c, resp, stats)
		if err != nil {
			return fmt.Errorf("saveResults: %w", err)
		}
		params.offset += i

		log.Printf("stats.total %d | stats.ctr %d | stats.sz %d bytes\n", stats.total, stats.ctr, stats.sz)
		if stats.total == stats.ctr {
			break
		}
	}
	return nil
}

func newMarvelAPIReqParams() (*MarvelAPIReqParams, error) {
	var err error
	apiParams := MarvelAPIReqParams{
		ts:    uuid.New().String(),
		pbkey: os.Getenv("PUBKEY"),
		url:   os.Getenv("API_URL"),
	}
	apiParams.hash = utils.GetAPIKeyHash(apiParams.ts, os.Getenv("PRIVKEY"), apiParams.pbkey)
	apiParams.offset, err = strconv.Atoi(os.Getenv("API_INIT_OFFSET"))
	if err != nil {
		return nil, fmt.Errorf("failed to Atoi API_INIT_OFFSET: %w", err)
	}
	apiParams.limit, err = strconv.Atoi(os.Getenv("API_LIMIT"))
	if err != nil {
		return nil, fmt.Errorf("failed to Atoi API_LIMIT: %w", err)
	}

	return &apiParams, nil
}

func fetchDataFromMarvelAPI(a MarvelAPIReqParams, rsp map[string]interface{}, s *MarvelAPISvcStats) error {
	b, err := callMarvelAPI(a.ts, a.pbkey, a.hash, strconv.Itoa(a.offset), strconv.Itoa(a.limit), a.url)
	if err != nil {
		return fmt.Errorf("failed callMarvelAPI: %w", err)
	}
	// update estimated cache size
	s.sz += (unsafe.Sizeof(b) * uintptr(len(b)))

	err = json.Unmarshal(b, &rsp)
	if err != nil {
		return fmt.Errorf("failed to unmarshal body: %w", err)
	}

	return nil
}

func callMarvelAPI(ts, pbkey, hash, offset, limit, url string) ([]byte, error) {
	paramStr := fmt.Sprintf("ts=%s&apikey=%s&hash=%s&offset=%s&limit=%s", ts, pbkey, hash, offset, limit)
	rsp, err := http.Get(url + paramStr)
	if err != nil {
		return nil, fmt.Errorf("failed http request: %w", err)
	}

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}
	rsp.Body.Close()

	return body, nil
}

func saveMarvelAPIResults(c service.Cache, r map[string]interface{}, s *MarvelAPISvcStats) (int, error) {
	data := r["data"].(map[string]interface{})
	results := data["results"].([]interface{})

	for _, result := range results {
		r := result.(map[string]interface{})
		err := c.Set(strconv.Itoa(int(r["id"].(float64))), r)
		if err != nil {
			return 0, fmt.Errorf("failed to set cache: %w", err)
		}
	}

	count := int(data["count"].(float64))
	s.total = int(data["total"].(float64))
	s.ctr += count
	return count, nil
}
