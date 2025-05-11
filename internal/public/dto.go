package public

import "time"

type AuthorPublicResponse struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	AvatarUrl      string `json:"avatar_url"`
	AvatarFileName string `json:"avatar_file_name"`
}

type AuthorPublicParams struct {
	Page  int
	Limit int
	Sort  string
	Order string
	Name  string
}

type AboutPublicResponse struct {
	ID              int    `json:"about_id"`
	Title           string `json:"title"`
	DescriptionHTML string `json:"description_html"`
	AvatarUrl       string `json:"avatar_url"`
	AvatarFileName  string `json:"avatar_file_name"`
}

type CurrentWorkPublicResponse struct {
	Position       string `json:"position"`
	CompanyName    string `json:"company_name"`
	WorkType       string `json:"work_type"`
	Country        string `json:"country"`
	City           string `json:"city"`
	CompWebsiteUrl string `json:"comp_website_url"`
}

type TechnologyProfilePublicResponse struct {
	ID           int    `json:"tech_id"`
	Name         string `json:"name"`
	LogoUrl      string `json:"logo_url"`
	LogoFileName string `json:"logo_file_name"`
}

type ExperiencesPublicResponse struct {
	Position          string `json:"position"`
	CompanyName       string `json:"company_name"`
	WorkType          string `json:"work_type"`
	Country           string `json:"country"`
	City              string `json:"city"`
	CompWebsiteUrl    string `json:"comp_website_url"`
	SummaryHTML       string `json:"summary_html"`
	FromDate          string `json:"from_date"`
	ToDate            string `json:"to_date"`
	CompImageUrl      string `json:"comp_image_url"`
	CompImageFileName string `json:"comp_image_file_name"`
	IsCurrent         string `json:"is_current"`
}

type ProfilePublicResponse struct {
	About        AboutPublicResponse               `json:"about"`
	CurrentWork  CurrentWorkPublicResponse         `json:"current_work"`
	Technologies []TechnologyProfilePublicResponse `json:"technologies"`
	Experiences  []ExperiencesPublicResponse       `json:"experiences"`
}

type BlogPublicParams struct {
	Page   int    `binding:"required"`
	Limit  int    `binding:"required"`
	Sort   string `binding:"required,oneof=published_at id views"`
	Order  string `binding:"required"`
	Search string
	Topics []int
}

type BlogPaginatePublicRaw struct {
	ID             int    `json:"id"`
	Title          string `json:"title"`
	StatisticID    int    `json:"statistic_id"`
	StatisticLikes int    `json:"statistic_likes"`
	StatisticViews int    `json:"statistic_views"`
	StatisticType  string `json:"statistic_type"`
}

type BlogPublicRaw struct {
	ID                          int        `json:"id"`
	Title                       string     `json:"title"`
	BannerUrl                   string     `json:"banner_url"`
	BannerFileName              string     `json:"banner_file_name"`
	Summary                     string     `json:"summary"`
	Status                      string     `json:"status"`
	Slug                        string     `json:"slug"`
	PublishedAt                 *time.Time `json:"published_at"`
	AuthorID                    int        `json:"author_id"`
	AuthorName                  string     `json:"author_name"`
	ReadingTimeID               int        `json:"reading_time_id"`
	ReadingTimeMinutes          int        `json:"reading_time_minutes"`
	ReadingTimeTextLength       int        `json:"reading_time_text_length"`
	ReadingTimeEstimatedSeconds float64    `json:"reading_time_estimated_seconds"`
	ReadingTimeWordCount        int        `json:"reading_time_word_count"`
	ReadingTimeType             string     `json:"reading_time_type"`
	StatisticID                 int        `json:"statistic_id"`
	StatisticLikes              int        `json:"statistic_likes"`
	StatisticViews              int        `json:"statistic_views"`
	StatisticType               string     `json:"statistic_type"`
	TopicID                     int        `json:"topic_id"`
	TopicName                   string     `json:"topic_name"`
}

type BlogTopicPublicRaw struct {
	BlogID    int    `json:"blog_id"`
	TopicID   int    `json:"topic_id"`
	TopicName string `json:"topic_name"`
}

type BlogPublicAuthorResponse struct {
	AuthorID   int    `json:"id"`
	AuthorName string `json:"name"`
}

