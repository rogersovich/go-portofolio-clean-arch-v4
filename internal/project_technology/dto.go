package project_technology

type CreateProjectTechnologyRequest struct {
	ProjectID    int `json:"project_id" binding:"required"`
	TechnologyID int `json:"technology_id" binding:"required"`
}

type UpdateProjectTechnologyRequest struct {
	Id           int `json:"id" binding:"required"`
	ProjectID    int `json:"project_id" binding:"required"`
	TechnologyID int `json:"technology_id" binding:"required"`
}

type ProjectTechnologyResponse struct {
	ID           int    `json:"id"`
	ProjectID    int    `json:"project_id"`
	TechnologyID int    `json:"technology_id"`
	CreatedAt    string `json:"created_at"`
}

type ProjectTechnologyUpdateResponse struct {
	ID           int `json:"id"`
	ProjectID    int `json:"project_id"`
	TechnologyID int `json:"technology_id"`
}

type ProjectTechnologyDeleteRequest struct {
	ID int `json:"id" binding:"required"`
}

func ToProjectTechnologyResponse(p ProjectTechnology) ProjectTechnologyResponse {
	return ProjectTechnologyResponse{
		ID:           p.ID,
		ProjectID:    p.ProjectID,
		TechnologyID: p.TechnologyID,
		CreatedAt:    p.CreatedAt.Format("2006-01-02"),
	}
}

func ToProjectTechnologyUpdateResponse(p ProjectTechnology) ProjectTechnologyUpdateResponse {
	return ProjectTechnologyUpdateResponse{
		ID:           p.ID,
		ProjectID:    p.ProjectID,
		TechnologyID: p.TechnologyID,
	}
}
