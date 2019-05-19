package testutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/ryantking/marina/pkg/prisma"
)

type exportRes struct {
	Out struct {
		JSONElems interface{} `json:"jsonElements"`
		Size      int         `json:"size"`
	} `json:"out"`
	Cursor map[string]int `json:"cursor"`
	IsFull bool           `json:"isFull"`
}

var (
	client = prisma.New(nil)
	mus    = map[string]*sync.Mutex{}
	musL   sync.Mutex

	data  = map[string]map[string]interface{}{}
	dataL sync.Mutex
)

func lockTable(table string) {
	musL.Lock()
	defer musL.Unlock()

	mu, ok := mus[table]
	if !ok {
		mu = new(sync.Mutex)
		mus[table] = mu
	}
	mu.Lock()
}

func unlockTable(table string) {
	musL.Lock()
	defer musL.Unlock()
	mus[table].Unlock()
}

func Aquire(tables ...string) {
	dataL.Lock()
	defer dataL.Unlock()

	for _, table := range tables {
		lockTable(table)

		data[table] = map[string]interface{}{}
		for _, ft := range []string{"nodes", "lists", "relations"} {
			reqBody := fmt.Sprintf(`{"fileType": "%s", "cursor": {"table": 0, "row": 0, "field": 0, "array": 0}}`, ft)
			res, err := http.Post("http://localhost:4466/marina/dev/export", "application/json", strings.NewReader(reqBody))
			if err != nil {
				panic(err.Error())
			}
			if res.StatusCode != http.StatusOK {
				panic("bad response from prisma")
			}

			elems := new(exportRes)
			err = json.NewDecoder(res.Body).Decode(elems)
			if err != nil {
				panic(err.Error())
			}
			data[table][ft] = elems.Out.JSONElems
		}
	}
}

func Clean(tables ...string) {
	dataL.Lock()
	defer dataL.Unlock()

	for _, table := range tables {
		for _, ft := range []string{"nodes", "lists", "relations"} {
			reqBody := map[string]interface{}{}
			reqBody["valueType"] = ft
			reqBody["values"] = data[table][ft]
			b, err := json.Marshal(reqBody)
			if err != nil {
				panic(err.Error())
			}
			res, err := http.Post("http://localhost:4466/marina/dev/import", "application/json", bytes.NewReader(b))
			if err != nil {
				panic(err.Error())
			}
			if res.StatusCode != http.StatusOK {
				panic("bad response from prisma")
			}
		}

		delete(data, table)
		unlockTable(table)
	}
}
