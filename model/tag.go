package model

type Tag struct {
	IdModel
	Name         string                 `json:"name" gorm:"default:'';not null;"`
	UserId       uint                   `json:"user_id" gorm:"default:0;not null;index"`
	Color        uint                   `json:"color" gorm:"default:0;not null;"` //color is a flutter color value, from 0x00000000 to 0xFFFFFFFF; the first two digits are the alpha, the last 6 the color; can be converted to rgba
	CollectionId uint                   `json:"collection_id" gorm:"default:0;not null;index"`
	Collection   *AddressBookCollection `json:"collection,omitempty"`
	TimeModel
}

type TagList struct {
	Tags []*Tag `json:"list"`
	Pagination
}
