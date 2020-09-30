package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"u3.com/u3query/hack"
	"u3.com/u3query/models"
)

const SplitLength = hack.SplitLength

// Operations about Unit
type UnitController struct {
	beego.Controller
}

// @Title CreateUser
// @Description create Unit
// @Param	body		body 	models.Unit	true		"body for unit content"
// @Success 200 {int} models.Unit.Id
// @Failure 403 body is empty
// @router / [post]
func (u *UnitController) Post() {
	var unit models.Unit
	json.Unmarshal(u.Ctx.Input.RequestBody, &unit)

	id, err := models.InsertUnit(&unit)
	if err != nil {
		u.Data["json"] = err.Error()
	} else {
		u.Data["json"] = map[string]int{"id": id}
	}
	u.ServeJSON()
}

// @Title Get100
// @Description get top100 Unit
// @Success 200 {object} models.Unit
// @router / [get]
func (u *UnitController) GetAll() {
	rlt := make(map[int]*models.Unit)
	bt, err := models.CacheBt.GetCacheBt("0-100000")
	if err != nil {
		u.Data["json"] = err.Error()
	} else {
		if err == nil {
			for i := 0; i < 100; i++ {
				unit, ok := bt.Search(i)
				if ok {
					one := unit.(*models.Unit)
					rlt[i] = one
				}
			}
		}
		u.Data["json"] = rlt
	}
	u.ServeJSON()
}

// @Title Get
// @Description get unit by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Unit
// @Failure 403 :id is empty
// @router /:id [get]
func (u *UnitController) Get() {
	id, err := u.GetInt(":id")
	if err == nil {
		unit, err := models.GetUnit(id)
		if err != nil {
			u.Data["json"] = err.Error()
		} else {
			u.Data["json"] = unit
		}
	} else {
		u.Data["json"] = err.Error()
	}
	u.ServeJSON()
}

