package rbac

import (
	"encoding/json"
	"github.com/astaxie/beego/orm"
	c "github.com/osgochina/admin/controllers"
	m "github.com/osgochina/admin/models/rbacmodels"
)

type NodeController struct {
	c.CommonController
}

func (this *NodeController) Rsp(status bool, str string) {
	this.Data["json"] = &map[string]interface{}{"status": status, "info": str}
	this.ServeJson()
}

func (this *NodeController) Index() {
	if this.IsAjax() {
		page, _ := this.GetInt("page")
		page_size, _ := this.GetInt("rows")
		sort := this.GetString("sort")
		order := this.GetString("order")
		if len(order) > 0 {
			if order == "desc" {
				sort = "-" + sort
			}
		} else {
			sort = "Id"
		}
		nodes, count := m.GetNodelist(page, page_size, sort)
		for i := 0; i < len(nodes); i++ {
			if nodes[i]["Pid"] != 0 {
				nodes[i]["_parentId"] = nodes[i]["Pid"]
			} else {
				nodes[i]["state"] = "closed"
			}
		}
		if len(nodes) < 1 {
			nodes = []orm.Params{}
		}
		this.Data["json"] = &map[string]interface{}{"total": count, "rows": &nodes}
		this.ServeJson()
		return
	} else {
		grouplist := m.GroupList()
		b, _ := json.Marshal(grouplist)
		this.Data["grouplist"] = string(b)
		this.TplNames = "easyui/rbac/node.tpl"
	}

}
func (this *NodeController) AddAndEdit() {
	n := m.Node{}
	if err := this.ParseForm(&n); err != nil {
		//handle error
		this.Rsp(false, err.Error())
		return
	}
	var id int64
	var err error
	Nid, _ := this.GetInt("Id")
	if Nid > 0 {
		id, err = m.UpdateNode(&n)
	} else {
		group_id, _ := this.GetInt("Group_id")
		group := new(m.Group)
		group.Id = group_id
		n.Group = group
		if n.Pid != 0 {
			n1, _ := m.ReadNode(n.Pid)
			n.Level = n1.Level + 1
		} else {
			n.Level = 1
		}
		id, err = m.AddNode(&n)
	}
	if err == nil && id > 0 {
		this.Rsp(true, "Success")
		return
	} else {
		this.Rsp(false, err.Error())
		return
	}

}

func (this *NodeController) DelNode() {
	Id, _ := this.GetInt("Id")
	status, err := m.DelNodeById(Id)
	if err == nil && status > 0 {
		this.Rsp(true, "Success")
		return
	} else {
		this.Rsp(false, err.Error())
		return
	}
}
