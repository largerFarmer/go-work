package gorm

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

/**

	题目1：模型定义
	假设你要开发一个博客系统，有以下几个实体： User （用户）、 Post （文章）、 Comment （评论）。
	要求 ：
	使用Gorm定义 User 、 Post 和 Comment 模型，其中 User 与 Post 是一对多关系（一个用户可以发布多篇文章）， Post 与 Comment 也是一对多关系
	（一篇文章可以有多个评论）。
	编写Go代码，使用Gorm创建这些模型对应的数据库表。
	题目2：关联查询
	基于上述博客系统的模型定义。
	要求 ：
	编写Go代码，使用Gorm查询某个用户发布的所有文章及其对应的评论信息。
	编写Go代码，使用Gorm查询评论数量最多的文章信息。
	题目3：钩子函数
	继续使用博客系统的模型。
	要求 ：
	为 Post 模型添加一个钩子函数，在文章创建时自动更新用户的文章数量统计字段。
	为 Comment 模型添加一个钩子函数，在评论删除时检查文章的评论数量，如果评论数量为 0，则更新文章的评论状态为 "无评论"。

**/

// User 用户模型
type User struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Username  string    `gorm:"size:64;not null;unique" json:"username"`
	PostCount int       `gorm:"default:0" json:"post_count"` // 文章数量
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// 关联
	Posts []Post `gorm:"foreignKey:UserID" json:"posts,omitempty"`
}

// Post 文章模型
type Post struct {
	ID            uint      `gorm:"primarykey" json:"id"`
	Title         string    `gorm:"size:200;not null" json:"title"`
	Content       string    `gorm:"type:text;not null" json:"content"`
	CommentStatus string    `gorm:"size:20;default:'无评论'" json:"comment_status"` // 有评论, 无评论
	UserID        uint      `gorm:"not null" json:"user_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	// 关联
	User     User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Comments []Comment `gorm:"foreignKey:PostID" json:"comments,omitempty"`
}

// Comment 评论模型
type Comment struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	PostID    uint      `gorm:"not null" json:"post_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// 关联
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Post Post `gorm:"foreignKey:PostID" json:"post,omitempty"`
}

// 钩子函数

// Post创建后，更新用户的文章数量
func (p *Post) AfterCreate(tx *gorm.DB) error {
	// 更新用户的文章数量
	err := tx.Model(&User{}).Where("id = ?", p.UserID).UpdateColumn("post_count", gorm.Expr("post_count + 1")).Error
	if err != nil {
		return fmt.Errorf("更新用户文章数量失败: %v", err)
	}
	return nil
}

// Comment删除后的钩子函数 - 检查并更新文章的评论状态
func (c *Comment) AfterDelete(tx *gorm.DB) error {
	// 评论数量>0，更新为“有状态”
	var commentCount int64
	err := tx.Model(&Comment{}).Where("post_id = ?", c.PostID).Count(&commentCount).Error
	if err != nil {
		return fmt.Errorf("查询文章评论数量失败: %v", err)
	}
	if commentCount > 0 {
		err = tx.Model(&Post{}).Where("id = ?", c.PostID).Update("comment_status", "有评论").Error
		if err != nil {
			return fmt.Errorf("更新文章评论状态失败: %v", err)
		}
	}
	return nil
}

func GormSql() {

	db, err := gorm.Open(mysql.Open("root:root@123@tcp(127.0.0.1:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local"))
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}
	// 创建表
	db.AutoMigrate(&User{}, &Post{}, &Comment{})

	// 添加测试数据

	// 用户
	users := []User{{Username: "name1"}, {Username: "name2"}, {Username: "name3"}}
	result := db.Create(&users)
	if result.Error != nil {
		log.Fatal("创建用户失败：", result.Error)
	}
	// 文章
	posts := []Post{
		{Title: "标题1", Content: "内容...", UserID: users[0].ID},
		{Title: "标题2", Content: "内容...", UserID: users[0].ID},
		{Title: "标题3", Content: "内容...", UserID: users[1].ID},
		{Title: "标题4", Content: "内容...", UserID: users[2].ID},
	}
	result = db.Create(&posts)
	if result.Error != nil {
		log.Fatal("创建文章失败：", result.Error)
	}
	// 评论
	comments := []Comment{
		{Content: "评论1", UserID: users[1].ID, PostID: posts[0].ID},
		{Content: "评论2", UserID: users[2].ID, PostID: posts[0].ID},
		{Content: "评论3", UserID: users[1].ID, PostID: posts[1].ID},
		{Content: "评论4", UserID: users[2].ID, PostID: posts[2].ID},
		{Content: "评论5", UserID: users[0].ID, PostID: posts[3].ID},
	}
	result = db.Create(&comments)
	if result.Error != nil {
		log.Fatal("创建评论失败：", result.Error)
	}

	// 编写Go代码，使用Gorm查询某个用户发布的所有文章及其对应的评论信息。
	var user User
	userName := "name1"
	err1 := db.Preload("Posts.Comments.User").
		Preload("Posts.User").
		Where("username = ?", userName).
		First(&user).Error
	if err1 != nil {
		log.Fatal("查询用户文章和评论失败: ", err1)
		return
	}
	fmt.Printf("用户: %s\n", user.Username)
	for _, post := range user.Posts {
		fmt.Printf("文章: %s\n", post.Title)
		for _, comment := range post.Comments {
			fmt.Printf("评论: %s\n", comment.Content)
		}
	}

	// 编写Go代码，使用Gorm查询评论数量最多的文章信息。
	var post Post
	err2 := db.Model(&Post{}).
		Select("posts.*, COUNT(comments.id) c ").
		Joins("LEFT JOIN comments ON comments.post_id = posts.id").
		Group("posts.id").
		Order("c DESC").
		First(&post).Error
	if err2 != nil {
		log.Fatal("查询评论最多的文章失败: ", err2)
		return
	}
	fmt.Println("")
	fmt.Println("评论最多的文章：", post.Title)

}
