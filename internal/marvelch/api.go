package marvelch

import (
	"log"
	"marvel-chars/internal/service"
)

// GetCharacterIDs return IDs of all marvel characters
func GetCharacterIDs(c service.Cache) ([]int64, error) {
	ret := make([]int64, 0)
	res := c.Get("ALL")

	data := res.(map[string]interface{})
	for _, v := range data {
		char := v.(map[string]interface{})
		ret = append(ret, int64(char["id"].(float64)))
	}
	log.Printf("%d\n", len(ret))
	return ret, nil
}

// GetCharacterByID returns marvel character
// given its ID
func GetCharacterByID(c service.Cache, id string) MarvelAPIRsp {
	res := c.Get(id)
	char := res.(map[string]interface{})
	return MarvelAPIRsp{
		ID:   int64(char["id"].(float64)),
		Name: char["name"].(string),
		Desc: char["description"].(string),
	}
}
