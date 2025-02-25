package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"crmcandidate/vendors/db"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/mattn/go-sqlite3"
)

var (
	router *chi.Mux
	dbConn *db.DB
)

func catch(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func main() {
	router = chi.NewRouter()
	router.Use(middleware.Recoverer)

	var err error
	dbConn, err = db.NewDB() // Используем NewDB вместо Connect
	catch(err)
	defer dbConn.Close() // Закрываем соединение с базой данных

	router.Use(ChangeMethod)
	router.Get("/", GetAllArticles)
	router.Post("/upload", UploadHandler)
	router.Get("/images/*", ServeImages)
	router.Route("/articles", func(r chi.Router) {
		r.Get("/", NewArticle)
		r.Post("/", CreateArticle)
		r.Route("/{articleID}", func(r chi.Router) {
			r.Use(ArticleCtx)
			r.Get("/", GetArticle)       // GET /articles/1234
			r.Put("/", UpdateArticle)    // PUT /articles/1234
			r.Delete("/", DeleteArticle) // DELETE /articles/1234
			r.Get("/edit", EditArticle)  // GET /articles/1234/edit
		})
	})

	err = http.ListenAndServe(":8005", router)
	catch(err)
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	const MAX_UPLOAD_SIZE = 10 << 20
	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)
	if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
		http.Error(w, "The uploaded file is too big. Please choose an file that's less than 10MB in size", http.StatusBadRequest)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer file.Close()

	// Create the uploads folder if it doesn't already exist
	err = os.MkdirAll("./images", os.ModePerm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create a new file in the uploads directory
	filename := fmt.Sprintf("/images/%d%s", time.Now().UnixNano(), filepath.Ext(fileHeader.Filename))
	dst, err := os.Create("." + filename)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer dst.Close()

	// Copy the uploaded file to  the specified destination
	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println(filename)
	response, _ := json.Marshal(map[string]string{"location": filename})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}

func ServeImages(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL)
	fs := http.StripPrefix("/images/", http.FileServer(http.Dir("./images")))
	fs.ServeHTTP(w, r)
}

func ChangeMethod(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			switch method := r.PostFormValue("_method"); method {
			case http.MethodPut:
				fallthrough
			case http.MethodPatch:
				fallthrough
			case http.MethodDelete:
				r.Method = method
			default:
			}
		}
		next.ServeHTTP(w, r)
	})
}

type contextKey string

const articleKey contextKey = "article"

func ArticleCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		articleID := chi.URLParam(r, "articleID")
		article, err := dbConn.GetArticle(articleID)
		if err != nil {
			fmt.Println(err)
			http.Error(w, http.StatusText(404), 404)
			return
		}
		ctx := context.WithValue(r.Context(), articleKey, article)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetAllArticles(w http.ResponseWriter, r *http.Request) {
	articles, err := dbConn.GetAllArticles()
	catch(err)

	t, err := template.ParseFiles("templates/base.html", "templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, articles)
	catch(err)
}

func NewArticle(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/base.html", "templates/new.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, nil)
	catch(err)
}

func CreateArticle(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	content := r.FormValue("content")
	article := &db.Article{
		Title:   title,
		Content: template.HTML(content),
	}

	err := dbConn.CreateArticle(article)
	catch(err)
	http.Redirect(w, r, "/", http.StatusFound)
}

func GetArticle(w http.ResponseWriter, r *http.Request) {
	article := r.Context().Value(articleKey).(*db.Article)
	t, err := template.ParseFiles("templates/base.html", "templates/article.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, article)
	catch(err)
}

func EditArticle(w http.ResponseWriter, r *http.Request) {
	article := r.Context().Value(articleKey).(*db.Article)

	t, err := template.ParseFiles("templates/base.html", "templates/edit.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, article)
	catch(err)
}

func UpdateArticle(w http.ResponseWriter, r *http.Request) {
	article := r.Context().Value(articleKey).(*db.Article)

	title := r.FormValue("title")
	content := r.FormValue("content")
	newArticle := &db.Article{
		Title:   title,
		Content: template.HTML(content),
	}
	fmt.Println(newArticle.Content)
	err := dbConn.UpdateArticle(strconv.Itoa(article.ID), newArticle)
	catch(err)
	http.Redirect(w, r, fmt.Sprintf("/articles/%d", article.ID), http.StatusFound)
}

func DeleteArticle(w http.ResponseWriter, r *http.Request) {
	article := r.Context().Value(articleKey).(*db.Article)
	err := dbConn.DeleteArticle(strconv.Itoa(article.ID))
	catch(err)

	http.Redirect(w, r, "/", http.StatusFound)
}
