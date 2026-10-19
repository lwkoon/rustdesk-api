package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	request "github.com/lejianwen/rustdesk-api/v2/http/request/api"
	"github.com/lejianwen/rustdesk-api/v2/http/response"
	"github.com/lejianwen/rustdesk-api/v2/model"
	"github.com/lejianwen/rustdesk-api/v2/service"
	"time"
)

type Audit struct {
}

// AuditConn
// @Tags Audit
// @Summary Audit connection
// @Description Audit connection
// @Accept  json
// @Produce  json
// @Param body body request.AuditConnForm true "Audit connection"
// @Success 200 {string} string ""
// @Failure 500 {object} response.Response
// @Router /audit/conn [post]
func (a *Audit) AuditConn(c *gin.Context) {
	af := &request.AuditConnForm{}
	err := c.ShouldBindBodyWith(af, binding.JSON)
	if err != nil {
		response.Error(c, response.TranslateMsg(c, "ParamsError")+err.Error())
		return
	}
	/*ttt := &gin.H{}
	c.ShouldBindBodyWith(ttt, binding.JSON)
	fmt.Println(ttt)*/
	// Same unauthenticated-endpoint problem as CVE-2026-111324: only accept audit records from a known device
	// that presents its own uuid, so nobody can forge or close connection logs for other devices.
	if !auditFromKnownPeer(af.Id, af.Uuid) {
		response.Success(c, "")
		return
	}
	ac := af.ToAuditConn()
	if af.Action == model.AuditActionNew {
		service.AllService.AuditService.CreateAuditConn(ac)
	} else if af.Action == model.AuditActionClose {
		ex := service.AllService.AuditService.InfoByPeerIdAndConnId(af.Id, af.ConnId)
		if ex.Id != 0 {
			ex.CloseTime = time.Now().Unix()
			service.AllService.AuditService.UpdateAuditConn(ex)
		}
	} else if af.Action == "" {
		ex := service.AllService.AuditService.InfoByPeerIdAndConnId(af.Id, af.ConnId)
		if ex.Id != 0 {
			up := &model.AuditConn{
				IdModel:   model.IdModel{Id: ex.Id},
				FromPeer:  ac.FromPeer,
				FromName:  ac.FromName,
				SessionId: ac.SessionId,
				Type:      ac.Type,
			}
			service.AllService.AuditService.UpdateAuditConn(up)
		}
	}
	response.Success(c, "")
}

// AuditFile
// @Tags Audit
// @Summary Audit file
// @Description Audit file
// @Accept  json
// @Produce  json
// @Param body body request.AuditFileForm true "Audit file"
// @Success 200 {string} string ""
// @Failure 500 {object} response.Response
// @Router /audit/file [post]
func (a *Audit) AuditFile(c *gin.Context) {
	aff := &request.AuditFileForm{}
	err := c.ShouldBindBodyWith(aff, binding.JSON)
	if err != nil {
		response.Error(c, response.TranslateMsg(c, "ParamsError")+err.Error())
		return
	}
	//ttt := &gin.H{}
	//c.ShouldBindBodyWith(ttt, binding.JSON)
	//fmt.Println(ttt)
	if !auditFromKnownPeer(aff.Id, aff.Uuid) {
		response.Success(c, "")
		return
	}
	af := aff.ToAuditFile()
	service.AllService.AuditService.CreateAuditFile(af)
	response.Success(c, "")
}

// auditFromKnownPeer reports whether id is a registered device whose stored uuid equals uuid.
func auditFromKnownPeer(id, uuid string) bool {
	if id == "" || uuid == "" {
		return false
	}
	pe := service.AllService.PeerService.FindById(id)
	return pe != nil && pe.RowId != 0 && pe.Uuid == uuid
}
