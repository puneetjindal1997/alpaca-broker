package user

import (
	"net/http"

	"github.com/alpacahq/ribbit-backend/apperr"
	"github.com/alpacahq/ribbit-backend/model"
	"github.com/alpacahq/ribbit-backend/repository/platform/query"
	"github.com/alpacahq/ribbit-backend/repository/platform/structs"

	"github.com/gin-gonic/gin"
)

// NewUserService create a new user application service
func NewUserService(userRepo model.UserRepo, auth model.AuthService, rbac model.RBACService) *Service {
	return &Service{
		userRepo: userRepo,
		auth:     auth,
		rbac:     rbac,
	}
}

// Service represents the user application service
type Service struct {
	userRepo model.UserRepo
	auth     model.AuthService
	rbac     model.RBACService
}

// List returns list of users
func (s *Service) List(c *gin.Context, p *model.Pagination) ([]model.User, error) {
	u := s.auth.User(c)
	q, err := query.List(u)
	if err != nil {
		return nil, err
	}
	return s.userRepo.List(q, p)
}

// View returns single user
func (s *Service) View(c *gin.Context, id int) (*model.User, error) {
	if !s.rbac.EnforceUser(c, id) {
		return nil, apperr.New(http.StatusForbidden, "Forbidden")
	}
	return s.userRepo.View(id)
}

// Search username return single user
func (s *Service) Search(c *gin.Context, searchKey string) (*model.User, error) {
	u := s.auth.User(c)
	if !s.rbac.EnforceUser(c, u.ID) {
		return nil, apperr.New(http.StatusForbidden, "Forbidden")
	}
	user, err := s.userRepo.FindByUsername(searchKey)
	if err != nil {
		user, err = s.userRepo.FindByEmail(searchKey)
	}
	if user.ID <= 0 {
		return nil, apperr.New(http.StatusOK, "Not found")
	}
	return user, err
}

// Update contains user's information used for updating
type Update struct {
	ID        int
	FirstName *string
	LastName  *string
	Mobile    *string
	Phone     *string
	Address   *string
}

// Update updates user's contact information
func (s *Service) Update(c *gin.Context, update *Update) (*model.User, error) {
	if !s.rbac.EnforceUser(c, update.ID) {
		return nil, apperr.New(http.StatusForbidden, "Forbidden")
	}
	u, err := s.userRepo.View(update.ID)
	if err != nil {
		return nil, err
	}
	structs.Merge(u, update)
	return s.userRepo.Update(u)
}

// Delete deletes a user
func (s *Service) Delete(c *gin.Context, id int) error {
	u, err := s.userRepo.View(id)
	if err != nil {
		return err
	}
	if !s.rbac.IsLowerRole(c, u.Role.AccessLevel) {
		return apperr.New(http.StatusForbidden, "Forbidden")
	}
	u.Delete()
	return s.userRepo.Delete(u)
}

func (s *Service) CreatePost(c *gin.Context, post *model.Post) (*model.Post, error) {
	u := s.auth.User(c)
	if !s.rbac.EnforceUser(c, u.ID) {
		return nil, apperr.New(http.StatusForbidden, "Forbidden")
	}
	resp, err := s.userRepo.CreateComment(post)
	return resp, err
}

// List returns list of public post
func (s *Service) ListPost(c *gin.Context, p *model.Pagination) ([]model.Post, error) {
	return s.userRepo.ListPost(&model.ListQuery{}, p)
}

// List returns list of private post
func (s *Service) ListPrivatePost(c *gin.Context, p *model.Pagination) ([]model.Post, error) {
	u := s.auth.User(c)
	q, err := query.List(u)
	if err != nil {
		return nil, err
	}
	return s.userRepo.ListPrivatePost(q, p)
}

func (s *Service) GetCommentCount(ids []int) ([]model.CommentCount, error) {
	return s.userRepo.ListPostCommentsCount(ids)
}

type UpdatePost struct {
	ID           int
	Description  *string
	Title        *string
	IsPrivate    *bool
	Attachements *[]interface{}
}

// Update updates user's contact information
func (s *Service) UpdatePost(c *gin.Context, update *UpdatePost) (*model.Post, error) {

	p, err := s.userRepo.ViewPost(update.ID)
	if err != nil {
		return nil, err
	}

	structs.Merge(p, update)
	return s.userRepo.UpdatePost(p)
}

// Delete deletes a user posts
func (s *Service) DeletePost(c *gin.Context, id int) error {
	p, err := s.userRepo.ViewPost(id)
	if err != nil {
		return err
	}
	p.Delete()
	return s.userRepo.DeletePost(p)
}

// List returns list of public post
func (s *Service) ListPostComments(c *gin.Context, p *model.Pagination, postId int) ([]model.Post, error) {
	return s.userRepo.ListPostComments(&model.ListQuery{}, p, postId)
}

// Update updates user's contact information
func (s *Service) UpdatePostComment(c *gin.Context, update *UpdatePost) (*model.Post, error) {

	p, err := s.userRepo.ViewPostComment(update.ID)
	if err != nil {
		return nil, err
	}

	structs.Merge(p, update)
	return s.userRepo.UpdatePostComment(p)
}

// Delete deletes a user posts
func (s *Service) DeletePostComment(c *gin.Context, id int) error {
	p, err := s.userRepo.ViewPostComment(id)
	if err != nil {
		return err
	}
	p.Delete()
	return s.userRepo.DeletePostComment(p)
}
