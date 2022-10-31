package repository

import (
	"net/http"

	"github.com/go-pg/pg/v9/orm"
	"go.uber.org/zap"

	"github.com/alpacahq/ribbit-backend/apperr"
	"github.com/alpacahq/ribbit-backend/model"
)

const notDeleted = "deleted_at is null"
const isPrivate = "is_private is false"
const isPrivateTrue = "is_private is true"
const comment = "comment"
const post = "post"

// NewUserRepo returns a new UserRepo instance
func NewUserRepo(db orm.DB, log *zap.Logger) *UserRepo {
	return &UserRepo{db, log}
}

// UserRepo is the client for our user model
type UserRepo struct {
	db  orm.DB
	log *zap.Logger
}

// View returns single user by ID
func (u *UserRepo) View(id int) (*model.User, error) {
	var user = new(model.User)
	sql := `SELECT "user".*, "role"."id" AS "role__id", "role"."access_level" AS "role__access_level", "role"."name" AS "role__name" 
	FROM "users" AS "user" LEFT JOIN "roles" AS "role" ON "role"."id" = "user"."role_id" 
	WHERE ("user"."id" = ? and deleted_at is null)`
	_, err := u.db.QueryOne(user, sql, id)
	if err != nil {
		u.log.Warn("UserRepo Error", zap.Error(err))
		return nil, apperr.New(http.StatusNotFound, "400 not found")
	}
	return user, nil
}

// View returns single user by referral code
func (u *UserRepo) FindByReferralCode(referralCode string) (*model.ReferralCodeVerifyResponse, error) {
	var user = new(model.ReferralCodeVerifyResponse)
	sql := `SELECT "user"."first_name", "user"."last_name", "user"."referral_code", "user"."username"
	FROM "users" AS "user" 
	WHERE ("user"."referral_code" = ? and deleted_at is null)`
	_, err := u.db.QueryOne(user, sql, referralCode)
	if err != nil {
		u.log.Warn("UserRepo Error", zap.Error(err))
		return nil, apperr.New(http.StatusNotFound, "400 not found")
	}
	return user, nil
}

// FindByUsername queries for a single user by username
func (u *UserRepo) FindByUsername(username string) (*model.User, error) {
	user := new(model.User)
	sql := `SELECT "user".*, "role"."id" AS "role__id", "role"."access_level" AS "role__access_level", "role"."name" AS "role__name" 
	FROM "users" AS "user" LEFT JOIN "roles" AS "role" ON "role"."id" = "user"."role_id" 
	WHERE ("user"."username" = ? and deleted_at is null)`
	_, err := u.db.QueryOne(user, sql, username)
	if err != nil {
		u.log.Warn("UserRepo Error", zap.String("Error:", err.Error()))
		return nil, apperr.New(http.StatusNotFound, "400 not found")
	}
	return user, nil
}

// FindByEmail queries for a single user by email
func (u *UserRepo) FindByEmail(email string) (*model.User, error) {
	user := new(model.User)
	sql := `SELECT "user".*, "role"."id" AS "role__id", "role"."access_level" AS "role__access_level", "role"."name" AS "role__name" 
	FROM "users" AS "user" LEFT JOIN "roles" AS "role" ON "role"."id" = "user"."role_id" 
	WHERE ("user"."email" = ? and deleted_at is null)`
	_, err := u.db.QueryOne(user, sql, email)
	if err != nil {
		u.log.Warn("UserRepo Error", zap.String("Error:", err.Error()))
		return nil, apperr.New(http.StatusNotFound, "400 not found")
	}
	return user, nil
}

// FindByMobile queries for a single user by mobile (and country code)
func (u *UserRepo) FindByMobile(countryCode, mobile string) (*model.User, error) {
	user := new(model.User)
	sql := `SELECT "user".*, "role"."id" AS "role__id", "role"."access_level" AS "role__access_level", "role"."name" AS "role__name" 
	FROM "users" AS "user" LEFT JOIN "roles" AS "role" ON "role"."id" = "user"."role_id" 
	WHERE ("user"."country_code" = ? and "user"."mobile" = ? and deleted_at is null)`
	_, err := u.db.QueryOne(user, sql, countryCode, mobile)
	if err != nil {
		u.log.Warn("UserRepo Error", zap.String("Error:", err.Error()))
		return nil, apperr.New(http.StatusNotFound, "400 not found")
	}
	return user, nil
}

