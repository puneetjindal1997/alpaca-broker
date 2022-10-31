package service

import (
	"net/http"

	"github.com/alpacahq/ribbit-backend/apperr"
	"github.com/alpacahq/ribbit-backend/model"
	"github.com/alpacahq/ribbit-backend/repository/user"
	"github.com/alpacahq/ribbit-backend/request"

	"github.com/gin-gonic/gin"
)

// User represents the user http service
type User struct {
	svc *user.Service
}

// UserRouter declares the orutes for users router group
func UserRouter(svc *user.Service, r *gin.RouterGroup) {
	u := User{
		svc: svc,
	}
	ur := r.Group("/users")
	ur.GET("", u.list)
	ur.GET("/:id", u.view)
	ur.PATCH("/:id", u.update)
	ur.DELETE("/:id", u.delete)
	ur.GET("/search", u.search)
	// post
	ur.POST("/post", u.post)
	ur.GET("/post", u.listpost)
	ur.GET("/post/private", u.listprivatepost)
	ur.PATCH("/post/:id", u.updatepost)
	ur.DELETE("/post-delete/:id", u.deletepost)

	ur.POST("/comment", u.createcomment)
	ur.GET("/comment/:post_id", u.listpostcomments)
	ur.PATCH("/comment/:post_id", u.updatepostcomment)
	ur.DELETE("/comment/:post_id", u.deletepostcomment)
}

type listResponse struct {
	Users []model.User `json:"users"`
	Page  int          `json:"page"`
}

type listPostResponse struct {
	Post []model.Post `json:"users"`
	Page int          `json:"page"`
}

func (u *User) list(c *gin.Context) {
	p, err := request.Paginate(c)
	if err != nil {
		return
	}
	result, err := u.svc.List(c, &model.Pagination{
		Limit: p.Limit, Offset: p.Offset,
	})
	if err != nil {
		apperr.Response(c, err)
		return
	}
	c.JSON(http.StatusOK, listResponse{
		Users: result,
		Page:  p.Page,
	})
}

func (u *User) view(c *gin.Context) {
	id, err := request.ID(c)
	if err != nil {
		return
	}
	result, err := u.svc.View(c, id)
	if err != nil {
		apperr.Response(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (u *User) update(c *gin.Context) {
	updateUser, err := request.UserUpdate(c)
	if err != nil {
		return
	}
	userUpdate, err := u.svc.Update(c, &user.Update{
		ID:        updateUser.ID,
		FirstName: updateUser.FirstName,
		LastName:  updateUser.LastName,
		Mobile:    updateUser.Mobile,
		Phone:     updateUser.Phone,
		Address:   updateUser.Address,
	})
	if err != nil {
		apperr.Response(c, err)
		return
	}
	c.JSON(http.StatusOK, userUpdate)
}

func (u *User) delete(c *gin.Context) {
	id, err := request.ID(c)
	if err != nil {
		return
	}
	if err := u.svc.Delete(c, id); err != nil {
		apperr.Response(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

func (u *User) search(c *gin.Context) {
	s := c.DefaultQuery("s", "")
	user, err := u.svc.Search(c, s)
	if err != nil {
		return
	}
	c.JSON(http.StatusOK, user)
}

func (u *User) post(c *gin.Context) {
	r, err := request.PostCreate(c)
	if err != nil {
		return
	}
	post := &model.Post{
		Description:  r.Description,
		Title:        r.Title,
		Email:        r.Email,
		IsPrivate:    r.IsPrivate,
		Attachements: r.Attachements,
	}
	respPost, err := u.svc.CreatePost(c, post)
	if err != nil {
		apperr.Response(c, err)
		return
	}
	c.JSON(http.StatusOK, respPost)
}

func (u *User) listpost(c *gin.Context) {
	var finalResult []model.Post
	p, err := request.Paginate(c)
	if err != nil {
		return
	}
	result, err := u.svc.ListPost(c, &model.Pagination{
		Limit: p.Limit, Offset: p.Offset,
	})
	if err != nil {
		apperr.Response(c, err)
		return
	}
	var postIds []int
	for _, singlePost := range result {
		postIds = append(postIds, singlePost.ID)
	}
	// get comments for particular post
	comments, err := u.svc.GetCommentCount(postIds)
	for _, singlePost := range result {
		for _, comment := range comments {
			if singlePost.ID == comment.ParentID {
				singlePost.CommentCount = comment.Comments
			} else {
				singlePost.CommentCount = 0
			}
		}
		finalResult = append(finalResult, singlePost)
	}
	if err != nil {
		apperr.Response(c, err)
		return
	}
	c.JSON(http.StatusOK, listPostResponse{
		Post: finalResult,
		Page: p.Page,
	})
}

func (u *User) listprivatepost(c *gin.Context) {
	p, err := request.Paginate(c)
	if err != nil {
		return
	}
	result, err := u.svc.ListPrivatePost(c, &model.Pagination{
		Limit: p.Limit, Offset: p.Offset,
	})
	if err != nil {
		apperr.Response(c, err)
		return
	}
	c.JSON(http.StatusOK, listPostResponse{
		Post: result,
		Page: p.Page,
	})
}

func (u *User) updatepost(c *gin.Context) {
	updatePost, err := request.PostUpdate(c)
	if err != nil {
		return
	}
	userUpdate, err := u.svc.UpdatePost(c, &user.UpdatePost{
		ID:           updatePost.ID,
		Description:  updatePost.Description,
		Title:        updatePost.Title,
		IsPrivate:    updatePost.IsPrivate,
		Attachements: updatePost.Attachements,
	})
	if err != nil {
		apperr.Response(c, err)
		return
	}
	c.JSON(http.StatusOK, userUpdate)
}

func (u *User) deletepost(c *gin.Context) {
	id, err := request.ID(c)
	if err != nil {
		return
	}
	if err := u.svc.DeletePost(c, id); err != nil {
		apperr.Response(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

func (u *User) createcomment(c *gin.Context) {
	r, err := request.PostCreate(c)
	if err != nil {
		return
	}
	post := &model.Post{
		ParentID:     r.ParentID,
		Description:  r.Description,
		Title:        r.Title,
		Email:        r.Email,
		Attachements: r.Attachements,
	}
	respPost, err := u.svc.CreatePost(c, post)
	if err != nil {
		apperr.Response(c, err)
		return
	}
	c.JSON(http.StatusOK, respPost)
}

func (u *User) listpostcomments(c *gin.Context) {
	id, err := request.PostID(c)
	if err != nil {
		return
	}
	p, err := request.Paginate(c)
	if err != nil {
		return
	}
	result, err := u.svc.ListPostComments(c, &model.Pagination{
		Limit: p.Limit, Offset: p.Offset,
	}, id)
	if err != nil {
		apperr.Response(c, err)
		return
	}
	c.JSON(http.StatusOK, listPostResponse{
		Post: result,
		Page: p.Page,
	})
}

func (u *User) updatepostcomment(c *gin.Context) {
	updatePost, err := request.PostUpdate(c)
	if err != nil {
		return
	}
	userUpdate, err := u.svc.UpdatePostComment(c, &user.UpdatePost{
		ID:          updatePost.ID,
		Description: updatePost.Description,
	})
	if err != nil {
		apperr.Response(c, err)
		return
	}
	c.JSON(http.StatusOK, userUpdate)
}

func (u *User) deletepostcomment(c *gin.Context) {
	id, err := request.PostID(c)
	if err != nil {
		return
	}
	if err := u.svc.DeletePostComment(c, id); err != nil {
		apperr.Response(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
