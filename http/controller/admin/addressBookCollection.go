package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/lejianwen/rustdesk-api/v2/global"
	"github.com/lejianwen/rustdesk-api/v2/http/request/admin"
	"github.com/lejianwen/rustdesk-api/v2/http/response"
	"github.com/lejianwen/rustdesk-api/v2/model"
	"github.com/lejianwen/rustdesk-api/v2/service"
	"gorm.io/gorm"
	"strconv"
)

type AddressBookCollection struct {
}

// Detail address book name
// @Tags Address book name
// @Summary Address book name detail
// @Description Address book name detail
// @Accept  json
// @Produce  json
// @Param id path int true "ID"
// @Success 200 {object} response.Response{data=model.AddressBookCollection}
// @Failure 500 {object} response.Response
// @Router /admin/address_book_collection/detail/{id} [get]
// @Security token
func (abc *AddressBookCollection) Detail(c *gin.Context) {
	id := c.Param("id")
	iid, _ := strconv.Atoi(id)
	t := service.AllService.AddressBookService.CollectionInfoById(uint(iid))
	if t.Id > 0 {
		response.Success(c, t)
		return
	}
	response.Fail(c, 101, response.TranslateMsg(c, "ItemNotFound"))
	return
}

// Create create address book name
// @Tags Address book name
// @Summary Create address book name
// @Description Create address book name
// @Accept  json
// @Produce  json
// @Param body body model.AddressBookCollection true "Address book name info"
// @Success 200 {object} response.Response{data=model.AddressBookCollection}
// @Failure 500 {object} response.Response
// @Router /admin/address_book_collection/create [post]
// @Security token
func (abc *AddressBookCollection) Create(c *gin.Context) {
	f := &model.AddressBookCollection{}
	if err := c.ShouldBindJSON(f); err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError")+err.Error())
		return
	}
	errList := global.Validator.ValidStruct(c, f)
	if len(errList) > 0 {
		response.Fail(c, 101, errList[0])
		return
	}
	if f.UserId == 0 {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError"))
		return
	}
	t := f
	err := service.AllService.AddressBookService.CreateCollection(t)
	if err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "OperationFailed")+err.Error())
		return
	}
	response.Success(c, nil)
}

// List list
// @Tags Address book name
// @Summary Address book name list
// @Description Address book name list
// @Accept  json
// @Produce  json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param is_my query int false "Whether it is mine"
// @Param user_id query int false "User id"
// @Success 200 {object} response.Response{data=model.AddressBookCollectionList}
// @Failure 500 {object} response.Response
// @Router /admin/address_book_collection/list [get]
// @Security token
func (abc *AddressBookCollection) List(c *gin.Context) {
	query := &admin.AddressBookCollectionQuery{}
	if err := c.ShouldBindQuery(query); err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError")+err.Error())
		return
	}
	res := service.AllService.AddressBookService.ListCollection(query.Page, query.PageSize, func(tx *gorm.DB) {
		if query.UserId > 0 {
			tx.Where("user_id = ?", query.UserId)
		}
	})
	response.Success(c, res)
}

// Update edit
// @Tags Address book name
// @Summary Address book name edit
// @Description Address book name edit
// @Accept  json
// @Produce  json
// @Param body body model.AddressBookCollection true "Address book name info"
// @Success 200 {object} response.Response{data=model.AddressBookCollection}
// @Failure 500 {object} response.Response
// @Router /admin/address_book_collection/update [post]
// @Security token
func (abc *AddressBookCollection) Update(c *gin.Context) {
	f := &model.AddressBookCollection{}
	if err := c.ShouldBindJSON(f); err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError")+err.Error())
		return
	}
	errList := global.Validator.ValidStruct(c, f)
	if len(errList) > 0 {
		response.Fail(c, 101, errList[0])
		return
	}
	if f.Id == 0 {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError"))
		return
	}
	t := f //f.ToAddressBookCollection()
	err := service.AllService.AddressBookService.UpdateCollection(t)
	if err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "OperationFailed")+err.Error())
		return
	}
	response.Success(c, nil)
}

// Delete delete
// @Tags Address book name
// @Summary Address book name delete
// @Description Address book name delete
// @Accept  json
// @Produce  json
// @Param body body model.AddressBookCollection true "Address book name info"
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /admin/address_book_collection/delete [post]
// @Security token
func (abc *AddressBookCollection) Delete(c *gin.Context) {
	f := &model.AddressBookCollection{}
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
	ex := service.AllService.AddressBookService.CollectionInfoById(f.Id)
	if ex.Id == 0 {
		response.Fail(c, 101, response.TranslateMsg(c, "ItemNotFound"))
		return
	}
	err := service.AllService.AddressBookService.DeleteCollection(ex)
	if err == nil {
		response.Success(c, nil)
		return
	}
	response.Fail(c, 101, response.TranslateMsg(c, "OperationFailed")+err.Error())
}