// FindByToken queries for single user by token
func (u *UserRepo) FindByToken(token string) (*model.User, error) {
	var user = new(model.User)
	sql := `SELECT "user".*, "role"."id" AS "role__id", "role"."access_level" AS "role__access_level", "role"."name" AS "role__name" 
	FROM "users" AS "user" LEFT JOIN "roles" AS "role" ON "role"."id" = "user"."role_id" 
	WHERE ("user"."token" = ? and deleted_at is null)`
	_, err := u.db.QueryOne(user, sql, token)
	if err != nil {
		u.log.Warn("UserRepo Error", zap.String("Error:", err.Error()))
		return nil, apperr.New(http.StatusNotFound, "400 not found")
	}
	return user, nil
}

// UpdateLogin updates last login and refresh token for user
func (u *UserRepo) UpdateLogin(user *model.User) error {
	user.UpdateLastLogin() // update user object's last_login field
	_, err := u.db.Model(user).Column("last_login", "token").WherePK().Update()
	if err != nil {
		u.log.Warn("UserRepo Error", zap.Error(err))
	}
	return err
}

// List returns list of all users retreivable for the current user, depending on role
func (u *UserRepo) List(qp *model.ListQuery, p *model.Pagination) ([]model.User, error) {
	var users []model.User
	q := u.db.Model(&users).Column("user.*", "Role").Where(notDeleted).Order("user.id desc").Limit(p.Limit).Offset(p.Offset)
	if qp != nil {
		q.Where(qp.Query, qp.ID)
	}
	if err := q.Select(); err != nil {
		u.log.Warn("UserDB Error", zap.Error(err))
		return nil, err
	}
	return users, nil
}

// Update updates user's contact info
func (u *UserRepo) Update(user *model.User) (*model.User, error) {
	_, err := u.db.Model(user).Column(
		"first_name",
		"last_name",
		"username",
		"mobile",
		"country_code",
		"address",
		"account_id",
		"account_number",
		"account_currency",
		"account_status",
		"dob",
		"city",
		"state",
		"country",
		"tax_id_type",
		"tax_id",
		"funding_source",
		"employment_status",
		"investing_experience",
		"public_shareholder",
		"another_brokerage",
		"device_id",
		"profile_completion",
		"bio",
		"facebook_url",
		"twitter_url",
		"instagram_url",
		"public_portfolio",
		"employer_name",
		"occupation",
		"unit_apt",
		"zip_code",
		"stock_symbol",
		"brokerage_firm_name",
		"brokerage_firm_employee_name",
		"brokerage_firm_employee_relationship",
		"shareholder_company_name",
		"avatar",
		"referred_by",
		"watchlist_id",
		"active",
		"verified",
		"updated_at",
	).WherePK().Update()
	if err != nil {
		u.log.Warn("UserDB Error", zap.Error(err))
	}
	return user, err
}

// Delete sets deleted_at for a user
func (u *UserRepo) Delete(user *model.User) error {
	user.Delete()
	_, err := u.db.Model(user).Column("deleted_at").WherePK().Update()
	if err != nil {
		u.log.Warn("UserRepo Error", zap.Error(err))
	}
	return err
}

// Create creates a new post in our database
func (u *UserRepo) CreatePost(post *model.Post) (*model.Post, error) {
	post.Type = post.PostType(1)
	user := new(model.Post)
	sql := `SELECT id FROM posts WHERE id = ? and type = ? and deleted_at is null`
	res, err := u.db.Query(user, sql, post.ID, post)
	if err != nil {
		u.log.Error("CreatePost Error: ", zap.Error(err))
		return nil, apperr.DB
	}
	if res.RowsReturned() != 0 {
		return nil, apperr.New(http.StatusBadRequest, "Post already exists.")
	}
	if err := u.db.Insert(post); err != nil {
		u.log.Warn("CreatePost error: ", zap.Error(err))
		return nil, apperr.DB
	}
	return post, nil
}

// List returns list of all post retreivable for the current user, depending on role
func (u *UserRepo) ListPost(qp *model.ListQuery, p *model.Pagination) ([]model.Post, error) {
	var posts []model.Post
	q := u.db.Model(&posts).Column("post.*").Where(isPrivate + " and " + notDeleted).Order("posts.id desc").Limit(p.Limit).Offset(p.Offset)
	if qp != nil {
		q.Where(qp.Query, qp.ID)
	}
	if err := q.Select(); err != nil {
		u.log.Warn("PostDB Error", zap.Error(err))
		return nil, err
	}
	return posts, nil
}

func (u *UserRepo) ListPostCommentsCount(postIds []int) ([]model.CommentCount, error) {
	var comments []model.CommentCount
	sql := `SELECT "post".* AS comments "post"."id" AS "post_id" , "post"."parent_id" AS "parent_id" 
	FROM "posts" AS "posts" WHERE (parent_id IN ? and deleted_at is null) GROUP BY parent_id`
	_, err := u.db.QueryOne(comments, sql, postIds)
	if err != nil {
		u.log.Warn("CommentRepo Error", zap.String("Error:", err.Error()))
		return nil, apperr.New(http.StatusNotFound, "400 not found")
	}
	return comments, nil
}

