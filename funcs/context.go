package funcs

import (
	"context"
	"fmt"
	"time"
	"unsafe"
)

func WriteKeys() {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "f1", time.Now())
	ctx = context.WithValue(ctx, "req", "sdfsdfwrrew3")
	ctx = context.WithValue(ctx, "tokenCtxKey", "tokenString")
	kws := getKeyValues(ctx)
	fmt.Println(kws)

	nc := RewriteKVToNewContext(ctx)
	fmt.Println("-------------")
	fmt.Println(getKeyValues(nc))
}

func getKeyValues(ctx context.Context) map[interface{}]interface{} {
	m := make(map[interface{}]interface{})
	getKeyValue(ctx, m)
	return m
}

type iface struct {
	itab, data uintptr
}

type emptyCtx int

type valueCtx struct {
	context.Context
	key, val interface{}
}

func GetKeyValues(ctx context.Context) map[interface{}]interface{} {
	m := make(map[interface{}]interface{})
	getKeyValue(ctx, m)
	return m
}

func getKeyValue(ctx context.Context, m map[interface{}]interface{}) {
	ictx := *(*iface)(unsafe.Pointer(&ctx))
	if ictx.data == 0 || int(*(*emptyCtx)(unsafe.Pointer(ictx.data))) == 0 {
		return
	}

	valCtx := (*valueCtx)(unsafe.Pointer(ictx.data))
	if valCtx != nil && valCtx.key != nil {
		m[valCtx.key] = valCtx.val
	}
	getKeyValue(valCtx.Context, m)
}

func RewriteKVToNewContext(ctx context.Context) context.Context {
	kvm := GetKeyValues(ctx)
	newCtx := context.Background()
	for k, v := range kvm {
		newCtx = context.WithValue(newCtx, k, v)
	}

	return newCtx
}
