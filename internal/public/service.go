package public

import (
	"slices"
	"sync"

	"github.com/rogersovich/go-portofolio-clean-arch-v4/pkg/utils"
)

type Service interface {
	GetProfile() (ProfilePublicResponse, error)
	GetPublicBlogs(params BlogPublicParams) ([]BlogPublicResponse, error)
	GetPublicBlogBySlug(slug string) (SingleBlogPublicResponse, error)
	GetPublicTestimonials() ([]TestimonialPublicResponse, error)
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{
		repo: r,
	}
}

func (s *service) GetProfile() (ProfilePublicResponse, error) {
	// Create channels to collect results and errors
	aboutCh := make(chan AboutPublicResponse, 1)
	technologiesCh := make(chan []TechnologyProfilePublicResponse, 1)
	currentWorkCh := make(chan CurrentWorkPublicResponse, 1)
	experiencesCh := make(chan []ExperiencesPublicResponse, 1)

	// Create a WaitGroup to ensure all goroutines finish
	var wg sync.WaitGroup

	// Error channel to capture any errors from goroutines
	errCh := make(chan error, 1)

	// Goroutines to call repository functions concurrently
	wg.Add(4) // We are launching 4 goroutines

	// Goroutines to call repository functions concurrently
	go func() {
		defer wg.Done()
		about, err := s.repo.GetAboutPublic()
		if err != nil {
			errCh <- err
			return
		}
		aboutCh <- about
	}()

	go func() {
		defer wg.Done()
		technologies, err := s.repo.GetTechnologiesPublic()
		if err != nil {
			errCh <- err
			return
		}
		technologiesCh <- technologies
	}()

	go func() {
		defer wg.Done()
		currentWork, err := s.repo.GetCurrentWork()
		if err != nil {
			errCh <- err
			return
		}
		currentWorkCh <- currentWork
	}()

	go func() {
		defer wg.Done()
		experiences, err := s.repo.GetExperiencesPublic()
		if err != nil {
			errCh <- err
			return
		}
		experiencesCh <- experiences
	}()

	// Wait for all goroutines to finish
	wg.Wait()

	// Handle any errors from the goroutines
	select {
	case err := <-errCh:
		return ProfilePublicResponse{}, err
	default:
		// No errors, continue processing
	}

	// Collect results from channels
	about := <-aboutCh
	technologies := <-technologiesCh
	currentWork := <-currentWorkCh
	experiences := <-experiencesCh

	var experiences_formatted []ExperiencesPublicResponse
	for _, experience := range experiences {
		fromDateTime, _ := utils.ParseStringToTime(experience.FromDate, "2006-01-02T15:04:05-07:00")
		experiences_formatted = append(experiences_formatted, ExperiencesPublicResponse{
			Position:          experience.Position,
			CompanyName:       experience.CompanyName,
			WorkType:          experience.WorkType,
			Country:           experience.Country,
			City:              experience.City,
			CompWebsiteUrl:    experience.CompWebsiteUrl,
			SummaryHTML:       experience.SummaryHTML,
			FromDate:          fromDateTime.Format("2006-01-02"),
			ToDate:            experience.ToDate,
			CompImageUrl:      experience.CompImageUrl,
			CompImageFileName: experience.CompImageFileName,
			IsCurrent:         experience.IsCurrent,
		})
	}

	data := ProfilePublicResponse{
		About:        about,
		Technologies: technologies,
		CurrentWork:  currentWork,
		Experiences:  experiences_formatted,
	}
	return data, nil
}

