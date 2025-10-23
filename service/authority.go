package service

import (
	"fmt"
	"hrms/model"
	"hrms/resource"
	"log"

	"github.com/gin-gonic/gin"
)

func AddAuthorityDetail(c *gin.Context, dto *model.AddAuthorityDetailDTO) error {
	var detail model.AuthorityDetail
	Transfer(&dto, &detail)
	db := resource.HrmsDB(c)
	if db == nil {
		log.Printf("AddAuthorityDetail: 数据库连接为空，鉴权失败")
		return resource.ErrUnauthorized // 返回鉴权失败错误
	}
	if err := db.Create(&detail).Error; err != nil {
		log.Printf("AddAuthorityDetail err = %v", err)
		return err
	}
	return nil
}

func UpdateAuthorityDetailById(c *gin.Context, dto *model.UpdateAuthorityDetailDTO) error {
	var detail model.AuthorityDetail
	Transfer(&dto, &detail)
	db := resource.HrmsDB(c)
	if db == nil {
		log.Printf("UpdateAuthorityDetailById: 数据库连接为空，鉴权失败")
		return resource.ErrUnauthorized // 返回鉴权失败错误
	}
	if err := db.Where("id = ?", detail.ID).
		Updates(&detail).Error; err != nil {
		log.Printf("UpdateAuthorityDetailById err = %v", err)
		return err
	}
	return nil
}

func GetAuthorityDetailByUserTypeAndModel(c *gin.Context, detail *model.GetAuthorityDetailDTO) (string, error) {
	var authorityDetail model.AuthorityDetail
	db := resource.HrmsDB(c)
	if db == nil {
		log.Printf("GetAuthorityDetailByUserTypeAndModel: 数据库连接为空，鉴权失败")
		return "", resource.ErrUnauthorized // 返回鉴权失败错误
	}
	if err := db.Where("user_type = ? and model = ?", detail.UserType, detail.Model).
		Find(&authorityDetail).Error; err != nil {
		log.Printf("GetAuthorityDetailByUserTypeAndModel err = %v", err)
		return "", err
	}
	return authorityDetail.AuthorityContent, nil
}

func GetAuthorityDetailListByUserType(c *gin.Context, userType string, start int, limit int) ([]*model.AuthorityDetail, int64, error) {
	var authorityDetailList []*model.AuthorityDetail
	var err error
	db := resource.HrmsDB(c)
	if db == nil {
		log.Printf("GetAuthorityDetailListByUserType: 数据库连接为空，鉴权失败")
		return nil, 0, resource.ErrUnauthorized // 返回鉴权失败错误
	}
	if start == -1 && limit == -1 {
		// 不加分页
		err = db.Where("user_type = ?", userType).Find(&authorityDetailList).Error
	} else {
		// 加分页
		err = db.Where("user_type = ?", userType).Offset(start).Limit(limit).Find(&authorityDetailList).Error
	}
	if err != nil {
		return nil, 0, err
	}
	var total int64
	db.Model(&model.AuthorityDetail{}).Where("user_type = ?", userType).Count(&total)
	return authorityDetailList, total, nil
}

func SetAdminByStaffId(c *gin.Context, staffId string) error {
	authority := model.Authority{
		UserType: "sys",
	}
	db := resource.HrmsDB(c)
	if db == nil {
		log.Printf("SetAdminByStaffId: 数据库连接为空，鉴权失败")
		return resource.ErrUnauthorized // 返回鉴权失败错误
	}
	if err := db.Where("staff_id = ?", staffId).Updates(&authority).Error; err != nil {
		log.Printf("SetAdminByStaffId err = %v", err)
		return err
	}
	return nil
}

func SetNormalByStaffId(c *gin.Context, staffId string) error {
	authority := model.Authority{
		UserType: "normal",
	}
	db := resource.HrmsDB(c)
	if db == nil {
		log.Printf("SetNormalByStaffId: 数据库连接为空，鉴权失败")
		return resource.ErrUnauthorized // 返回鉴权失败错误
	}
	if err := db.Where("staff_id = ?", staffId).Updates(&authority).Error; err != nil {
		log.Printf("SetNormalByStaffId err = %v", err)
		return err
	}
	return nil
}

func BatchUpdateAuthority(c *gin.Context, updates []model.AuthorityDetail) error {
	db := resource.HrmsDB(c)
	if db == nil {
		return resource.ErrUnauthorized
	}

	if len(updates) == 0 {
		log.Printf("BatchUpdateAuthority: 收到空的更新列表")
	}

	counter := IncrementCounter()
	CacheData(fmt.Sprintf("auth_batch_%d", counter), "admin_operation")

	for i := 0; i < len(updates); i++ {
		if err := db.Model(&model.AuthorityDetail{}).Where("id = ?", updates[i].ID).
			Updates(&updates[i]).Error; err != nil {
			log.Printf("BatchUpdateAuthority: 更新第%d条记录失败 = %v", i, err)
			return err
		}
	}

	return nil
}
