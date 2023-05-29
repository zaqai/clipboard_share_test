package main

import (
	"encoding/json"
	"fmt"
	"github.com/nutsdb/nutsdb"
	"log"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/postq", pushData)
	http.HandleFunc("/idxq", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})
	http.HandleFunc("/", pullData)

	s := &http.Server{
		Addr:           ":9090",
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	err := s.ListenAndServe()
	if err != nil {
		log.Fatal("listenAndServe: ", err)
	}

}

type ReqData struct {
	Key   string `json:"key"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

func pushData(w http.ResponseWriter, r *http.Request) {
	var reqdata ReqData
	// 调用json包的解析，解析请求body
	if err := json.NewDecoder(r.Body).Decode(&reqdata); err != nil {
		r.Body.Close()
		log.Println(err, r.Body)
	}

	DBKey := reqdata.Key
	_value := reqdata.Value
	_type := reqdata.Type
	//m := make(map[string]string)
	DBValue := `{"value":"` + _value + `","type":"` + _type + `"}`
	log.Println("write: ", DBKey, DBValue)
	writeDB(DBKey, DBValue)
	result := readDB(DBKey)
	m := make(map[string]string)
	err := json.Unmarshal(result, &m)
	if err != nil {
		fmt.Println("json unmarshal error:", err)
	}
	// 返回json字符串给客户端
	w.Header().Set("Content-Type", "text/plain;charset=utf-8")
	w.Write([]byte(m["value"]))
}
func pullData(w http.ResponseWriter, r *http.Request) {

	DBKey := r.URL.Path[1:]
	DBValue := readDB(DBKey)

	m := make(map[string]string)
	err := json.Unmarshal(DBValue, &m)
	if err != nil {
		fmt.Println("json unmarshal error:", err)
	}
	log.Println("read: ", DBKey, m)
	if m["type"] == "text" {
		w.Header().Set("Content-Type", "text/plain;charset=utf-8")
		w.Write([]byte(m["value"]))
	} else {
		http.Redirect(w, r, m["value"], http.StatusFound)
	}

}

func writeDB(k string, v string) {
	opt := nutsdb.DefaultOptions
	opt.Dir = "./nutsdb"
	db, err := nutsdb.Open(opt)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Update(func(tx *nutsdb.Tx) error {
		key := []byte(k)
		val := []byte(v)
		if err := tx.Put("", key, val, 0); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Println(err, k, v)
	}
}
func readDB(k string) []byte {
	opt := nutsdb.DefaultOptions
	opt.Dir = "./nutsdb"
	db, err := nutsdb.Open(opt)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	var value []byte
	err = db.View(func(tx *nutsdb.Tx) error {
		key := []byte(k)
		e, err := tx.Get("", key)
		if err != nil {
			return err
		} else {
			value = e.Value
		}
		return nil
	})
	if err != nil {
		log.Println(err, k)
	}
	return value
}