type BlogPublicReadingTimeResponse struct {
	ReadingTimeID               int     `json:"id"`
	ReadingTimeMinutes          int     `json:"minutes"`
	ReadingTimeTextLength       int     `json:"text_length"`
	ReadingTimeEstimatedSeconds float64 `json:"estimated_seconds"`
	ReadingTimeWordCount        int     `json:"word_count"`
	ReadingTimeType             string  `json:"type"`
}

type BlogPublicStatisticResponse struct {
	StatisticID    int    `json:"id"`
	StatisticLikes int    `json:"likes"`
	StatisticViews int    `json:"views"`
	StatisticType  string `json:"type"`
}

type BlogPublicTopicResponse struct {
	TopicID   int    `json:"id"`
	TopicName string `json:"name"`
}

type BlogPublicContentImageResponse struct {
	ContentImageID       int    `json:"id"`
	ContentImageUrl      string `json:"url"`
	ContentImageFileName string `json:"file_name"`
}

type BlogPublicResponse struct {
	ID             int                            `json:"id"`
	Title          string                         `json:"title"`
	BannerUrl      string                         `json:"banner_url"`
	BannerFileName string                         `json:"banner_file_name"`
	Summary        string                         `json:"summary"`
	Status         string                         `json:"status"`
	Slug           string                         `json:"slug"`
	PublishedAt    *time.Time                     `json:"published_at"`
	Author         *BlogPublicAuthorResponse      `json:"author"`
	ReadingTime    *BlogPublicReadingTimeResponse `json:"reading_time"`
	Statistic      *BlogPublicStatisticResponse   `json:"statistic"`
	Topics         []BlogPublicTopicResponse      `json:"topics"`
}

type SingleBlogPublicRaw struct {
	ID                          int        `json:"id"`
	Title                       string     `json:"title"`
	DescriptionHTML             string     `json:"description_html"`
	BannerUrl                   string     `json:"banner_url"`
	BannerFileName              string     `json:"banner_file_name"`
	Summary                     string     `json:"summary"`
	Status                      string     `json:"status"`
	Slug                        string     `json:"slug"`
	PublishedAt                 *time.Time `json:"published_at"`
	AuthorID                    int        `json:"author_id"`
	AuthorName                  string     `json:"author_name"`
	ReadingTimeID               int        `json:"reading_time_id"`
	ReadingTimeMinutes          int        `json:"reading_time_minutes"`
	ReadingTimeTextLength       int        `json:"reading_time_text_length"`
	ReadingTimeEstimatedSeconds float64    `json:"reading_time_estimated_seconds"`
	ReadingTimeWordCount        int        `json:"reading_time_word_count"`
	ReadingTimeType             string     `json:"reading_time_type"`
	StatisticID                 int        `json:"statistic_id"`
	StatisticLikes              int        `json:"statistic_likes"`
	StatisticViews              int        `json:"statistic_views"`
	StatisticType               string     `json:"statistic_type"`
	ContentImageID              int        `json:"content_image_id"`
	ContentImageUrl             string     `json:"content_image_url"`
	ContentImageFileName        string     `json:"content_image_file_name"`
	TopicID                     int        `json:"topic_id"`
	TopicName                   string     `json:"topic_name"`
}

type SingleBlogPublicResponse struct {
	ID              int                              `json:"id"`
	Title           string                           `json:"title"`
	DescriptionHTML string                           `json:"description_html"`
	BannerUrl       string                           `json:"banner_url"`
	BannerFileName  string                           `json:"banner_file_name"`
	Summary         string                           `json:"summary"`
	Status          string                           `json:"status"`
	Slug            string                           `json:"slug"`
	PublishedAt     *string                          `json:"published_at"`
	Author          *BlogPublicAuthorResponse        `json:"author"`
	ReadingTime     *BlogPublicReadingTimeResponse   `json:"reading_time"`
	Statistic       *BlogPublicStatisticResponse     `json:"statistic"`
	Topics          []BlogPublicTopicResponse        `json:"topics"`
	ContentImages   []BlogPublicContentImageResponse `json:"content_image"`
}

type TestimonialPublicResponse struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Via       *string `json:"via"`
	Role      *string `json:"role"`
	WorkingAt *string `json:"working_at"`
	IsUsed    string  `json:"is_used"`
	CreatedAt string  `json:"created_at"`
}

type TopicPublicResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type ProjectPublicParams struct {
	Page   int    `binding:"required"`
	Limit  int    `binding:"required"`
	Sort   string `binding:"required,oneof=published_at created_at id"`
	Order  string `binding:"required"`
	Search string
}

type ProjectPaginatePublicRaw struct {
	ID            int        `json:"id"`
	Title         string     `json:"title"`
	Summary       string     `json:"summary"`
	ImageURL      string     `json:"image_url"`
	ImageFileName string     `json:"image_file_name"`
	RepositoryURL *string    `json:"repository_url"`
	PublishedAt   *time.Time `json:"published_at"`
	Slug          string     `json:"slug"`
}

type ProjectTechnologyPublicRaw struct {
	ProjectID   int    `json:"project_id"`
	TechID      int    `json:"tech_id"`
	TechName    string `json:"tech_name"`
	TechLogoURL string `json:"tech_logo_url"`
	TechLink    string `json:"tech_link"`
}

type ProjectTechnologyPublicResponse struct {
	TechID      int    `json:"tech_id"`
	TechName    string `json:"tech_name"`
	TechLogoURL string `json:"tech_logo_url"`
	TechLink    string `json:"tech_link"`
}

type ProjectPublicContentImageResponse struct {
	ContentImageID       int    `json:"id"`
	ContentImageUrl      string `json:"url"`
	ContentImageFileName string `json:"file_name"`
}

type ProjectPublicStatisticResponse struct {
	StatisticID    int    `json:"id"`
	StatisticLikes int    `json:"likes"`
	StatisticViews int    `json:"views"`
	StatisticType  string `json:"type"`
}

type ProjectPublicResponse struct {
	ID            int                               `json:"id"`
	Title         string                            `json:"title"`
	Summary       string                            `json:"summary"`
	ImageURL      string                            `json:"image_url"`
	ImageFileName string                            `json:"image_file_name"`
	RepositoryURL *string                           `json:"repository_url"`
	Slug          string                            `json:"slug"`
	PublishedAt   *time.Time                        `json:"published_at"`
	Technologies  []ProjectTechnologyPublicResponse `json:"technologies"`
}

type SingleProjectPublicRaw struct {
	ID                   int        `json:"id"`
	Title                string     `json:"title"`
	Description          string     `json:"description"`
	ImageUrl             string     `json:"image_url"`
	ImageFileName        string     `json:"image_file_name"`
	RepositoryUrl        *string    `json:"repository_url"`
	Summary              string     `json:"summary"`
	Status               string     `json:"status"`
	Slug                 string     `json:"slug"`
	PublishedAt          *time.Time `json:"published_at"`
	StatisticID          int        `json:"statistic_id"`
	StatisticLikes       int        `json:"statistic_likes"`
	StatisticViews       int        `json:"statistic_views"`
	StatisticType        string     `json:"statistic_type"`
	ContentImageID       int        `json:"content_image_id"`
	ContentImageUrl      string     `json:"content_image_url"`
	ContentImageFileName string     `json:"content_image_file_name"`
	TechID               int        `json:"tech_id"`
	TechName             string     `json:"tech_name"`
	TechLogoURL          string     `json:"tech_logo_url"`
	TechLink             string     `json:"tech_link"`
}

type SingleProjectPublicResponse struct {
	ID            int                                 `json:"id"`
	Title         string                              `json:"title"`
	Description   string                              `json:"description"`
	ImageUrl      string                              `json:"image_url"`
	ImageFileName string                              `json:"image_file_name"`
	RepositoryUrl *string                             `json:"repository_url"`
	Summary       string                              `json:"summary"`
	Status        string                              `json:"status"`
	Slug          string                              `json:"slug"`
	PublishedAt   *string                             `json:"published_at"`
	Statistic     *ProjectPublicStatisticResponse     `json:"statistic"`
	ContentImages []ProjectPublicContentImageResponse `json:"content_image"`
	Technologies  []ProjectTechnologyPublicResponse   `json:"topics"`
}

type TechnologyPublicResponse struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	LogoURL      string `json:"logo_url"`
	LogoFileName string `json:"logo_file_name"`
	IsMajor      string `json:"is_major"`
}
