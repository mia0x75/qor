package admin

import (
	"fmt"
	"github.com/qor/qor"

	"reflect"
)

func (admin *Admin) Dashboard(app *qor.Context) {
	admin.Render("dashboard", app)
}

func (admin *Admin) Index(context *qor.Context) {
	resource := admin.Resources[context.ResourceName]
	sliceType := reflect.SliceOf(reflect.Indirect(reflect.ValueOf(resource.Model)).Type())
	slice := reflect.MakeSlice(sliceType, 0, 0)
	slicePtr := reflect.New(sliceType)
	slicePtr.Elem().Set(slice)
	admin.DB.Debug().Find(slicePtr.Interface())

	fmt.Println(slicePtr.Interface())
	data := map[string]interface{}{"Admin": admin, "Context": context, "Resource": resource, "Result": slicePtr.Interface()}
	if err := admin.Render("index", context).Execute(context.Writer, data); err != nil {
		fmt.Println(err)
	}
}

func (admin *Admin) Show(context *qor.Context) {
	resource := admin.Resources[context.ResourceName]
	res := reflect.New(reflect.Indirect(reflect.ValueOf(resource.Model)).Type())
	admin.DB.First(res.Interface(), context.ResourceID)

	admin.Render("show", context)
}

func (admin *Admin) New(context *qor.Context) {
	admin.Render("new", context)
}

func (admin *Admin) Create(context *qor.Context) {
}

func (admin *Admin) Update(context *qor.Context) {
}

func (admin *Admin) Delete(context *qor.Context) {
}