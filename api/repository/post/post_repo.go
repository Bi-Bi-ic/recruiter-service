package post

import (
	"errors"
	"net/http"
	"regexp"
	"strings"

	"github.com/lib/pq"
	"github.com/mitchellh/mapstructure"
	"github.com/rgrs-x/service/api/factory"
	"github.com/rgrs-x/service/api/models"
	repo "github.com/rgrs-x/service/api/repository"
	"github.com/rgrs-x/service/api/repository/company"
	u "github.com/rgrs-x/service/api/utils"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type postStorage struct {
	DB *gorm.DB
}

//NewPostRepository ... We can implement to use PostRepository interface{} there
func NewPostRepository(db *gorm.DB) repo.PostRepository {
	return &postStorage{
		DB: db,
	}
}

// valid if Introduction is exist
func (storage *postStorage) checkIntroductionExist(creatorID string) repo.RepoResponse {
	var result bool
	//@ check if Introduction is existed
	commonDB, _ := storage.DB.DB()
	row := commonDB.QueryRow("SELECT EXISTS(SELECT introduction FROM posts WHERE creator_id = $1 AND type = 'introduction')", creatorID)
	row.Scan(&result)
	if result {
		logrus.WithFields(logrus.Fields{
			"creator": creatorID,
		}).Info("Introduction existed !!!")

		return repo.RepoResponse{Status: false}
	}

	return repo.RepoResponse{Status: true}
}

// Validate incoming details ...
func (storage *postStorage) Validate(post models.Post) bool {
	//Valid for Main details
	if result := storage.ValidateBlank(post.Title, post.Description); !result {
		return false
	}

	return true
}

//Check if field is blank
func (storage *postStorage) ValidateBlank(details ...string) bool {
	for _, detail := range details {
		re := regexp.MustCompile(`(?m)^\s*$`)
		result := re.MatchString(detail)
		if result == true {
			return false
		}
	}
	return true
}

//Create make a valid Post
func (storage *postStorage) Create(content models.Post, creatorID string) (post models.Post, status repo.Status) {

	if ok := storage.Validate(content); !ok {
		status = repo.CanNotCreate
		return
	}

	content.CreatorID = creatorID
	var factory factory.FileInfoFactoty

	file := factory.Create(content.Cover+".jpeg", content.Cover)
	content.Cover = file.Link

	content.TotalLike = 0

	// Set Type Post
	content.Type = "normal"

	queryStmt := storage.DB.Create(&content)
	err := queryStmt.Error
	if err != nil {
		status = repo.CanNotCreate
		return
	}

	post = content
	status = repo.Created
	return
}

// CreateIntroduction ...
func (storage *postStorage) CreateIntroduction(content models.Post, creatorID string) (post models.Post, status repo.Status) {
	if ok := storage.Validate(content); !ok {
		status = repo.CanNotCreate
		return
	}

	content.CreatorID = creatorID
	var factory factory.FileInfoFactoty

	file := factory.Create(content.Cover+".jpeg", content.Cover)
	content.Cover = file.Link

	content.TotalLike = 0

	// Set Type Post
	content.Type = "introduction"

	queryStmt := storage.DB.Create(&content)
	err := queryStmt.Error
	if err != nil {
		status = repo.CanNotCreate
		return
	}

	post = content
	status = repo.Created
	return
}

// //GetPost return post by Post-ID
func (storage *postStorage) GetPost(id uuid.UUID) (u.ResultRepository, int) {
	var post models.Post
	post.ID = id

	queryStatement := storage.DB.Table("posts").First(&post)

	err := queryStatement.Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.ResultRepository{Result: []string{}, Error: ErrContentsNotFound}, http.StatusNotFound
		}
		return u.ResultRepository{Result: []string{}, Error: repo.ErrRequestTooLong}, http.StatusRequestTimeout
	}

	message, statusCode := storage.FetchCreator(post.CreatorID, &post)
	if statusCode != http.StatusOK {
		return message, statusCode
	}

	return u.ResultRepository{Result: post, Message: "Found Post"}, http.StatusOK
}

