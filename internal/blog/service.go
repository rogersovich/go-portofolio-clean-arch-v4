package blog

import (
	"context"
	"fmt"
	"time"

	"github.com/rogersovich/go-portofolio-clean-arch-v4/internal/author"
	"github.com/rogersovich/go-portofolio-clean-arch-v4/internal/blog_topic"
	"github.com/rogersovich/go-portofolio-clean-arch-v4/internal/reading_time"
	"github.com/rogersovich/go-portofolio-clean-arch-v4/internal/statistic"
	"github.com/rogersovich/go-portofolio-clean-arch-v4/internal/topic"
	"github.com/rogersovich/go-portofolio-clean-arch-v4/pkg/utils"
	"gorm.io/gorm"
)

type Service interface {
	GetAllBlogs() ([]BlogResponse, error)
	GetBlogByIdWithRelations(id int) (BlogRelationResponse, error)
	CreateBlog(p CreateBlogRequest) (BlogResponse, error)
}

type service struct {
	authorService      author.Service
	topicService       topic.Service
	statisticService   statistic.Service
	readingTimeService reading_time.Service
	blogTopicService   blog_topic.Service
	blogRepo           Repository
	db                 *gorm.DB
}

func NewService(
	authorSvc author.Service,
	topicSvc topic.Service,
	statisticSvc statistic.Service,
	readingTimeSvc reading_time.Service,
	blogTopicSvc blog_topic.Service,
	r Repository,
	db *gorm.DB) Service {
	return &service{
		authorService:      authorSvc,
		topicService:       topicSvc,
		statisticService:   statisticSvc,
		readingTimeService: readingTimeSvc,
		blogTopicService:   blogTopicSvc,
		blogRepo:           r,
		db:                 db,
	}
}

func (s *service) GetAllBlogs() ([]BlogResponse, error) {
	datas, err := s.blogRepo.FindAll()
	if err != nil {
		return nil, err
	}

	var result []BlogResponse
	for _, p := range datas {
		result = append(result, ToBlogResponse(p))
	}
	return result, nil
}

func (s *service) GetBlogByIdWithRelations(id int) (BlogRelationResponse, error) {
	data, err := s.blogRepo.FindByIdWithRelations(id)

	if err != nil {
		return BlogRelationResponse{}, err
	}

	// Mapping result
	blogMap := map[int]*BlogRelationResponse{}

	for _, row := range data {
		blogID := int(row.ID)

		//? "Comma-ok" itu fitur spesial
		_, exists := blogMap[blogID]
		if !exists {
			var publishedAtPointer *string
			if row.PublishedAt != nil {
				formattedPublishedAt := row.PublishedAt.Format("2006-01-02 15:04:05")
				publishedAtPointer = &formattedPublishedAt
			}

			blogAuthor := BlogAuthorDTO{
				AuthorID:   row.AuthorID,
				AuthorName: row.AuthorName,
			}

			blogReadingTime := BlogReadingTimeDTO{
				ReadingTimeID:               row.ReadingTimeID,
				ReadingTimeMinutes:          row.ReadingTimeMinutes,
				ReadingTimeTextLength:       row.ReadingTimeTextLength,
				ReadingTimeEstimatedSeconds: row.ReadingTimeEstimatedSeconds,
				ReadingTimeWordCount:        row.ReadingTimeWordCount,
				ReadingTimeType:             row.ReadingTimeType,
			}

			blogStatistic := BlogStatisticDTO{
				StatisticID:    row.StatisticID,
				StatisticLikes: row.StatisticLikes,
				StatisticViews: row.StatisticViews,
				StatisticType:  row.StatisticType,
			}

			blogMap[blogID] = &BlogRelationResponse{
				ID:              blogID,
				Title:           row.Title,
				DescriptionHTML: row.DescriptionHTML,
				BannerUrl:       row.BannerUrl,
				BannerFileName:  row.BannerFileName,
				Summary:         row.Summary,
				Status:          row.Status,
				PublishedAt:     publishedAtPointer,
				CreatedAt:       row.CreatedAt.Format("2006-01-02 15:04:05"),
				Author:          blogAuthor,
				ReadingTime:     blogReadingTime,
				Statistic:       blogStatistic,
			}
		}

		blogMap[blogID].Topics = append(blogMap[blogID].Topics, BlogTopicDTO{
			TopicID:   row.TopicID,
			TopicName: row.TopicName,
		})
	}

	// Convert Map to Struct
	var result BlogRelationResponse
	for _, v := range blogMap {
		result = *v
		break
	}

	return result, nil
}

