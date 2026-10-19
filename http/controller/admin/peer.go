package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/lejianwen/rustdesk-api/v2/global"
	"github.com/lejianwen/rustdesk-api/v2/http/request/admin"
	"github.com/lejianwen/rustdesk-api/v2/http/response"
	"github.com/lejianwen/rustdesk-api/v2/service"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type Peer struct {
}

// Detail device
// @Tags Device
// @Summary Device detail
// @Description Device detail
// @Accept  json
// @Produce  json
// @Param id path int true "ID"
// @Success 200 {object} response.Response{data=model.Peer}
// @Failure 500 {object} response.Response
// @Router /admin/peer/detail/{id} [get]
// @Security token
func (ct *Peer) Detail(c *gin.Context) {
	id := c.Param("id")
	iid, _ := strconv.Atoi(id)
	u := service.AllService.PeerService.InfoByRowId(uint(iid))
	if u.RowId > 0 {
		response.Success(c, u)
		return
	}
	response.Fail(c, 101, response.TranslateMsg(c, "ItemNotFound"))
	return
}

// Create create device
// @Tags Device
// @Summary Create device
// @Description Create device
// @Accept  json
// @Produce  json
// @Param body body admin.PeerForm true "Device info"
// @Success 200 {object} response.Response{data=model.Peer}
// @Failure 500 {object} response.Response
// @Router /admin/peer/create [post]
// @Security token
func (ct *Peer) Create(c *gin.Context) {
	f := &admin.PeerForm{}
	if err := c.ShouldBindJSON(f); err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError")+err.Error())
		return
	}
	errList := global.Validator.ValidStruct(c, f)
	if len(errList) > 0 {
		response.Fail(c, 101, errList[0])
		return
	}
	p := f.ToPeer()
	err := service.AllService.PeerService.Create(p)
	if err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "OperationFailed")+err.Error())
		return
	}
	response.Success(c, nil)
}

// List list
// @Tags Device
// @Summary Device list
// @Description Device list
// @Accept  json
// @Produce  json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param time_ago query int false "Time"
// @Param id query string false "ID"
// @Param hostname query string false "Hostname"
// @Param uuids query string false "uuids, comma separated"
// @Success 200 {object} response.Response{data=model.PeerList}
// @Failure 500 {object} response.Response
// @Router /admin/peer/list [get]
// @Security token
func (ct *Peer) List(c *gin.Context) {
	query := &admin.PeerQuery{}
	if err := c.ShouldBindQuery(query); err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError")+err.Error())
		return
	}
	res := service.AllService.PeerService.List(query.Page, query.PageSize, func(tx *gorm.DB) {
		if query.TimeAgo > 0 {
			lt := time.Now().Unix() - int64(query.TimeAgo)
			tx.Where("last_online_time < ?", lt)
		}
		if query.TimeAgo < 0 {
			lt := time.Now().Unix() + int64(query.TimeAgo)
			tx.Where("last_online_time > ?", lt)
		}
		if query.Id != "" {
			tx.Where("id like ?", "%"+query.Id+"%")
		}
		if query.Hostname != "" {
			tx.Where("hostname like ?", "%"+query.Hostname+"%")
		}
		if query.Uuids != "" {
			tx.Where("uuid in (?)", query.Uuids)
		}
		if query.Username != "" {
			tx.Where("username like ?", "%"+query.Username+"%")
		}
		if query.Ip != "" {
			tx.Where("last_online_ip like ?", "%"+query.Ip+"%")
		}
		if query.Alias != "" {
			tx.Where("alias like ?", "%"+query.Alias+"%")
		}
	})
	response.Success(c, res)
}

// Update edit
// @Tags Device
// @Summary Device edit
// @Description Device edit
// @Accept  json
// @Produce  json
// @Param body body admin.PeerForm true "Device info"
// @Success 200 {object} response.Response{data=model.Peer}
// @Failure 500 {object} response.Response
// @Router /admin/peer/update [post]
// @Security token
func (ct *Peer) Update(c *gin.Context) {
	f := &admin.PeerForm{}
	if err := c.ShouldBindJSON(f); err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError")+err.Error())
		return
	}
	if f.RowId == 0 {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError"))
		return
	}
	errList := global.Validator.ValidStruct(c, f)
	if len(errList) > 0 {
		response.Fail(c, 101, errList[0])
		return
	}
	u := f.ToPeer()
	err := service.AllService.PeerService.Update(u)
	if err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "OperationFailed")+err.Error())
		return
	}
	response.Success(c, nil)
}

// Delete delete
// @Tags Device
// @Summary Device delete
// @Description Device delete
// @Accept  json
// @Produce  json
// @Param body body admin.PeerForm true "Device info"
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /admin/peer/delete [post]
// @Security token
func (ct *Peer) Delete(c *gin.Context) {
	f := &admin.PeerForm{}
	if err := c.ShouldBindJSON(f); err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError")+err.Error())
		return
	}
	id := f.RowId
	errList := global.Validator.ValidVar(c, id, "required,gt=0")
	if len(errList) > 0 {
		response.Fail(c, 101, errList[0])
		return
	}
	u := service.AllService.PeerService.InfoByRowId(f.RowId)
	if u.RowId > 0 {
		err := service.AllService.PeerService.Delete(u)
		if err == nil {
			response.Success(c, nil)
			return
		}
		response.Fail(c, 101, response.TranslateMsg(c, "OperationFailed")+err.Error())
		return
	}
	response.Fail(c, 101, response.TranslateMsg(c, "ItemNotFound"))
}

// BatchDelete batch delete
// @Tags Device
// @Summary Batch delete devices
// @Description Batch delete devices
// @Accept  json
// @Produce  json
// @Param body body admin.PeerBatchDeleteForm true "Device ids"
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /admin/peer/batchDelete [post]
// @Security token
func (ct *Peer) BatchDelete(c *gin.Context) {
	f := &admin.PeerBatchDeleteForm{}
	if err := c.ShouldBindJSON(f); err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError")+err.Error())
		return
	}
	if len(f.RowIds) == 0 {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError"))
		return
	}
	err := service.AllService.PeerService.BatchDelete(f.RowIds)
	if err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "OperationFailed")+err.Error())
		return
	}
	response.Success(c, nil)
}

func (ct *Peer) SimpleData(c *gin.Context) {
	f := &admin.SimpleDataQuery{}
	if err := c.ShouldBindJSON(f); err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError")+err.Error())
		return
	}
	if len(f.Ids) == 0 {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError"))
		return
	}
	res := service.AllService.PeerService.List(1, 99999, func(tx *gorm.DB) {
		//information that can be made public
		tx.Select("id,version")
		tx.Where("id in (?)", f.Ids)
	})
	response.Success(c, res)
}
