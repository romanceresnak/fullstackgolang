package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Post struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Title     string    `gorm:"size:255;not null;unique" json:"title"`
	Content   string    `gorm:"size:255;not null;" json:"content"`
	Author    User      `json:"author"`
	AuthorID  uint32    `gorm:"not null" json:"author_id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (p *Post) Prepare() {
	p.ID = 0
	p.Title = html.EscapeString(strings.TrimSpace(p.Title))
	p.Content = html.EscapeString(strings.TrimSpace(p.Content))
	p.Author = User{}
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}

//A method is a function with a special receiver argument. (p *Post)
func (p *Post) Validate() error {
	if p.Title == "" {
		return errors.New("Required Title")
	}
	if p.Content == "" {
		return errors.New("Required Content")
	}
	if p.AuthorID < 1 {
		return errors.New("Required Author")
	}
	return nil
}

//Save post
func (p *Post) SavePost(db *gorm.DB) (*Post, error) {
	var err error
	//try to create Post object
	err = db.Debug().Model(&Post{}).Create(&p).Error

	//if it is error return empty Post object and error
	if err != nil {
		return &Post{}, err
	}

	if p.ID != 0 {
		//need to include Author by AuthorID
		err = db.Debug().Model(&User{}).Where("id = ?", p.AuthorID).Take(p.Author).Error
		if err != nil {
			return &Post{}, err
		}
	}

	//return saved object and nil
	return p, nil
}

func (p *Post) FindAllPosts(db *gorm.DB) (*[]Post, error) {
	var err error
	posts := []Post{}

	err = db.Debug().Model(&Post{}).Limit(100).Find(&posts).Error
	if err != nil {
		return &[]Post{}, err
	}

	if len(posts) > 0 {
		for i, _ := range posts {
			err := db.Debug().Model(&User{}).Where("id = ?", posts[i].AuthorID).Take(&posts[i].Author).Error
			if err != nil {
				return &[]Post{}, err
			}
		}

	}

	return &posts, nil
}

//Toto mapové [reťazcové] rozhranie {} bude obsahovať mapu reťazcov podľa ľubovoľných typov údajov.
func (p *Post) FindPostByID(db *gorm.DB, uid uint64) (*Post, error) {
	var err error
	db = db.Debug().Model(&Post{}).Where("id = ?", uid).Take(&Post{}).UpdateColumns(
		map[string]interface{}{
			"title":      p.ID,
			"content":    p.Content,
			"updated_at": time.Now(),
		},
	)
	err = db.Debug().Model(&Post{}).Where("id = ?", uid).Take(&p).Error
	if err != nil {
		return &Post{}, err
	}

	return p, nil
}

func (p *Post) UpdateAPost(db *gorm.DB, pid uint32) (*Post, error) {
	var err error

	db = db.Debug().Model(&Post{}).Where("id = ?", pid).Take(&Post{}).UpdateColumns(
		map[string]interface{}{
			"title":      p.Title,
			"content":    p.Content,
			"updated_at": time.Now(),
		},
	)
	err = db.Debug().Model(&Post{}).Where("id = ?", pid).Take(&p).Error
	if err != nil {
		return &Post{}, err
	}

	if p.ID != 0 {
		err = db.Debug().Model(&Post{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
		if err != nil {
			return &Post{}, err
		}
	}

	return p, nil
}

func (p *Post) DeleteAPost(db *gorm.DB, pid uint64, uid uint32) (int64, error) {

	db = db.Debug().Model(&Post{}).Where("id = ? and author_id = ?", pid, uid).Take(&Post{}).Delete(&Post{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Post not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