func (s *service) CreateBlog(p CreateBlogRequest) (BlogResponse, error) {
	//todo: Check Author Id
	//ex: author_id = 1
	_, err := s.authorService.GetAuthorById(p.AuthorID)

	if err != nil {
		err = fmt.Errorf("author_id %d not found", p.AuthorID)
		return BlogResponse{}, err
	}

	//todo: Check Topic Ids
	//ex: topic_ids = [1, 2, 3]
	topic_ids, _ := utils.ConvertStringSliceToIntSlice(p.TopicIds)
	_, err = s.topicService.CheckTopicIds(topic_ids)
	if err != nil {
		return BlogResponse{}, err
	}

	var publishedAt *time.Time
	var status string
	if p.IsPublished == "Y" {
		now := time.Now()
		publishedAt = &now
		status = "PUBLISHED"
	} else if p.IsPublished == "N" {
		status = "UNPUBLISHED"
	}

	tx := s.db.Begin()

	//todo: Create Statistic
	zero := 0
	pStatistic := statistic.CreateStatisticRequest{
		Likes: &zero,
		Views: &zero,
		Type:  "Blog",
	}
	dataStatistic, err := s.statisticService.CreateStatisticWithTx(pStatistic, tx)

	if err != nil {
		tx.Rollback()
		return BlogResponse{}, err
	}

	//todo: Create Reading Time
	readingTimeStats := utils.ExtractHTMLtoStatistics(p.DescriptionHTML)
	pReadingTime := reading_time.CreateReadingTimeRequest{
		Minutes:          readingTimeStats.Minutes,
		TextLength:       readingTimeStats.TextLength,
		EstimatedSeconds: readingTimeStats.EstimatedSeconds,
		WordCount:        readingTimeStats.WordCount,
		Type:             "Blog",
	}

	dataReadingTime, err := s.readingTimeService.CreateReadingTime(pReadingTime, tx)

	if err != nil {
		tx.Rollback()
		return BlogResponse{}, err
	}

	//todo: Upload Banner
	bannerRes, err := utils.HandlUploadFile(p.BannerFile, "blog")
	if err != nil {
		tx.Rollback()
		return BlogResponse{}, err
	}

	uploadedImageFilName := bannerRes.FileName

	//todo: Create Blog
	payload := CreateBlogDTO{
		AuthorID:        p.AuthorID,
		StatisticID:     dataStatistic.ID,
		ReadingTimeID:   dataReadingTime.ID,
		TopicIds:        p.TopicIds,
		Title:           p.Title,
		DescriptionHTML: p.DescriptionHTML,
		BannerUrl:       bannerRes.FileURL,
		BannerFileName:  bannerRes.FileName,
		Summary:         p.Summary,
		Status:          status,
		PublishedAt:     publishedAt,
	}

	data, err := s.blogRepo.CreateBlog(payload, tx)
	if err != nil {
		tx.Rollback()
		//? Delete banner image
		if uploadedImageFilName != "" {
			_ = utils.DeleteFromMinio(context.Background(), uploadedImageFilName)
		}
		return BlogResponse{}, err
	}

	//todo: Create Blog Topic
	err = s.blogTopicService.BulkCreateBlogTopic(topic_ids, data.ID, tx)
	if err != nil {
		tx.Rollback()
		//? Delete banner image
		if uploadedImageFilName != "" {
			_ = utils.DeleteFromMinio(context.Background(), uploadedImageFilName)
		}
		return BlogResponse{}, err
	}

	if err := tx.Commit().Error; err != nil {
		err = fmt.Errorf("error commit transaction")
		return BlogResponse{}, err
	}

	// return BlogResponse{}, nil
	return ToBlogResponse(data), nil
}
