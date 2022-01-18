package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "mysql-master"
	"net/http"
)

type data_karyawan struct {
	Id     string
	Nama   string
	Umur   int
	Posisi string
}

func koneksi() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:@tcp(localhost)/tugas16_golang")
	if err != nil {
		return nil, err
	}
	return db, nil
}

var data []data_karyawan

func main() {
	ambil_data()
	http.HandleFunc("/karyawan", ambil_karyawan)
	http.HandleFunc("/cari_karyawan", cari_karyawan)
	http.HandleFunc("/hapus_karyawan", hapus_karyawan)
	http.HandleFunc("/update_karyawan", update_karyawan)
	fmt.Println("menjalankan Web server")
	http.ListenAndServe(":8080", nil)
}

func ambil_karyawan(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "POST" {
		var result, err = json.Marshal(data)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(result)
		return
	}
	http.Error(w, "", http.StatusBadRequest)
}

// cari karyawan

func cari_karyawan(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "POST" {
		var idkaryawan = r.FormValue("Id")
		var result []byte
		var err error

		for _, each := range data {
			if each.Id == idkaryawan {
				result, err = json.Marshal(each)

				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.Write(result)
				return
			}
		}
		http.Error(w, "Karyawan Tidak Terdaftar", http.StatusBadRequest)
		return
	}
	http.Error(w, "", http.StatusBadRequest)
}

// Hapus Karyawan
func hapus_karyawan(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "aplication/json")
	db, err := koneksi()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()

	if r.Method == "POST" {
		var idkaryawan = r.FormValue("Id")
		var err error

		// hapus Data
		_, err = db.Query("delete from tb_karyawan where Id =? ", idkaryawan)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		http.Error(w, "Karyawan dihapus", http.StatusBadRequest)
		return
	}
	http.Error(w, "", http.StatusBadRequest)
}

// Update Karyawan
func update_karyawan(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "aplication/json")
	db, err := koneksi()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()

	if r.Method == "POST" {
		var idkaryawan = r.FormValue("Id")
		var namakaryawan = r.FormValue("Nama")
		var err error

		// hapus Data
		_, err = db.Query("update tb_karyawan set Nama = ? where Id = ?", namakaryawan, idkaryawan)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		http.Error(w, "Karyawan diupdate", http.StatusBadRequest)
		return
	}
	http.Error(w, "", http.StatusBadRequest)
}

func ambil_data() {
	db, err := koneksi()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()

	rows, err := db.Query("Select * from tb_karyawan")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer rows.Close()

	for rows.Next() {
		var each = data_karyawan{}
		var err = rows.Scan(&each.Id, &each.Nama, &each.Umur, &each.Posisi)

		if err != nil {
			fmt.Println(err.Error())
			return
		}
		data = append(data, each)
	}
	if err = rows.Err(); err != nil {
		fmt.Println(err.Error())
		return
	}

}