func (s *service) GetPublicBlogs(params BlogPublicParams) ([]BlogPublicResponse, error) {
	//todo: Get Raw Paginate Blogs
	rawPaginateBlogs, err := s.repo.GetRawPublicPaginateBlogs(params)

	if err != nil {
		return []BlogPublicResponse{}, err
	}

	if len(rawPaginateBlogs) == 0 {
		return []BlogPublicResponse{}, nil
	}

	//?: Slice unique paginate blog ids
	var uniquePaginateBlogIDs []int

	for _, data := range rawPaginateBlogs {
		if !slices.Contains(uniquePaginateBlogIDs, data.ID) {
			uniquePaginateBlogIDs = append(uniquePaginateBlogIDs, data.ID)
		}
	}

	//todo: Get Raw Blogs
	rawBlogs, err := s.repo.GetRawPublicBlogs(params, uniquePaginateBlogIDs)

	if err != nil {
		return []BlogPublicResponse{}, err
	}

	//?: Slice unique blog ids
	var uniqueBlogIDs []int

	for _, data := range rawBlogs {
		if !slices.Contains(uniqueBlogIDs, data.ID) {
			uniqueBlogIDs = append(uniqueBlogIDs, data.ID)
		}
	}

	//todo: Get Raw Blog Topics
	rawBlogTopics, err := s.repo.GetRawPublicBlogTopics(params, uniqueBlogIDs)

	if err != nil {
		return []BlogPublicResponse{}, err
	}

	//?: Map Blog Topics
	mappedBlogTopics := make(map[int][]BlogTopicPublicRaw)
	for _, topic := range rawBlogTopics {
		mappedBlogTopics[topic.BlogID] = append(mappedBlogTopics[topic.BlogID], topic)
	}

	// Create a slice to hold the formatted BlogPublicResponse data
	var blogResponses []BlogPublicResponse

	// Iterate through the raw blog data and map it to BlogPublicResponse
	for _, raw := range rawBlogs {
		// Map the raw data to BlogPublicResponse
		mapBlogTopic := mappedBlogTopics[raw.ID]
		blogResponse := s.MapBlogRawToResponse(raw, mapBlogTopic)

		// Append the result to the slice
		blogResponses = append(blogResponses, blogResponse)
	}

	return blogResponses, nil
}

func (s *service) MapBlogRawToResponse(raw BlogPublicRaw, blogTopics []BlogTopicPublicRaw) BlogPublicResponse {
	// Mapping the BlogPublicRaw to BlogPublicResponse
	blogResponse := BlogPublicResponse{
		ID:             raw.ID,
		Title:          raw.Title,
		BannerUrl:      raw.BannerUrl,
		BannerFileName: raw.BannerFileName,
		Summary:        raw.Summary,
		Status:         raw.Status,
		Slug:           raw.Slug,
		PublishedAt:    raw.PublishedAt,
	}

	// Mapping Author (Assuming Author info comes from the same raw data)
	var blogAuthor *BlogPublicAuthorResponse
	if raw.AuthorID != 0 {
		blogAuthor = &BlogPublicAuthorResponse{
			AuthorID:   raw.AuthorID,
			AuthorName: raw.AuthorName,
		}
		blogResponse.Author = blogAuthor
	} else {
		blogResponse.Author = nil
	}

	// Mapping Reading Time (You might have a separate function or struct for reading time)
	var blogReadingTime *BlogPublicReadingTimeResponse
	if raw.ReadingTimeID != 0 {
		blogReadingTime = &BlogPublicReadingTimeResponse{
			ReadingTimeID:               raw.ReadingTimeID,
			ReadingTimeMinutes:          raw.ReadingTimeMinutes,
			ReadingTimeTextLength:       raw.ReadingTimeTextLength,
			ReadingTimeEstimatedSeconds: raw.ReadingTimeEstimatedSeconds,
			ReadingTimeWordCount:        raw.ReadingTimeWordCount,
			ReadingTimeType:             raw.ReadingTimeType,
		}

		blogResponse.ReadingTime = blogReadingTime
	} else {
		blogResponse.ReadingTime = nil
	}

	// Mapping the Statistic (Single mapping for statistic-related fields)
	var blogStatistic *BlogPublicStatisticResponse
	if raw.StatisticID != 0 {
		blogStatistic = &BlogPublicStatisticResponse{
			StatisticID:    raw.StatisticID,
			StatisticLikes: raw.StatisticLikes,
			StatisticViews: raw.StatisticViews,
			StatisticType:  raw.StatisticType,
		}
		blogResponse.Statistic = blogStatistic
	} else {
		blogResponse.Statistic = nil
	}

	// Mapping the Topics (Assuming topics is an array, you may want to append more topic records here)
	if len(blogTopics) != 0 {
		for _, topic := range blogTopics {
			blogResponse.Topics = append(blogResponse.Topics, BlogPublicTopicResponse{
				TopicID:   topic.TopicID,
				TopicName: topic.TopicName,
			})
		}
	} else {
		blogResponse.Topics = []BlogPublicTopicResponse{}
	}

	return blogResponse
}

