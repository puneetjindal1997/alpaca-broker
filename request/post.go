package request

import (
	apperr "github.com/alpacahq/ribbit-backend/apperr"
	"github.com/gin-gonic/gin"
)

type UserPost struct {
	ID           int           `json:"-"`
	Description  string        `json:"description,omitempty" binding:"omitempty,min=2"`
	Title        string        `json:"title,omitempty" binding:"omitempty,min=2"`
	AddedBy      string        `json:"added_by,omitempty"`
	ParentID     int           `json:"parent_id,omitempty"`
	Type         string        `json:"type,omitempty"`
	Email        string        `json:"email,omitempty"`
	IsPrivate    bool          `json:"is_private,omitempty"`
	CommentCount int           `json:"comment_count,omitempty"`
	Attachements []interface{} `json:"attachements,omitempty"`
}

// AccountCreate validates account creation request
func PostCreate(c *gin.Context) (*UserPost, error) {
	var r UserPost
	if err := c.ShouldBindJSON(&r); err != nil {
		apperr.Response(c, err)
		return nil, err
	}
	return &r, nil
}
