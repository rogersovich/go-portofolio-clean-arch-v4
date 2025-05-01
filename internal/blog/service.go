package blog

import (
	"context"
	"fmt"
	"time"

	"github.com/rogersovich/go-portofolio-clean-arch-v4/internal/author"
	"github.com/rogersovich/go-portofolio-clean-arch-v4/internal/blog_content_image"
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
	GetBlogById(id int) (BlogResponse, error)
	CreateBlog(p CreateBlogRequest) (BlogResponse, error)
	UpdateBlog(p UpdateBlogRequest) (BlogUpdateResponse, error)
	DeleteBlog(id int) (Blog, error)
}

type service struct {
	authorService           author.Service
	topicService            topic.Service
	statisticService        statistic.Service
	readingTimeService      reading_time.Service
	blogTopicService        blog_topic.Service
	blogContentImageService blog_content_image.Service
	blogRepo                Repository
	db                      *gorm.DB
}

func NewService(
	authorSvc author.Service,
	topicSvc topic.Service,
	statisticSvc statistic.Service,
	readingTimeSvc reading_time.Service,
	blogTopicSvc blog_topic.Service,
	blogContentImageSvc blog_content_image.Service,
	r Repository,
	db *gorm.DB) Service {
	return &service{
		authorService:           authorSvc,
		topicService:            topicSvc,
		statisticService:        statisticSvc,
		readingTimeService:      readingTimeSvc,
		blogTopicService:        blogTopicSvc,
		blogContentImageService: blogContentImageSvc,
		blogRepo:                r,
		db:                      db,
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

func (s *service) GetBlogById(id int) (BlogResponse, error) {
	data, err := s.blogRepo.FindById(id)
	if err != nil {
		return BlogResponse{}, err
	}
	return data, nil
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
				ContentImages:   []BlogContentImageDTO{},
				Topics:          []BlogTopicDTO{},
			}
		}

		blogMap[blogID].Topics = append(blogMap[blogID].Topics, BlogTopicDTO{
			TopicID:   row.TopicID,
			TopicName: row.TopicName,
		})

		if row.BlogContentImageID != 0 {
			seen := make(map[int]bool) // Map to check if the ID is already seen
			for _, img := range blogMap[blogID].ContentImages {
				seen[img.BlogContentImageID] = true
			}

			if !seen[row.BlogContentImageID] {
				blogMap[blogID].ContentImages = append(blogMap[blogID].ContentImages, BlogContentImageDTO{
					BlogContentImageID:       row.BlogContentImageID,
					BlogContentImageUrl:      row.BlogContentImageUrl,
					BlogContentImageFileName: row.BlogContentImageFileName,
				})
			}
		}
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
	topic_ids := p.TopicIds
	_, err = s.topicService.CheckTopicIds(topic_ids)
	if err != nil {
		return BlogResponse{}, err
	}

	//todo: Check Content Images
	err = s.blogContentImageService.CountUnlinkedImages(p.ContentImages)
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

	//todo: Update Blog Content Images
	err = s.blogContentImageService.MarkImagesUsedByBlog(p.ContentImages, data.ID, tx)
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

	return ToBlogResponse(data), nil
}

