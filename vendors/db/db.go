package db

import (
	"database/sql"
	"html/template"

	_ "github.com/mattn/go-sqlite3"
)

type Article struct {
	ID      int
	Title   string
	Content template.HTML // Изменен тип на template.HTML
}

type DB struct {
	db *sql.DB
}

func NewDB() (*DB, error) {
	db, err := sql.Open("sqlite3", "./data.sqlite")
	if err != nil {
		return nil, err
	}

	sqlStmt := `
	create table if not exists articles (id integer not null primary key autoincrement, title text, content text);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return nil, err
	}

	return &DB{db: db}, nil
}

func (d *DB) Close() error {
	return d.db.Close()
}

func (d *DB) CreateArticle(article *Article) error {
	query, err := d.db.Prepare("insert into articles(title,content) values (?,?)")
	if err != nil {
		return err
	}
	defer query.Close()

	_, err = query.Exec(article.Title, article.Content)
	if err != nil {
		return err
	}

	return nil
}

func (d *DB) GetAllArticles() ([]*Article, error) {
	query, err := d.db.Prepare("select id, title, content from articles")
	if err != nil {
		return nil, err
	}
	defer query.Close()

	result, err := query.Query()
	if err != nil {
		return nil, err
	}
	defer result.Close()

	articles := make([]*Article, 0)
	for result.Next() {
		data := new(Article)
		err := result.Scan(
			&data.ID,
			&data.Title,
			&data.Content,
		)
		if err != nil {
			return nil, err
		}
		articles = append(articles, data)
	}

	return articles, nil
}

func (d *DB) GetArticle(articleID string) (*Article, error) {
	query, err := d.db.Prepare("select id, title, content from articles where id = ?")
	if err != nil {
		return nil, err
	}
	defer query.Close()

	result := query.QueryRow(articleID)
	data := new(Article)
	err = result.Scan(&data.ID, &data.Title, &data.Content)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (d *DB) UpdateArticle(id string, article *Article) error {
	query, err := d.db.Prepare("update articles set title = ?, content = ? where id = ?")
	if err != nil {
		return err
	}
	defer query.Close()

	_, err = query.Exec(article.Title, article.Content, id)
	if err != nil {
		return err
	}

	return nil
}

func (d *DB) DeleteArticle(id string) error {
	query, err := d.db.Prepare("delete from articles where id = ?")
	if err != nil {
		return err
	}
	defer query.Close()

	_, err = query.Exec(id)
	if err != nil {
		return err
	}

	return nil
}
