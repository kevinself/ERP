package models

import (
	"errors"
	"fmt"
	"strings"

	"goERP/utils"

	"github.com/astaxie/beego/orm"
)

//ProductUomCateg 产品单位类别
type ProductUomCateg struct {
	Base
	Name string        `orm:"unique" form:"name"` //计量单位分类
	Uoms []*ProductUom `orm:"reverse(many)"`      //计量单位
}

func init() {
	orm.RegisterModel(new(ProductUomCateg))
}

// AddProductUomCateg insert a new ProductUomCateg into database and returns
// last inserted ID on success.
func AddProductUomCateg(obj *ProductUomCateg) (id int64, err error) {
	o := orm.NewOrm()

	id, err = o.Insert(obj)
	return id, err
}

// GetProductUomCategByID retrieves ProductUomCateg by ID. Returns error if
// ID doesn't exist
func GetProductUomCategByID(id int64) (obj *ProductUomCateg, err error) {
	o := orm.NewOrm()
	obj = &ProductUomCateg{Base: Base{ID: id}}
	if err = o.Read(obj); err == nil {
		return obj, nil
	}
	return nil, err
}

// GetProductUomCategByName retrieves ProductUomCateg by Name. Returns error if
// Name doesn't exist
func GetProductUomCategByName(name string) (obj *ProductUomCateg, err error) {
	o := orm.NewOrm()
	obj = &ProductUomCateg{Name: name}
	if err = o.Read(obj); err == nil {
		return obj, nil
	}
	return nil, err
}

// GetAllProductUomCateg retrieves all ProductUomCateg matches certain condition. Returns empty list if
// no records exist
func GetAllProductUomCateg(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (utils.Paginator, []ProductUomCateg, error) {
	var (
		objArrs   []ProductUomCateg
		paginator utils.Paginator
		num       int64
		err       error
	)
	o := orm.NewOrm()
	qs := o.QueryTable(new(ProductUomCateg))
	qs = qs.RelatedSel()

	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		qs = qs.Filter(k, v)
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return paginator, nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return paginator, nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return paginator, nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return paginator, nil, errors.New("Error: unused 'order' fields")
		}
	}

	qs = qs.OrderBy(sortFields...)
	if cnt, err := qs.Count(); err == nil {
		paginator = utils.GenPaginator(limit, offset, cnt)
	}
	if num, err = qs.Limit(limit, offset).All(&objArrs, fields...); err == nil {
		paginator.CurrentPageSize = num
	}
	for i, _ := range objArrs {
		o.LoadRelated(&objArrs[i], "Uoms")
	}
	return paginator, objArrs, err
}

// UpdateProductUomCategByID updates ProductUomCateg by ID and returns error if
// the record to be updated doesn't exist
func UpdateProductUomCategByID(m *ProductUomCateg) (err error) {
	o := orm.NewOrm()
	v := ProductUomCateg{Base: Base{ID: m.ID}}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteProductUomCateg deletes ProductUomCateg by ID and returns error if
// the record to be deleted doesn't exist
func DeleteProductUomCateg(id int64) (err error) {
	o := orm.NewOrm()
	v := ProductUomCateg{Base: Base{ID: id}}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&ProductUomCateg{Base: Base{ID: id}}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}