func (s *service) UpdateBlog(p UpdateBlogRequest) (BlogUpdateResponse, error) {
	blog, err := s.GetBlogById(p.ID)
	if err != nil {
		return BlogUpdateResponse{}, err
	}

	//* set oldFileName
	oldFileName := ""
	if blog.BannerFileName != "" {
		oldFileName = blog.BannerFileName
	}

	//todo: Check Topic Ids
	//ex: topic_ids = [1, 2, 3]
	var topic_ids []int
	for _, item := range p.TopicIds {
		topic_ids = append(topic_ids, item.TopicID)
	}
	_, err = s.topicService.CheckTopicIds(topic_ids)
	if err != nil {
		return BlogUpdateResponse{}, err
	}

	//todo: Check Author Id
	_, err = s.authorService.GetAuthorById(p.AuthorID)
	if err != nil {
		return BlogUpdateResponse{}, err
	}

	//! todo: Begin Transaction

	tx := s.db.Begin()

	// todo: Update Blog Topics
	err = s.blogTopicService.BatchUpdateBlogTopic(topic_ids, p.ID, tx)
	if err != nil {
		tx.Rollback()
		return BlogUpdateResponse{}, err
	}

	// todo: Sync Blog Images
	oldImageBlogs, err := s.blogContentImageService.SyncBlogImages(p.ContentImages, p.ID, tx)

	if err != nil {
		tx.Rollback()
		return BlogUpdateResponse{}, err
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

	isUpdateReadingTime := false
	if p.DescriptionHTML != blog.DescriptionHTML {
		isUpdateReadingTime = true
	}

	if isUpdateReadingTime {
		//todo: Extract Reading Time
		readingTimeStats := utils.ExtractHTMLtoStatistics(p.DescriptionHTML)
		pReadingTime := reading_time.UpdateReadingTimeRequest{
			ID:               blog.ReadingTimeID,
			Minutes:          readingTimeStats.Minutes,
			TextLength:       readingTimeStats.TextLength,
			EstimatedSeconds: readingTimeStats.EstimatedSeconds,
			WordCount:        readingTimeStats.WordCount,
			Type:             "Blog",
		}

		//todo: Update Reading Time
		_, err := s.readingTimeService.UpdateReadingTime(pReadingTime, tx)

		if err != nil {
			tx.Rollback()
			return BlogUpdateResponse{}, err
		}
	}

	//todo: Handle Upload Banner
	var newFileURL string
	var newFileName string

	if p.BannerFile != nil {
		imageRes, err := utils.HandlUploadFile(p.BannerFile, "blog")
		if err != nil {
			return BlogUpdateResponse{}, err
		}

		newFileURL = imageRes.FileURL
		newFileName = imageRes.FileName
	} else {
		newFileURL = blog.BannerUrl // keep existing if not updated
		newFileName = blog.BannerFileName
	}

	payload := UpdateBlogDTO{
		ID:              p.ID,
		TopicIds:        p.TopicIds,
		AuthorID:        p.AuthorID,
		StatisticID:     blog.StatisticID,
		ReadingTimeID:   blog.ReadingTimeID,
		Title:           p.Title,
		DescriptionHTML: p.DescriptionHTML,
		BannerUrl:       newFileURL,
		BannerFileName:  newFileName,
		Summary:         p.Summary,
		Status:          status,
		PublishedAt:     publishedAt,
	}

	//todo: Update Blog
	_, err = s.blogRepo.UpdateBlog(payload, tx)
	if err != nil {
		tx.Rollback()
		if oldFileName != newFileName {
			_ = utils.DeleteFromMinio(context.Background(), oldFileName)
		}
		return BlogUpdateResponse{}, err
	}

	//todo: Delete Old Blog Images
	if len(oldImageBlogs) > 0 {
		slice_image_urls := []string{}
		for _, item := range oldImageBlogs {
			slice_image_urls = append(slice_image_urls, item.ImageUrl)
		}

		err = s.blogContentImageService.BulkDeleteHardByImageUrls(slice_image_urls, tx)
		if err != nil {
			tx.Rollback()
			return BlogUpdateResponse{}, err
		}
	}

	//todo: Commit Transaction
	if err := tx.Commit().Error; err != nil {
		err = fmt.Errorf("error commit transaction")
		return BlogUpdateResponse{}, err
	}

	response := BlogUpdateResponse{
		ID:              p.ID,
		Title:           p.Title,
		DescriptionHTML: p.DescriptionHTML,
		BannerUrl:       newFileURL,
		BannerFileName:  newFileName,
		Summary:         p.Summary,
		Status:          status,
	}

	return response, nil
}

func (s *service) DeleteBlog(id int) (Blog, error) {
	data, err := s.blogRepo.DeleteBlog(id)
	if err != nil {
		return Blog{}, err
	}
	return data, nil
}
