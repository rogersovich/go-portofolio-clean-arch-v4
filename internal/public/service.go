package public

import (
	"sync"

	"github.com/rogersovich/go-portofolio-clean-arch-v4/internal/author"
	"github.com/rogersovich/go-portofolio-clean-arch-v4/pkg/utils"
)

type Service interface {
	GetAllPublicAuthors(params AuthorPublicParams) ([]AuthorPublicResponse, error)
	GetProfile() (ProfilePublicResponse, error)
	GetPublicBlogs(params BlogPublicParams) ([]BlogPublicResponse, error)
}

type service struct {
	authorService author.Service
	repo          Repository
}

func NewService(authorSvc author.Service, r Repository) Service {
	return &service{
		authorService: authorSvc,
		repo:          r,
	}
}

func (s *service) GetAllPublicAuthors(params AuthorPublicParams) ([]AuthorPublicResponse, error) {
	data, err := s.repo.FindAllAuthors(params)
	if err != nil {
		return nil, err
	}

	return data, nil
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
	rawBlogs, err := s.repo.GetPublicBlogs(params)

	// Create a slice to hold the formatted BlogPublicResponse data
	var blogResponses []BlogPublicResponse

	// Iterate through the raw blog data and map it to BlogPublicResponse
	for _, raw := range rawBlogs {
		blogResponse := s.MapBlogRawToResponse(raw)
		blogResponses = append(blogResponses, blogResponse)
	}

	if err != nil {
		return []BlogPublicResponse{}, err
	}
	return blogResponses, nil
}

func (s *service) MapBlogRawToResponse(raw BlogPublicRaw) BlogPublicResponse {
	// Mapping the BlogPublicRaw to BlogPublicResponse
	blogResponse := BlogPublicResponse{
		ID:              raw.ID,
		Title:           raw.Title,
		DescriptionHTML: raw.DescriptionHTML,
		BannerUrl:       raw.BannerUrl,
		BannerFileName:  raw.BannerFileName,
		Summary:         raw.Summary,
		Status:          raw.Status,
		PublishedAt:     raw.PublishedAt,
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
	var blogTopics []BlogPublicTopicResponse
	if raw.TopicID != 0 {
		blogTopics = append(blogTopics, BlogPublicTopicResponse{
			TopicID:   raw.TopicID,
			TopicName: raw.TopicName,
		})
		blogResponse.Topics = blogTopics
	} else {
		blogResponse.Topics = []BlogPublicTopicResponse{}
	}

	// Mapping Content Images (Assuming there can be multiple content images)
	var blogImages []BlogPublicContentImageResponse
	if raw.ContentImageID != 0 {
		blogImages = append(blogImages, BlogPublicContentImageResponse{
			ContentImageID:  raw.ContentImageID,
			ContentImageUrl: raw.ContentImageUrl,
			ContentFileName: raw.ContentImageFileName,
		})
		blogResponse.ContentImages = blogImages
	} else {
		blogResponse.ContentImages = []BlogPublicContentImageResponse{}
	}

	return blogResponse
}
