package models

import "github.com/google/uuid"

type DocumentsForEnlisting struct {
	DocumentTypeId string    `gorm:"type:uuid;default:uuid_generate_v4();primaryKey;index" json:"documentTypeId" example:"1"`
	DocumentType   string    `gorm:"column:document_type" json:"documentType" example:"pdf"` // pdf , doc , docx , jpg , jpeg , png , tiff , xls
	DocumentName   string    `gorm:"column:document_name" json:"documentName" example:"GST Non Enrollment Declaration"`
	Url            string    `gorm:"column:url" json:"url" example:"https://www.google.com/search?q=image+url&source=lnms&tbm=isch&sa=X&ved=2ahUKEwjX5bHvo7D-AhUixzgGHQ0jClQQ_AUoAXoECAEQAw&biw=1368&bih=800&dpr=1#imgrc=IflUsLqSUHqeoM"`
	BusinessId     uuid.UUID `gorm:"column:business_id" json:"businessId" example:"1"`
} //@name DocumentsForEnlisting

type DocumentsForEnlistingPatchRequest struct {
	BusinessId     string `json:"businessId" binding:"required" example:"1"`
	DocumentTypeId string `json:"documentTypeId" binding:"required" example:"1"`
	DocumentType   string `json:"documentType" binding:"required" example:"Passport"`
	DocumentName   string `json:"documentName" binding:"required" example:"Indian Passport"`
	Url            string `json:"url" binding:"required" example:"https://www.google.com"`
} // @name DocumentsForEnlistingPatchRequest
