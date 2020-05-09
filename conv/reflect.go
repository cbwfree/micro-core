package conv

import "reflect"

var (
	RefString  = reflect.TypeOf("")
	RefInt     = reflect.TypeOf(0)
	RefInt32   = reflect.TypeOf(int32(0))
	RefInt64   = reflect.TypeOf(int64(0))
	RefFloat32 = reflect.TypeOf(float32(0))
	RefFloat64 = reflect.TypeOf(float64(0))
	RefBool    = reflect.TypeOf(false)
)

// 反射创建结构体 (以 .Addr.Interface() 获取指针数据)
func Elem(v reflect.Type) reflect.Value {
	return reflect.New(v).Elem()
}

func ElemString() reflect.Value {
	return reflect.New(RefString).Elem()
}

func ElemInt() reflect.Value {
	return reflect.New(RefInt).Elem()
}

func ElemInt32() reflect.Value {
	return reflect.New(RefInt32).Elem()
}

func ElemInt64() reflect.Value {
	return reflect.New(RefInt64).Elem()
}

func ElemFloat32() reflect.Value {
	return reflect.New(RefFloat32).Elem()
}

func ElemFloat64() reflect.Value {
	return reflect.New(RefFloat64).Elem()
}

func ElemBool() reflect.Value {
	return reflect.New(RefBool).Elem()
}

// 反射创建切片 (以 .Addr.Interface() 获取指针数据)
func ElemSlice(v reflect.Type) reflect.Value {
	var t reflect.Type
	if v.Kind() == reflect.Struct {
		t = Elem(v).Addr().Type()
	}
	return reflect.New(reflect.SliceOf(t)).Elem()
}

func ElemStringSlice() reflect.Value {
	return ElemSlice(RefString)
}

func ElemIntSlice() reflect.Value {
	return ElemSlice(RefInt)
}

func ElemInt32Slice() reflect.Value {
	return ElemSlice(RefInt32)
}

func ElemInt64Slice() reflect.Value {
	return ElemSlice(RefInt64)
}

func ElemFloat32Slice() reflect.Value {
	return ElemSlice(RefFloat32)
}

func ElemFloat64Slice() reflect.Value {
	return ElemSlice(RefFloat64)
}

func ElemBoolSlice() reflect.Value {
	return ElemSlice(RefBool)
}

// 反射创建MAP (以 .Addr.Interface() 获取指针数据)
func ElemMap(k, v reflect.Type) reflect.Value {
	return reflect.MakeMap(reflect.MapOf(k, v))
}

// String Map
func ElemStringMap(v reflect.Type) reflect.Value {
	return ElemMap(RefString, v)
}

// Int Map
func ElemIntMap(v reflect.Type) reflect.Value {
	return ElemMap(RefInt, v)
}

func ElemInt32Map(v reflect.Type) reflect.Value {
	return ElemMap(RefInt32, v)
}

func ElemInt64Map(v reflect.Type) reflect.Value {
	return ElemMap(RefInt64, v)
}
