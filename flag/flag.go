package flag

import (
	"encoding/json"
	"ftxt-3-3/model"
	"net/http"

	"github.com/hashicorp/go-memdb"
)

type flagHandler struct {
	db *memdb.MemDB
}

func NewFlagHandler(db *memdb.MemDB) *flagHandler {
	return &flagHandler{db}
}

type Body struct {
	Flag string `json:"flag"`
}

func (fh *flagHandler) PutFlag(w http.ResponseWriter, r *http.Request) {
	var b Body
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		w.Write([]byte("error occured: " + err.Error()))
		return
	}

	txn := fh.db.Txn(true)

	flg := model.Flag{Flag: b.Flag}

	if err := txn.Insert("flag", &flg); err != nil {
		w.Write([]byte("error occured: " + err.Error()))
	}
	txn.Commit()
}

func (fh *flagHandler) GetFlag(w http.ResponseWriter, r *http.Request) {
	txn := fh.db.Txn(false)
	raw, err := txn.First("flag", "id")
	if err != nil {
		w.Write([]byte("error occured: " + err.Error()))
	}

	flg := raw.(*model.Flag)

	bytes, err := json.Marshal(flg)
	if err != nil {
		w.Write([]byte("error occured: " + err.Error()))
	}
	w.Write(bytes)
}