// GetIntroduction ...
func (storage *postStorage) GetIntroduction(creatorID string) (result repo.RepoResponse, status repo.Status) {
	var post models.Post
	queryStmt := storage.DB.Table("posts").Where("creator_id = ? AND type = 'introduction'", creatorID).Find(&post)
	err := queryStmt.Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result = repo.RepoResponse{Status: false}
			status = repo.NotFound

			return
		}

		result = repo.RepoResponse{Status: false}
		status = repo.GetError

		return
	}

	_, statusCode := storage.FetchCreator(creatorID, &post)
	if statusCode != http.StatusOK {
		result = repo.RepoResponse{Status: false}
		status = repo.GetError

		return
	}

	result = repo.RepoResponse{Status: true, Data: post}
	status = repo.Success

	return
}

// GetPostDetails fetch post by ID given
func (storage *postStorage) GetPostDetails(contentID string) (post models.Post, status repo.Status) {
	queryStmt := storage.DB.Where("id = ? AND type = 'normal'", contentID)
	find := queryStmt.First(&post)
	err := find.Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			status = repo.NotFound
			return
		}
		status = repo.CannotGet
		return
	}

	status = repo.Success
	return
}

// GetIntroductionDetails fetch introduction by ID given
func (storage *postStorage) GetIntroductionDetails(contentID string) (post models.Post, status repo.Status) {
	queryStmt := storage.DB.Where("id = ? AND type = 'introduction'", contentID)
	find := queryStmt.First(&post)
	err := find.Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			status = repo.NotFound
			return
		}
		status = repo.CannotGet
		return
	}

	status = repo.Success
	return
}

// //GetAllPosts Return an array of posts
func (storage *postStorage) GetAllPosts() []models.Post {

	var contents []models.Post
	if err := storage.DB.Where("delete_at is NULL").Find(&contents).Error; err != nil {
		panic(err)
	}

	for key, value := range contents {
		_, _ = storage.FetchCreator(value.CreatorID, &contents[key])
	}

	return contents
}

// GetPartnerContents ...
func (storage *postStorage) GetPartnerContents(id uuid.UUID) (u.ResultRepository, int) {
	var contents []models.Post

	queryStatement := storage.DB.Where("creator_id = ? AND delete_At is NULL", id).Find(&contents)
	err := queryStatement.Error
	if err != nil {
		return u.ResultRepository{Result: []string{}, Error: repo.ErrRequestTooLong}, http.StatusRequestTimeout
	}

	if len(contents) <= 0 {
		return u.ResultRepository{Result: []string{}, Error: ErrContentsNotFound}, http.StatusNotFound
	}

	for key, value := range contents {
		_, _ = storage.FetchCreator(value.CreatorID, &contents[key])
	}

	return u.ResultRepository{Result: contents, Message: "Found all Contents"}, models.Contents
}

//GetCompanyContents fetch all Company's Contents
func (storage *postStorage) GetCompanyContents(id uuid.UUID) (u.ResultRepository, int) {
	var result bool

	commonDB, _ := storage.DB.DB()
	commonDB.QueryRow("SELECT EXISTS(SELECT id FROM companies WHERE id = $1)", id).Scan(&result)
	if !result {
		return u.ResultRepository{Result: []string{}, Error: company.ErrCompanyNotFound}, http.StatusForbidden
	}

	var contents []models.Post
	queryStatement := storage.DB.Table("posts").Where("company_id = ?", id).Find(&contents)

	err := queryStatement.Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.ResultRepository{Result: []string{}, Error: ErrContentsNotFound}, http.StatusNotFound
		}
		return u.ResultRepository{Result: []string{}, Error: repo.ErrRequestTooLong}, http.StatusRequestTimeout
	}

	for key, value := range contents {
		_, _ = storage.FetchCreator(value.CreatorID, &contents[key])
	}
	return u.ResultRepository{Result: contents, Message: "Found all Contents"}, models.Contents
}

