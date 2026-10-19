package service

import (
	"github.com/lejianwen/rustdesk-api/v2/model"
	"gorm.io/gorm"
)

type PeerService struct {
}

// FindById find by id
func (ps *PeerService) FindById(id string) *model.Peer {
	p := &model.Peer{}
	DB.Where("id = ?", id).First(p)
	return p
}
func (ps *PeerService) FindByUuid(uuid string) *model.Peer {
	p := &model.Peer{}
	DB.Where("uuid = ?", uuid).First(p)
	return p
}
func (ps *PeerService) InfoByRowId(id uint) *model.Peer {
	p := &model.Peer{}
	DB.Where("row_id = ?", id).First(p)
	return p
}

// FindByUserIdAndUuid find a peer by user id and uuid
func (ps *PeerService) FindByUserIdAndUuid(uuid string, userId uint) *model.Peer {
	p := &model.Peer{}
	DB.Where("uuid = ? and user_id = ?", uuid, userId).First(p)
	return p
}

// UuidBindUserId bind the user id
func (ps *PeerService) UuidBindUserId(deviceId string, uuid string, userId uint) {
	peer := ps.FindByUuid(uuid)
	// update it if it exists
	if peer.RowId > 0 {
		peer.UserId = userId
		ps.Update(peer)
	} else {
		// create it if it does not exist
		/*if deviceId != "" {
			DB.Create(&model.Peer{
				Id:     deviceId,
				Uuid:   uuid,
				UserId: userId,
			})
		}*/
	}
}

// UuidUnbindUserId unbind the user id, used when a user logs out
func (ps *PeerService) UuidUnbindUserId(uuid string, userId uint) {
	peer := ps.FindByUserIdAndUuid(uuid, userId)
	if peer.RowId > 0 {
		DB.Model(peer).Update("user_id", 0)
	}
}

// EraseUserId clear the user id, used when a user is deleted
func (ps *PeerService) EraseUserId(userId uint) error {
	return DB.Model(&model.Peer{}).Where("user_id = ?", userId).Update("user_id", 0).Error
}

// ListByUserIds get the list by user ids
func (ps *PeerService) ListByUserIds(userIds []uint, page, pageSize uint) (res *model.PeerList) {
	res = &model.PeerList{}
	res.Page = int64(page)
	res.PageSize = int64(pageSize)
	tx := DB.Model(&model.Peer{})
	tx.Where("user_id in (?)", userIds)
	tx.Count(&res.Total)
	tx.Scopes(Paginate(page, pageSize))
	tx.Find(&res.Peers)
	return
}

func (ps *PeerService) List(page, pageSize uint, where func(tx *gorm.DB)) (res *model.PeerList) {
	res = &model.PeerList{}
	res.Page = int64(page)
	res.PageSize = int64(pageSize)
	tx := DB.Model(&model.Peer{})
	if where != nil {
		where(tx)
	}
	tx.Count(&res.Total)
	tx.Scopes(Paginate(page, pageSize))
	tx.Find(&res.Peers)
	return
}

// ListFilterByUserId filter the Peer list by user id
func (ps *PeerService) ListFilterByUserId(page, pageSize uint, where func(tx *gorm.DB), userId uint) (res *model.PeerList) {
	userWhere := func(tx *gorm.DB) {
		tx.Where("user_id = ?", userId)
		// if there is an extra filter condition, run it
		if where != nil {
			where(tx)
		}
	}
	return ps.List(page, pageSize, userWhere)
}

// Create create
func (ps *PeerService) Create(u *model.Peer) error {
	res := DB.Create(u).Error
	return res
}

// Delete delete; the tokens should be deleted too
func (ps *PeerService) Delete(u *model.Peer) error {
	uuid := u.Uuid
	err := DB.Delete(u).Error
	if err != nil {
		return err
	}
	// delete the token
	return AllService.UserService.FlushTokenByUuid(uuid)
}

// GetUuidListByIDs get the uuid list by ids
func (ps *PeerService) GetUuidListByIDs(ids []uint) ([]string, error) {
	var uuids []string
	err := DB.Model(&model.Peer{}).
		Where("row_id in (?)", ids).
		Pluck("uuid", &uuids).Error
	//filter out empty strings in uuids
	var newUuids []string
	for _, uuid := range uuids {
		if uuid != "" {
			newUuids = append(newUuids, uuid)
		}
	}
	return newUuids, err
}

// BatchDelete batch delete; the tokens should be deleted too
func (ps *PeerService) BatchDelete(ids []uint) error {
	uuids, err := ps.GetUuidListByIDs(ids)
	err = DB.Where("row_id in (?)", ids).Delete(&model.Peer{}).Error
	if err != nil {
		return err
	}
	// delete the token
	return AllService.UserService.FlushTokenByUuids(uuids)
}

// Update update
func (ps *PeerService) Update(u *model.Peer) error {
	return DB.Model(u).Updates(u).Error
}
