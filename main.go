// Informasi Database:
// Username: root
// Password: jarankeplese
// Database: form_db
// Alur API:
// 1. Server akan berjalan pada localhost di port 1111.
// 2. Endpoint "/" menampilkan form HTML untuk mengisi data.
// 3. Endpoint "/submit" menerima data dari form, melakukan validasi sederhana, dan menyimpannya ke database.
// 4. Jika data berhasil disimpan, pengguna akan diarahkan ke halaman sukses.
// 5. Jika terjadi kesalahan, pengguna akan diberi pesan error di halaman form.
// http://localhost:1111/

package main

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type FormData struct {
	ID            int    `json:"id"`
	GolonganDarah string `json:"golongan_darah"`
	GadoGado      bool   `json:"gado_gado"`
	Nama          string `json:"nama"`
	JenisKelamin  string `json:"jenis_kelamin"`
}

var db *sql.DB

func main() {

	var err error
	db, err = sql.Open("mysql", "root:jarankeplese@tcp(localhost:3306)/form_db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS form_data (
            id INT AUTO_INCREMENT PRIMARY KEY,
            golongan_darah VARCHAR(5),
            gado_gado BOOLEAN,
            nama VARCHAR(100),
            jenis_kelamin VARCHAR(20)
        )
    `)
	if err != nil {
		panic(err)
	}
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")

	r.GET("/", showForm)
	r.POST("/submit", submitForm)

	r.Run(":1111")
}

func showForm(c *gin.Context) {
	c.HTML(http.StatusOK, "form.html", nil)
}

func submitForm(c *gin.Context) {
	var formData FormData

	formData.GolonganDarah = c.PostForm("golongan_darah")
	formData.GadoGado = c.PostForm("gado_gado") == "on"
	formData.Nama = c.PostForm("nama")
	formData.JenisKelamin = c.PostForm("jenis_kelamin")

	//save

	_, err := db.Exec(`
        INSERT INTO form_data (golongan_darah, gado_gado, nama, jenis_kelamin)
        VALUES (?, ?, ?, ?)
    `, formData.GolonganDarah, formData.GadoGado, formData.Nama, formData.JenisKelamin)

	if err != nil {
		c.HTML(http.StatusInternalServerError, "form.html", gin.H{
			"error": "Gagal menyimpan data",
		})
		return
	}

	c.HTML(http.StatusOK, "sukses.html", nil)
}