//UpdatePost update an exist post
func (storage *postStorage) UpdatePost(post models.Post, postID string, creator string) (map[string]interface{}, int) {

	/*First init a new pointer to a null Post struct
	and fetch data Post from database
	*/
	content := &models.Post{}
	err := storage.DB.Table("posts").Where("id = ? AND delete_at is NULL", postID).First(content).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "Post Not Found to Update"), http.StatusNotFound
		}
		return u.Message(false, "Something went wrong . Please retry"), http.StatusRequestTimeout
	}

	//Well we will valid ID of the Creator of Post
	if content.CreatorID != creator {
		return u.Message(false, "Sorry! You are not the Author of this Post"), http.StatusForbidden
	}

	//Now fetch the Creator detail to keep Creator's informations are always updated
	message, statusCode := storage.FetchCreator(creator, &post)
	if statusCode != http.StatusOK {
		return u.Message(false, message.Error.Error()), statusCode
	}

	//OK... Now check if details update is valid
	if ok := storage.Validate(post); !ok {
		return u.Message(false, "Invalid Post request"), http.StatusBadRequest
	}

	/*Then assign Post ID */
	post.ID = content.ID

	//Update Post with Details receive from request
	//Time
	post.CreateAt = content.CreateAt

	storage.DB.Model(&post).Omit("id", "create_at", "creator_id", "company_id", "delete_at", "total_like").Updates(post)

	response := u.Message(true, "Post Updated")
	response["data"] = post

	return response, http.StatusOK
}

//DeletePost delete an existing Post
func (storage *postStorage) DeletePost(id string, creator string) (map[string]interface{}, int) {
	var post models.Post

	err := storage.DB.Table("posts").Where("id = ? AND delete_at is NULL", id).First(&post).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "Post Not Found"), http.StatusNotFound
		}
		return u.Message(false, "Something went wrong . Please retry"), http.StatusRequestTimeout
	}

	if post.CreatorID != creator {
		return u.Message(false, "Sorry! You are not the Author of this Post"), http.StatusForbidden
	}

	err = storage.DB.Delete(&post).Error

	if err != nil {
		panic(err)
	}

	response := u.Message(true, "Your post Successfully Deleted")
	return response, http.StatusOK
}

//Get AllTags show all Tags of Posts
func (storage *postStorage) GetAllTags() (map[string]interface{}, int) {

	var contents []models.Post

	err := storage.DB.Select("DISTINCT(tags)").Where("delete_at is NULL").Find(&contents).Error
	if err != nil {
		return u.Message(false, "Something went wrong. Please try again"), http.StatusRequestTimeout
	}

	if len(contents) <= 0 {
		return u.Message(false, "Can not find any Tags"), http.StatusNotFound
	}

	var result []pq.StringArray

	for _, value := range contents {
		result = append(result, value.Tags)
	}

	response := u.Message(true, "Got All Tags")
	response["data"] = result

	return response, http.StatusOK
}

//FetchCreator ...
func (storage *postStorage) FetchCreator(creatorID string, post *models.Post) (u.ResultRepository, int) {
	var partner models.Partner
	partner.ID, _ = uuid.FromString(creatorID)

	err := storage.DB.Where("id = ?", partner.ID).First(&partner).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.ResultRepository{Result: []string{}, Error: errors.New("Creator is not valid")}, http.StatusForbidden
		}
		return u.ResultRepository{Result: []string{}, Error: errors.New("Something went wrong . Please retry")}, http.StatusRequestTimeout
	}

	if post.CompanyID == "" && post.CompanyName == "" {
		post.CompanyID = partner.CompanyID
		post.CompanyName = partner.Name
	}

	post.UserName = partner.UserName
	post.PartnerName = partner.PartnerName

	post.Address = partner.Address
	post.CreatorAvatar = partner.Avatar

	post.MailContact = partner.MailContact
	post.Link = partner.Link
	post.Phone = partner.Phone
	return u.ResultRepository{Result: []string{}, Message: "Everything is OK !"}, http.StatusOK
}

//Pagination ... in development for page number post
func (storage *postStorage) CountContents(pagination *models.Pagination) error {
	var totalRows int64
	errCount := storage.DB.Model(&models.Post{}).Where("delete_at is NULL").Count(&totalRows).Error
	if errCount != nil {
		return errCount
	}

	pagination.TotalContents = int(totalRows)
	return nil
}