func (s *service) GetPublicBlogBySlug(slug string) (SingleBlogPublicResponse, error) {
	rawData, err := s.repo.GetPublicBlogBySlug(slug)

	datas := s.MapSingleBlogRawToResponse(rawData)

	if err != nil {
		return SingleBlogPublicResponse{}, err
	}
	return datas, nil
}

func (s *service) MapSingleBlogRawToResponse(rawData []SingleBlogPublicRaw) SingleBlogPublicResponse {
	// Mapping result
	blogMap := map[int]*SingleBlogPublicResponse{}

	for _, row := range rawData {
		blogID := int(row.ID)

		//? "Comma-ok" itu fitur spesial
		_, exists := blogMap[blogID]
		if !exists {
			var publishedAtPointer *string
			if row.PublishedAt != nil {
				formattedPublishedAt := row.PublishedAt.Format("2006-01-02 15:04:05")
				publishedAtPointer = &formattedPublishedAt
			}

			var blogAuthor *BlogPublicAuthorResponse
			if row.AuthorID != 0 {
				blogAuthor = &BlogPublicAuthorResponse{
					AuthorID:   row.AuthorID,
					AuthorName: row.AuthorName,
				}
			}

			var blogReadingTime *BlogPublicReadingTimeResponse
			if row.ReadingTimeID != 0 {
				blogReadingTime = &BlogPublicReadingTimeResponse{
					ReadingTimeID:               row.ReadingTimeID,
					ReadingTimeMinutes:          row.ReadingTimeMinutes,
					ReadingTimeTextLength:       row.ReadingTimeTextLength,
					ReadingTimeEstimatedSeconds: row.ReadingTimeEstimatedSeconds,
					ReadingTimeWordCount:        row.ReadingTimeWordCount,
					ReadingTimeType:             row.ReadingTimeType,
				}
			}

			var blogStatistic *BlogPublicStatisticResponse
			if row.StatisticID != 0 {
				blogStatistic = &BlogPublicStatisticResponse{
					StatisticID:    row.StatisticID,
					StatisticLikes: row.StatisticLikes,
					StatisticViews: row.StatisticViews,
					StatisticType:  row.StatisticType,
				}
			}

			blogMap[blogID] = &SingleBlogPublicResponse{
				ID:              blogID,
				Title:           row.Title,
				DescriptionHTML: row.DescriptionHTML,
				BannerUrl:       row.BannerUrl,
				BannerFileName:  row.BannerFileName,
				Summary:         row.Summary,
				Status:          row.Status,
				Slug:            row.Slug,
				Author:          blogAuthor,
				ReadingTime:     blogReadingTime,
				Statistic:       blogStatistic,
				PublishedAt:     publishedAtPointer,
				Topics:          []BlogPublicTopicResponse{},
				ContentImages:   []BlogPublicContentImageResponse{},
			}
		}

		if row.TopicID != 0 {
			seen := make(map[int]bool)
			for _, topic := range blogMap[blogID].Topics {
				seen[topic.TopicID] = true
			}

			if !seen[row.TopicID] {
				blogMap[blogID].Topics = append(blogMap[blogID].Topics, BlogPublicTopicResponse{
					TopicID:   row.TopicID,
					TopicName: row.TopicName,
				})
			}
		}

		if row.ContentImageID != 0 {
			seen := make(map[int]bool)
			for _, img := range blogMap[blogID].ContentImages {
				seen[img.ContentImageID] = true
			}

			if !seen[row.ContentImageID] {
				blogMap[blogID].ContentImages = append(blogMap[blogID].ContentImages, BlogPublicContentImageResponse{
					ContentImageID:       row.ContentImageID,
					ContentImageUrl:      row.ContentImageUrl,
					ContentImageFileName: row.ContentImageFileName,
				})
			}
		}
	}

	// Convert Map to Struct
	var result SingleBlogPublicResponse
	for _, v := range blogMap {
		result = *v
		break
	}

	return result
}

func (s *service) GetPublicTestimonials() ([]TestimonialPublicResponse, error) {
	datas, err := s.repo.GetPublicTestimonials()
	if err != nil {
		return []TestimonialPublicResponse{}, err
	}

	return datas, nil
}
