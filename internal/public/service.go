package public

import (
	"context"
	"fmt"
	"slices"

	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

type Service interface {
	GetProfile() (ProfilePublicResponse, error)
	GetPublicBlogs(params BlogPublicParams) ([]BlogPublicResponse, int, error)
	GetPublicBlogBySlug(slug string) (SingleBlogPublicResponse, error)
	GetPublicTestimonials() ([]TestimonialPublicResponse, error)
	GetPublicTopics() ([]TopicPublicResponse, error)
	GetPublicProjects(params ProjectPublicParams) ([]ProjectPublicResponse, int, error)
	GetPublicProjectBySlug(slug string) (SingleProjectPublicResponse, error)
	GetPublicTechnologies() ([]TechnologyPublicResponse, error)
	GetPublicAuthors() ([]AuthorPublicResponse, error)
	GetPublicExperiences() ([]ExperiencesPublicResponse, error)
	UpdatePublicProjectStatistic(p ProjectStatisticUpdatePublicRequest) (ProjectStatisticUpdatePubblicResponse, error)
	UpdatePublicBlogStatistic(p BlogStatisticUpdatePublicRequest) (BlogStatisticUpdatePubblicResponse, error)
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
	// base context
	ctx := context.Background()

	// errgroup to collect errors and wait
	eg, ctx := errgroup.WithContext(ctx)

	// limit concurrency to 2
	sem := semaphore.NewWeighted(2)

	// placeholders for results
	var (
		about        AboutPublicResponse
		technologies []TechnologyProfilePublicResponse
		currentWork  CurrentWorkPublicResponse
		experiences  []ExperiencesPublicResponse
	)

	// helper: run fn(), store into targetVar
	run := func(name string, fn func() (interface{}, error), store func(interface{})) {
		eg.Go(func() error {
			// acquire a semaphore slot
			if err := sem.Acquire(ctx, 1); err != nil {
				return fmt.Errorf("acquire %s: %w", name, err)
			}
			defer sem.Release(1)

			res, err := fn()
			if err != nil {
				return fmt.Errorf("%s failed: %w", name, err)
			}

			store(res)
			return nil
		})
	}

	// schedule all calls
	run("GetAboutPublic",
		func() (interface{}, error) { return s.repo.GetAboutPublic() },
		func(r interface{}) { about = r.(AboutPublicResponse) },
	)

	run("GetTechnologiesPublic",
		func() (interface{}, error) { return s.repo.GetTechnologiesPublic() },
		func(r interface{}) { technologies = r.([]TechnologyProfilePublicResponse) },
	)

	run("GetCurrentWork",
		func() (interface{}, error) { return s.repo.GetCurrentWork() },
		func(r interface{}) { currentWork = r.(CurrentWorkPublicResponse) },
	)

	run("GetExperiencesPublic",
		func() (interface{}, error) { return s.repo.GetExperiencesPublic() },
		func(r interface{}) { experiences = r.([]ExperiencesPublicResponse) },
	)

	// wait for all or first error
	if err := eg.Wait(); err != nil {
		return ProfilePublicResponse{}, err
	}

	// post-process experiences: reformat FromDate
	var formatted []ExperiencesPublicResponse
	for _, exp := range experiences {
		formatted = append(formatted, ExperiencesPublicResponse{
			Position:          exp.Position,
			CompanyName:       exp.CompanyName,
			WorkType:          exp.WorkType,
			Country:           exp.Country,
			City:              exp.City,
			CompWebsiteUrl:    exp.CompWebsiteUrl,
			SummaryHTML:       exp.SummaryHTML,
			FromDate:          exp.FromDate,
			ToDate:            exp.ToDate,
			CompImageUrl:      exp.CompImageUrl,
			CompImageFileName: exp.CompImageFileName,
			IsCurrent:         exp.IsCurrent,
		})
	}

	// assemble final response
	return ProfilePublicResponse{
		About:        about,
		Technologies: technologies,
		CurrentWork:  currentWork,
		Experiences:  formatted,
	}, nil
}

func (s *service) GetPublicExperiences() ([]ExperiencesPublicResponse, error) {
	datas, err := s.repo.GetExperiencesPublic()
	if err != nil {
		return []ExperiencesPublicResponse{}, err
	}

	var experiences_formatted []ExperiencesPublicResponse
	for _, experience := range datas {
		experiences_formatted = append(experiences_formatted, ExperiencesPublicResponse{
			Position:          experience.Position,
			CompanyName:       experience.CompanyName,
			WorkType:          experience.WorkType,
			Country:           experience.Country,
			City:              experience.City,
			CompWebsiteUrl:    experience.CompWebsiteUrl,
			SummaryHTML:       experience.SummaryHTML,
			FromDate:          experience.FromDate,
			ToDate:            experience.ToDate,
			CompImageUrl:      experience.CompImageUrl,
			CompImageFileName: experience.CompImageFileName,
			IsCurrent:         experience.IsCurrent,
		})
	}

	return experiences_formatted, nil
}

func (s *service) GetPublicBlogs(params BlogPublicParams) ([]BlogPublicResponse, int, error) {
	//todo: Get Raw Paginate Blogs
	rawPaginateBlogs, total_records, err := s.repo.GetRawPublicPaginateBlogs(params)

	if err != nil {
		return []BlogPublicResponse{}, 0, err
	}

	if len(rawPaginateBlogs) == 0 {
		return []BlogPublicResponse{}, 0, nil
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
		return []BlogPublicResponse{}, 0, err
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
		return []BlogPublicResponse{}, 0, err
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
		// Check if raw.ID already exists in blogResponses
		exists := false
		for _, existingResponse := range blogResponses {
			if existingResponse.ID == raw.ID {
				exists = true
				break
			}
		}

		// If it doesn't exist, map the raw data to BlogPublicResponse and append it
		if !exists {
			mapBlogTopic := mappedBlogTopics[raw.ID]
			blogResponse := s.MapBlogRawToResponse(raw, mapBlogTopic)
			blogResponses = append(blogResponses, blogResponse)
		}
	}

	return blogResponses, total_records, nil
}

func (s *service) MapBlogRawToResponse(raw BlogPublicRaw, blogTopics []BlogTopicPublicRaw) BlogPublicResponse {
	// Mapping the BlogPublicRaw to BlogPublicResponse
	blogResponse := BlogPublicResponse{
		ID:              raw.ID,
		Title:           raw.Title,
		BannerUrl:       raw.BannerUrl,
		BannerFileName:  raw.BannerFileName,
		Summary:         raw.Summary,
		DescriptionHTML: raw.DescriptionHTML,
		Status:          raw.Status,
		Slug:            raw.Slug,
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

func (s *service) GetPublicTopics() ([]TopicPublicResponse, error) {
	datas, err := s.repo.GetPublicTopics()
	if err != nil {
		return []TopicPublicResponse{}, err
	}

	return datas, nil
}

func (s *service) GetPublicProjects(params ProjectPublicParams) ([]ProjectPublicResponse, int, error) {
	//todo: Get Raw Paginate Project
	rawPaginateProjects, total_record, err := s.repo.GetRawPublicPaginateProjects(params)

	if err != nil {
		return []ProjectPublicResponse{}, 0, err
	}

	if len(rawPaginateProjects) == 0 {
		return []ProjectPublicResponse{}, 0, nil
	}

	//?: Slice unique paginate Project ids
	var uniquePaginateProjectIDs []int

	for _, data := range rawPaginateProjects {
		if !slices.Contains(uniquePaginateProjectIDs, data.ID) {
			uniquePaginateProjectIDs = append(uniquePaginateProjectIDs, data.ID)
		}
	}

	//todo: Get Raw Project Technlogies
	rawProjectTechnologies, err := s.repo.GetRawPublicProjectTechnologies(params, uniquePaginateProjectIDs)

	if err != nil {
		return []ProjectPublicResponse{}, 0, err
	}

	//?: Map Project Technlogies
	mappedProjectTechnologies := make(map[int][]ProjectTechnologyPublicRaw)
	for _, tech := range rawProjectTechnologies {
		mappedProjectTechnologies[tech.ProjectID] = append(mappedProjectTechnologies[tech.ProjectID], tech)
	}

	//? Create a slice to hold the formatted ProjectPublicResponse data
	var projectResponses []ProjectPublicResponse

	for _, raw := range rawPaginateProjects {
		mapProjectTech := mappedProjectTechnologies[raw.ID]
		projectResponse := s.MapProjectRawToResponse(raw, mapProjectTech)

		// Append the result to the slice
		projectResponses = append(projectResponses, projectResponse)
	}

	return projectResponses, total_record, nil
}

func (s *service) MapProjectRawToResponse(raw ProjectPaginatePublicRaw, projectTechnologies []ProjectTechnologyPublicRaw) ProjectPublicResponse {
	// Mapping the BlogPublicRaw to ProjectPublicResponse
	projectResponse := ProjectPublicResponse{
		ID:            raw.ID,
		Title:         raw.Title,
		Summary:       raw.Summary,
		ImageURL:      raw.ImageURL,
		ImageFileName: raw.ImageFileName,
		RepositoryURL: raw.RepositoryURL,
		PublishedAt:   raw.PublishedAt,
		Slug:          raw.Slug,
	}

	// Mapping the Topics (Assuming topics is an array, you may want to append more topic records here)
	if len(projectTechnologies) != 0 {
		for _, tech := range projectTechnologies {
			projectResponse.Technologies = append(projectResponse.Technologies, ProjectTechnologyPublicResponse{
				TechID:           tech.TechID,
				TechName:         tech.TechName,
				TechLogoURL:      tech.TechLogoURL,
				TechLogoFileName: tech.TechLogoFileName,
				TechLink:         tech.TechLink,
			})
		}
	} else {
		projectResponse.Technologies = []ProjectTechnologyPublicResponse{}
	}

	return projectResponse
}

func (s *service) GetPublicProjectBySlug(slug string) (SingleProjectPublicResponse, error) {
	rawData, err := s.repo.GetPublicProjectBySlug(slug)

	datas := s.MapSingleProjectRawToResponse(rawData)

	if err != nil {
		return SingleProjectPublicResponse{}, err
	}
	return datas, nil
}

func (s *service) MapSingleProjectRawToResponse(rawData []SingleProjectPublicRaw) SingleProjectPublicResponse {
	// Mapping result
	projectMap := map[int]*SingleProjectPublicResponse{}

	for _, row := range rawData {
		projectID := int(row.ID)

		//? "Comma-ok" itu fitur spesial
		_, exists := projectMap[projectID]
		if !exists {
			var publishedAtPointer *string
			if row.PublishedAt != nil {
				formattedPublishedAt := row.PublishedAt.Format("2006-01-02 15:04:05")
				publishedAtPointer = &formattedPublishedAt
			}

			var projectStatistic *ProjectPublicStatisticResponse
			if row.StatisticID != 0 {
				projectStatistic = &ProjectPublicStatisticResponse{
					StatisticID:    row.StatisticID,
					StatisticLikes: row.StatisticLikes,
					StatisticViews: row.StatisticViews,
					StatisticType:  row.StatisticType,
				}
			}

			projectMap[projectID] = &SingleProjectPublicResponse{
				ID:            projectID,
				Title:         row.Title,
				Description:   row.Description,
				ImageUrl:      row.ImageUrl,
				ImageFileName: row.ImageFileName,
				RepositoryUrl: row.RepositoryUrl,
				Summary:       row.Summary,
				Status:        row.Status,
				Slug:          row.Slug,
				PublishedAt:   publishedAtPointer,
				Statistic:     projectStatistic,
				ContentImages: []ProjectPublicContentImageResponse{},
				Technologies:  []ProjectTechnologyPublicResponse{},
			}
		}

		if row.TechID != 0 {
			seen := make(map[int]bool)
			for _, project := range projectMap[projectID].Technologies {
				seen[project.TechID] = true
			}

			if !seen[row.TechID] {
				projectMap[projectID].Technologies = append(projectMap[projectID].Technologies, ProjectTechnologyPublicResponse{
					TechID:      row.TechID,
					TechName:    row.TechName,
					TechLogoURL: row.TechLogoURL,
					TechLink:    row.TechLink,
				})
			}
		}

		if row.ContentImageID != 0 {
			seen := make(map[int]bool)
			for _, img := range projectMap[projectID].ContentImages {
				seen[img.ContentImageID] = true
			}

			if !seen[row.ContentImageID] {
				projectMap[projectID].ContentImages = append(projectMap[projectID].ContentImages, ProjectPublicContentImageResponse{
					ContentImageID:       row.ContentImageID,
					ContentImageUrl:      row.ContentImageUrl,
					ContentImageFileName: row.ContentImageFileName,
				})
			}
		}
	}

	// Convert Map to Struct
	var result SingleProjectPublicResponse
	for _, v := range projectMap {
		result = *v
		break
	}

	return result
}

func (s *service) GetPublicTechnologies() ([]TechnologyPublicResponse, error) {
	datas, err := s.repo.GetPublicTechnologies()
	if err != nil {
		return []TechnologyPublicResponse{}, err
	}

	for i := range datas {
		if datas[i].IsMajor == "0" {
			datas[i].IsMajor = "N"
		} else {
			datas[i].IsMajor = "Y"
		}
	}

	return datas, nil
}

func (s *service) GetPublicAuthors() ([]AuthorPublicResponse, error) {
	datas, err := s.repo.GetPublicAuthors()
	if err != nil {
		return []AuthorPublicResponse{}, err
	}

	return datas, nil
}

func (s *service) GetProjectById(id int) (ProjectByIdResponse, error) {
	data, err := s.repo.FindProjectById(id)
	if err != nil {
		return ProjectByIdResponse{}, err
	}
	return data, nil
}

func (s *service) UpdatePublicProjectStatistic(p ProjectStatisticUpdatePublicRequest) (ProjectStatisticUpdatePubblicResponse, error) {
	project, err := s.GetProjectById(p.ProjectID)
	if err != nil {
		return ProjectStatisticUpdatePubblicResponse{}, err
	}

	payload := ProjectStatisticUpdatePublicDTO{
		ProjectID:    p.ProjectID,
		ProjectTitle: project.Title,
		StatisticID:  p.StatisticID,
		Likes:        p.Likes,
		Views:        p.Views,
		Type:         p.Type,
	}

	data, err := s.repo.UpdatePublicProjectStatistic(payload)
	if err != nil {
		return ProjectStatisticUpdatePubblicResponse{}, err
	}
	return data, nil
}

func (s *service) GetBlogById(id int) (BlogByIdResponse, error) {
	data, err := s.repo.FindBlogById(id)
	if err != nil {
		return BlogByIdResponse{}, err
	}
	return data, nil
}

func (s *service) UpdatePublicBlogStatistic(p BlogStatisticUpdatePublicRequest) (BlogStatisticUpdatePubblicResponse, error) {
	blog, err := s.GetBlogById(p.BlogID)
	if err != nil {
		return BlogStatisticUpdatePubblicResponse{}, err
	}

	payload := BlogStatisticUpdatePublicDTO{
		BlogID:      p.BlogID,
		Title:       blog.Title,
		StatisticID: p.StatisticID,
		Likes:       p.Likes,
		Views:       p.Views,
		Type:        p.Type,
	}

	data, err := s.repo.UpdatePublicBlogStatistic(payload)
	if err != nil {
		return BlogStatisticUpdatePubblicResponse{}, err
	}
	return data, nil
}
