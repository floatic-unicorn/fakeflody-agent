package hashids

import "github.com/speps/go-hashids/v2"

const salt = "4f48c6bbb60ca18e214379a195222488"

func ToUid(id int) string {
	hd := hashids.NewData()
	hd.Salt = salt
	hd.MinLength = 8
	h, _ := hashids.NewWithData(hd)
	e, _ := h.Encode([]int{id})
	return e
}

func ToId(uid string) int {
	hd := hashids.NewData()
	hd.Salt = salt
	hd.MinLength = 8
	h, _ := hashids.NewWithData(hd)
	d, _ := h.DecodeWithError(uid)
	return d[0]
}
