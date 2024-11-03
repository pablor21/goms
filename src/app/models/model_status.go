package models

type PublishStatus string //@Name PublishStatus

const (
	PublishStatusDraft     PublishStatus = "DRAFT"
	PublishStatusPublished PublishStatus = "PUBLISHED"
	PublishStatusArchived  PublishStatus = "ARCHIVED"
)