// List returns list of all post retreivable for the current user, depending on role
func (u *UserRepo) ListPrivatePost(qp *model.ListQuery, p *model.Pagination) ([]model.Post, error) {
	var posts []model.Post
	q := u.db.Model(&posts).Column("post.*").Where(isPrivateTrue + " and " + notDeleted).Order("posts.id desc").Limit(p.Limit).Offset(p.Offset)
	if qp != nil {
		q.Where(qp.Query, qp.ID)
	}
	if err := q.Select(); err != nil {
		u.log.Warn("PostDB Error", zap.Error(err))
		return nil, err
	}
	return posts, nil
}

// View returns single post by ID
func (u *UserRepo) ViewPost(id int) (*model.Post, error) {
	var post = new(model.Post)
	sql := `SELECT post.* FROM posts WHERE id = ? and deleted_at is null`
	_, err := u.db.QueryOne(post, sql, id)
	if err != nil {
		u.log.Warn("UserRepo Post Error", zap.Error(err))
		return nil, apperr.New(http.StatusNotFound, "400 not found")
	}
	return post, nil
}

func (u *UserRepo) UpdatePost(post *model.Post) (*model.Post, error) {
	_, err := u.db.Model(post).Column(
		"description",
		"title",
		"is_private",
		"attachements",
		"updated_at",
	).WherePK().Update()
	if err != nil {
		u.log.Warn("UserDB Post Error", zap.Error(err))
	}
	return post, err
}

// Delete sets deleted_at for a user
func (u *UserRepo) DeletePost(post *model.Post) error {
	post.Delete()
	_, err := u.db.Model(post).Column("deleted_at").WherePK().Update()
	if err != nil {
		u.log.Warn("UserRepo Error", zap.Error(err))
	}
	return err
}

// Create creates a new post in our database
func (u *UserRepo) CreateComment(post *model.Post) (*model.Post, error) {
	post.Type = post.PostType(0)
	user := new(model.Post)
	sql := `SELECT id FROM posts WHERE id = ? and parent_id = ? and deleted_at is null`
	res, err := u.db.Query(user, sql, post.ID)
	if err != nil {
		u.log.Error("UserComment Error: ", zap.Error(err))
		return nil, apperr.DB
	}
	if res.RowsReturned() != 0 {
		return nil, apperr.New(http.StatusBadRequest, "Comment already exists.")
	}
	if err := u.db.Insert(post); err != nil {
		u.log.Warn("UserComment error: ", zap.Error(err))
		return nil, apperr.DB
	}
	return post, nil
}

// List returns list of all post comment retreivable for the current post comments, depending on role
func (u *UserRepo) ListPostComments(qp *model.ListQuery, p *model.Pagination, postId int) ([]model.Post, error) {
	var posts []model.Post
	q := u.db.Model(&posts).Column("post.*").Where(isPrivate+" and "+notDeleted+" and parent_id = ? and type = ?", postId, comment).Order("posts.id desc").Limit(p.Limit).Offset(p.Offset)
	if qp != nil {
		q.Where(qp.Query, qp.ID)
	}
	if err := q.Select(); err != nil {
		u.log.Warn("PostDB Error", zap.Error(err))
		return nil, err
	}
	return posts, nil
}

func (u *UserRepo) ViewPostComment(id int) (*model.Post, error) {
	var post = new(model.Post)
	sql := `SELECT post.* FROM posts WHERE id = ? and type = ? and deleted_at is null`
	_, err := u.db.QueryOne(post, sql, id, comment)
	if err != nil {
		u.log.Warn("UserRepo Post Error", zap.Error(err))
		return nil, apperr.New(http.StatusNotFound, "400 not found")
	}
	return post, nil
}

func (u *UserRepo) UpdatePostComment(post *model.Post) (*model.Post, error) {
	_, err := u.db.Model(post).Column(
		"description",
		"updated_at",
	).WherePK().Update()
	if err != nil {
		u.log.Warn("CommentsDB Post Error", zap.Error(err))
	}
	return post, err
}

// Delete sets deleted_at for a user
func (u *UserRepo) DeletePostComment(post *model.Post) error {
	post.Delete()
	_, err := u.db.Model(post).Column("deleted_at").WherePK().Update()
	if err != nil {
		u.log.Warn("UserRepo Error", zap.Error(err))
	}
	return err
}