func (storage *postStorage) Pagination(pagination *models.Pagination) u.ResultRepository {
	var posts []models.Post

	offset := pagination.Offset

	// check SortType
	sortType := pagination.Sort.AsString()
	if sortType == "" {
		return u.ResultRepository{Result: []string{}, Error: errors.New("Unknown Sort Type")}
	}
	// get data with limit, offset & order
	find := storage.DB.Where("delete_at is NULL").Limit(pagination.Limit).Offset(offset).Order(pagination.Sort.AsString())

	//execute query
	find = find.Find(&posts)

	// has error find data
	errFind := find.Error
	if errFind != nil {
		return u.ResultRepository{Error: errFind}
	}

	for key, value := range posts {
		_, _ = storage.FetchCreator(value.CreatorID, &posts[key])
	}
	pagination.Rows = posts
	return u.ResultRepository{Result: pagination}
}

func (storage *postStorage) Filter(filter *models.Filter) u.ResultRepository {
	var posts []models.Post

	query := storage.DB.Where("delete_at IS NULL")

	find := query.Find(&posts)

	errFind := find.Error

	if errFind != nil {
		return u.ResultRepository{Error: errFind}
	}

	for key, value := range posts {
		_, _ = storage.FetchCreator(value.CreatorID, &posts[key])
	}

	var contents []models.Post
	if len(filter.JobKind) != 0 {
		for _, post := range posts {
			for _, value := range post.JobKind {
				if storage.CheckFilter(value, filter.JobKind) {
					contents = append(contents, post)
				}
			}
		}
	}

	if len(filter.Position) != 0 {
		for _, post := range posts {
			if storage.CheckFilter(post.Position, filter.Position) {
				contents = append(contents, post)
			}
		}
	}

	if len(filter.District) != 0 {
		for _, post := range posts {
			if storage.CheckFilter(post.District, filter.District) {
				contents = append(contents, post)
			}
		}
	}

	for _, post := range contents {
		contents = storage.RemoveSamePost(post.ID.String(), contents)
	}

	filter.Rows = contents
	filter.TotalContents = len(contents)
	return u.ResultRepository{Result: filter, Message: "Found all Result !!!"}

}

func (storage *postStorage) CheckFilter(input string, test []string) bool {
	for _, value := range test {
		if ok := strings.Compare(value, input); ok == 0 {
			return true
		}
	}
	return false
}

func (storage *postStorage) RemoveSamePost(input string, test []models.Post) []models.Post {
	var count int
	var result []models.Post
	for _, post := range test {
		if input == post.ID.String() {
			count++
		}
		if count > 1 {
			count--
			continue
		}
		result = append(result, post)
	}
	return result
}

/* Like Post---------------------------------------------------*/
func (storage *postStorage) UpdatePostLike(id uuid.UUID) (u.ResultRepository, int) {
	message, statusCode := storage.GetPost(id)
	if statusCode != http.StatusOK {
		return message, statusCode
	}

	var content models.Post
	err := mapstructure.Decode(message.Result, &content)
	if err != nil {
		panic(err)
	}

	content.TotalLike++

	queryStatement := storage.DB.Model(&content).Updates(map[string]interface{}{"total_like": content.TotalLike})

	err = queryStatement.Error
	if err != nil {
		return u.ResultRepository{Result: []string{}, Error: repo.ErrRequestTooLong}, http.StatusRequestTimeout
	}

	return u.ResultRepository{Result: content, Message: "Liked Post"}, http.StatusOK
}

/* Tracking Post---------------------------------------------------*/
func (storage *postStorage) UpdatePostReview(id uuid.UUID) (u.ResultRepository, int) {
	message, statusCode := storage.GetPost(id)
	if statusCode != http.StatusOK {
		return message, statusCode
	}

	var content models.Post
	err := mapstructure.Decode(message.Result, &content)
	if err != nil {
		panic(err)
	}

	content.TotalView++

	queryStatement := storage.DB.Model(&content).Updates(map[string]interface{}{"total_view": content.TotalView})

	err = queryStatement.Error
	if err != nil {
		return u.ResultRepository{Result: []string{}, Error: repo.ErrRequestTooLong}, http.StatusRequestTimeout
	}

	return u.ResultRepository{Result: content, Message: "Viewed Post"}, http.StatusOK
}
