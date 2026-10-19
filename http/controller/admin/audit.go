package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/lejianwen/rustdesk-api/v2/global"
	"github.com/lejianwen/rustdesk-api/v2/http/request/admin"
	"github.com/lejianwen/rustdesk-api/v2/http/response"
	"github.com/lejianwen/rustdesk-api/v2/model"
	"github.com/lejianwen/rustdesk-api/v2/service"
	"gorm.io/gorm"
)

type Audit struct {
}

// ConnList list
// @Tags Connection log
// @Summary Connection log list
// @Description Connection log list
// @Accept  json
// @Produce  json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param peer_id query int false "Target device"
// @Param from_peer query int false "Source device"
// @Success 200 {object} response.Response{data=model.AuditConnList}
// @Failure 500 {object} response.Response
// @Router /admin/audit_conn/list [get]
// @Security token
func (a *Audit) ConnList(c *gin.Context) {
	query := &admin.AuditQuery{}
	if err := c.ShouldBindQuery(query); err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError")+err.Error())
		return
	}
	res := service.AllService.AuditService.AuditConnList(query.Page, query.PageSize, func(tx *gorm.DB) {
		if query.PeerId != "" {
			tx.Where("peer_id like ?", "%"+query.PeerId+"%")
		}
		if query.FromPeer != "" {
			tx.Where("from_peer like ?", "%"+query.FromPeer+"%")
		}
		tx.Order("id desc")
	})
	response.Success(c, res)
}

// ConnDelete delete
// @Tags Connection log
// @Summary Connection log delete
// @Description Connection log delete
// @Accept  json
// @Produce  json
// @Param body body model.AuditConn true "Connection log info"
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /admin/audit_conn/delete [post]
// @Security token
func (a *Audit) ConnDelete(c *gin.Context) {
	f := &model.AuditConn{}
	if err := c.ShouldBindJSON(f); err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError")+err.Error())
		return
	}
	id := f.Id
	errList := global.Validator.ValidVar(c, id, "required,gt=0")
	if len(errList) > 0 {
		response.Fail(c, 101, errList[0])
		return
	}
	l := service.AllService.AuditService.ConnInfoById(f.Id)
	if l.Id > 0 {
		err := service.AllService.AuditService.DeleteAuditConn(l)
		if err == nil {
			response.Success(c, nil)
			return
		}
		response.Fail(c, 101, err.Error())
		return
	}
	response.Fail(c, 101, response.TranslateMsg(c, "ItemNotFound"))
}

// BatchConnDelete delete
// @Tags Connection log
// @Summary Connection log batch delete
// @Description Connection log batch delete
// @Accept  json
// @Produce  json
// @Param body body admin.AuditConnLogIds true "Connection logs"
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /admin/audit_conn/batchDelete [post]
// @Security token
func (a *Audit) BatchConnDelete(c *gin.Context) {
	f := &admin.AuditConnLogIds{}
	if err := c.ShouldBindJSON(f); err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError")+err.Error())
		return
	}
	if len(f.Ids) == 0 {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError"))
		return
	}

	err := service.AllService.AuditService.BatchDeleteAuditConn(f.Ids)
	if err == nil {
		response.Success(c, nil)
		return
	}
	response.Fail(c, 101, err.Error())
	return
}

// FileList list
// @Tags File log
// @Summary File log list
// @Description File log list
// @Accept  json
// @Produce  json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param peer_id query int false "Target device"
// @Param from_peer query int false "Source device"
// @Success 200 {object} response.Response{data=model.AuditFileList}
// @Failure 500 {object} response.Response
// @Router /admin/audit_file/list [get]
// @Security token
func (a *Audit) FileList(c *gin.Context) {
	query := &admin.AuditQuery{}
	if err := c.ShouldBindQuery(query); err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError")+err.Error())
		return
	}
	res := service.AllService.AuditService.AuditFileList(query.Page, query.PageSize, func(tx *gorm.DB) {
		if query.PeerId != "" {
			tx.Where("peer_id like ?", "%"+query.PeerId+"%")
		}
		if query.FromPeer != "" {
			tx.Where("from_peer like ?", "%"+query.FromPeer+"%")
		}
		tx.Order("id desc")
	})
	response.Success(c, res)
}

// FileDelete delete
// @Tags File log
// @Summary File log delete
// @Description File log delete
// @Accept  json
// @Produce  json
// @Param body body model.AuditFile true "File log info"
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /admin/audit_file/delete [post]
// @Security token
func (a *Audit) FileDelete(c *gin.Context) {
	f := &model.AuditFile{}
	if err := c.ShouldBindJSON(f); err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError")+err.Error())
		return
	}
	id := f.Id
	errList := global.Validator.ValidVar(c, id, "required,gt=0")
	if len(errList) > 0 {
		response.Fail(c, 101, errList[0])
		return
	}
	l := service.AllService.AuditService.FileInfoById(f.Id)
	if l.Id > 0 {
		err := service.AllService.AuditService.DeleteAuditFile(l)
		if err == nil {
			response.Success(c, nil)
			return
		}
		response.Fail(c, 101, err.Error())
		return
	}
	response.Fail(c, 101, response.TranslateMsg(c, "ItemNotFound"))
}

// BatchFileDelete delete
// @Tags File log
// @Summary File log batch delete
// @Description File log batch delete
// @Accept  json
// @Produce  json
// @Param body body admin.AuditFileLogIds true "File logs"
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /admin/audit_file/batchDelete [post]
// @Security token
func (a *Audit) BatchFileDelete(c *gin.Context) {
	f := &admin.AuditFileLogIds{}
	if err := c.ShouldBindJSON(f); err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError")+err.Error())
		return
	}
	if len(f.Ids) == 0 {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError"))
		return
	}

	err := service.AllService.AuditService.BatchDeleteAuditFile(f.Ids)
	if err == nil {
		response.Success(c, nil)
		return
	}
	response.Fail(c, 101, err.Error())
	return
}
