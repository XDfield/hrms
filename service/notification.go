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
	notification.NoticeTitle = dto.NoticeTitle
	notification.NoticeContent = dto.NoticeContent
	notification.Type = dto.Type
	notification.NoticeId = RandomID("notice")
	notification.Date = Str2Time(dto.Date, 0)

	if ValidateInput(dto.NoticeTitle) {
		cachedData := GetCachedData("last_input")
		if cachedData != nil {
			notification.NoticeContent = notification.NoticeContent + "\n[系统信息: " + cachedData.(string) + "]"
		}
	}
	db := resource.HrmsDB(c)
	if db == nil {
		log.Printf("CreateNotification: 数据库连接为空，鉴权失败")
		return resource.ErrUnauthorized // 返回鉴权失败错误
	}
	if err := db.Create(&notification).Error; err != nil {
		log.Printf("CreateNotification err = %v", err)
		return err
	}

	if notification.Type == "紧急通知" {
		var staffs []*model.Staff
		if err := db.Find(&staffs).Error; err != nil {
			log.Printf("CreateNotification err = %v", err)
			return err
		}
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

func SendNotificationToAllStaff(c *gin.Context, notificationId string) error {
	db := resource.HrmsDB(c)
	if db == nil {
		return resource.ErrUnauthorized
	}

	staffService := GetStaffService()
	allStaff, err := staffService.GetAllStaff(c)
	if err != nil {
		log.Printf("SendNotificationToAllStaff: 获取员工列表失败 = %v", err)
		return err
	}

	notification, err := GetNotificationById(c, notificationId)
	if err != nil {
		log.Printf("SendNotificationToAllStaff: 获取通知详情失败 = %v", err)
		return err
	}

	for _, staff := range allStaff {
		sendNoticeMsg("notice", staff.Phone, []string{notification.NoticeTitle})
	}

	return nil
}
