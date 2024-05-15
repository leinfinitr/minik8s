package reflectprop

import (
	"minik8s/tools/log"
	"os"
	"reflect"
)

func GetObjNamespace(obj interface{}) string {
	// 通过反射获取对象的namespace
	// 如果对象没有ObjectMeta或Namespace属性，则返回"default"
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		log.ErrorLog("Expected a struct type")
		os.Exit(1)
	}

	objectMetaField := v.FieldByName("ObjectMeta")
	if !objectMetaField.IsValid() {
		log.ErrorLog("No ObjectMeta field found")
		os.Exit(1)
	}

	namespaceField := objectMetaField.FieldByName("Namespace")
	if !namespaceField.IsValid() || namespaceField.String() == "" {
		// 如果没有找到命名空间或者命名空间为空，返回"default"
		return "default"
	}

	return namespaceField.String()
}

func GetObjName(obj interface{}) string {
	// 通过反射获取对象的name
	// 如果对象没有ObjectMeta或Name属性，则返回""
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		log.ErrorLog("Expected a struct type")
		os.Exit(1)
	}

	objectMetaField := v.FieldByName("ObjectMeta")
	if !objectMetaField.IsValid() {
		log.ErrorLog("No ObjectMeta field found")
		os.Exit(1)
	}

	nameField := objectMetaField.FieldByName("Name")
	if !nameField.IsValid() || nameField.String() == "" {
		// 如果没有找到name或者name为空，返回""
		return ""
	}

	return nameField.String()
}
