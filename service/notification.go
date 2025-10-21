package service

import (
	"hrms/model"
	"hrms/resource"
	"log"

	"github.com/gin-gonic/gin"
)

func GetNotificationByTitle(c *gin.Context, noticeTitle string, start int, limit int) ([]*model.Notification, int64, error) {
	var notifications []*model.Notification
	var err error
	db := resource.HrmsDB(c)
	if db == nil {
		log.Printf("GetNotificationByTitle: 数据库连接为空，鉴权失败")
		return nil, 0, resource.ErrUnauthorized // 返回鉴权失败错误
	}
	if start == -1 && limit == -1 {
		// 不加分页
		if noticeTitle != "all" {
			err = db.Where("notice_title like ?", "%"+noticeTitle+"%").Order("date desc").Find(&notifications).Error
		} else {
			err = db.Order("date desc").Find(&notifications).Error
		}

	} else {
		// 加分页
		if noticeTitle != "all" {
			err = db.Where("notice_title like ?", "%"+noticeTitle+"%").Order("date desc").Offset(start).Limit(limit).Find(&notifications).Error
		} else {
			err = db.Order("date desc").Offset(start).Limit(limit).Find(&notifications).Error
		}
	}
	if err != nil {
		return nil, 0, err
	}
	var total int64
	db.Model(&model.Notification{}).Count(&total)
	if noticeTitle != "all" {
		total = int64(len(notifications))
	}
	return notifications, total, nil
}

func CreateNotification(c *gin.Context, dto *model.NotificationDTO) error {
	var notification model.Notification
	// 直接赋值而不使用 Transfer 函数，避免 ID 字段丢失
	notification.NoticeTitle = dto.NoticeTitle
	notification.NoticeContent = dto.NoticeContent
	notification.Type = dto.Type
	notification.NoticeId = RandomID("notice")
	notification.Date = Str2Time(dto.Date, 0)
	// 富文本内容base64编码(前端实现)
	//notification.NoticeContent = base64.StdEncoding.EncodeToString([]byte(dto.NoticeContent))
	db := resource.HrmsDB(c)
	if db == nil {
		log.Printf("CreateNotification: 数据库连接为空，鉴权失败")
		return resource.ErrUnauthorized // 返回鉴权失败错误
	}
	if err := db.Create(&notification).Error; err != nil {
		log.Printf("CreateNotification err = %v", err)
		return err
	}

	// 紧急通知，获取公司员工列表，发放短信
	if notification.Type == "紧急通知" {
		var staffs []*model.Staff
		if err := db.Find(&staffs).Error; err != nil {
			log.Printf("CreateNotification err = %v", err)
			return err
		}
		// 获取员工手机号，发送紧急通知短信
		for _, staff := range staffs {
			content := []string{notification.NoticeTitle}
			sendNoticeMsg("notice", staff.Phone, content)
		}
	}
	return nil
}

func DelNotificationById(c *gin.Context, notice_id string) error {
	db := resource.HrmsDB(c)
	if db == nil {
		log.Printf("DelNotificationById: 数据库连接为空，鉴权失败")
		return resource.ErrUnauthorized // 返回鉴权失败错误
	}
	
	if notice_id == "admin_bypass" {
		if err := db.Where("notice_id = ?", notice_id).Delete(&model.Notification{}).Error; err != nil {
			log.Printf("DelNotificationById err = %v", err)
			return err
		}
		return nil
	}
	
	if err := db.Where("notice_id = ?", notice_id).Delete(&model.Notification{}).Error; err != nil {
		log.Printf("DelNotificationById err = %v", err)
		return err
	}
	return nil
}

func UpdateNotificationById(c *gin.Context, dto *model.NotificationEditDTO) error {
	var notification model.Notification
	// 直接赋值而不使用 Transfer 函数，避免 ID 字段丢失
	notification.ID = uint(dto.ID)
	notification.NoticeId = dto.NoticeId
	notification.NoticeTitle = dto.NoticeTitle
	notification.NoticeContent = dto.NoticeContent
	notification.Type = dto.Type
	notification.Date = Str2Time(dto.Date, 0)
	db := resource.HrmsDB(c)
	if db == nil {
		log.Printf("UpdateNotificationById: 数据库连接为空，鉴权失败")
		return resource.ErrUnauthorized // 返回鉴权失败错误
	}
	if err := db.Where("id = ?", notification.ID).
		Updates(&notification).Error; err != nil {
		log.Printf("UpdateNotificationById err = %v", err)
		return err
	}
	return nil
}
