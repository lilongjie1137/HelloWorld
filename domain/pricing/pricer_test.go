package pricing

import (
	"errors"
	"testing"

	"github.com/lilongjie1137/HelloWorld/common/errcode"
	"github.com/lilongjie1137/HelloWorld/common/types"
)

// 珍珠奶茶：规格(大杯=15/中杯=12) 单选必填；加料(珍珠=3/椰果=2) 多选；做法(去冰/半糖) 多选不计价。
func milkTeaSPU() SPU {
	return SPU{
		ID: 1001, Name: "珍珠奶茶", BasePrice: 0,
		Groups: []Group{
			{ID: 1, Type: types.GroupSpec, SelectType: "SINGLE", Required: true, Options: map[int64]Option{
				11: {ID: 11, Name: "大杯", Price: 1500},
				12: {ID: 12, Name: "中杯", Price: 1200},
			}},
			{ID: 2, Type: types.GroupAddon, SelectType: "MULTI", Options: map[int64]Option{
				21: {ID: 21, Name: "加珍珠", Price: 300},
				22: {ID: 22, Name: "加椰果", Price: 200},
			}},
			{ID: 3, Type: types.GroupMethod, SelectType: "MULTI", Options: map[int64]Option{
				31: {ID: 31, Name: "去冰", Price: 0},
				32: {ID: 32, Name: "半糖", Price: 0},
			}},
		},
	}
}

// 炒菜：无规格，base_price 计价；做法(免葱) 不加价。
func dishSPU() SPU {
	return SPU{
		ID: 2001, Name: "鱼香肉丝", BasePrice: 2800,
		Groups: []Group{
			{ID: 9, Type: types.GroupMethod, SelectType: "MULTI", Options: map[int64]Option{
				91: {ID: 91, Name: "免葱", Price: 0},
			}},
		},
	}
}

func TestCalculate_MilkTea_AbsoluteSpecPlusAddon(t *testing.T) {
	// 大杯(15) + 加珍珠(3) = 18，数量 2 → 36。
	line, err := Calculate(milkTeaSPU(), Selection{
		Qty: 2, SpecOptionID: 11, AddonOptionIDs: []int64{21}, MethodOptionIDs: []int64{31, 32},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if line.UnitPrice != 1800 {
		t.Errorf("unitPrice = %d, want 1800", line.UnitPrice)
	}
	if line.Amount != 3600 {
		t.Errorf("amount = %d, want 3600", line.Amount)
	}
	// 快照应含 规格1 + 加料1 + 做法2 = 4 项。
	if len(line.Modifiers) != 4 {
		t.Errorf("modifiers = %d, want 4", len(line.Modifiers))
	}
}

func TestCalculate_MilkTea_MediumNoAddon(t *testing.T) {
	line, err := Calculate(milkTeaSPU(), Selection{Qty: 1, SpecOptionID: 12})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if line.UnitPrice != 1200 || line.Amount != 1200 {
		t.Errorf("got unit=%d amount=%d, want 1200/1200", line.UnitPrice, line.Amount)
	}
}

func TestCalculate_Dish_UsesBasePrice(t *testing.T) {
	line, err := Calculate(dishSPU(), Selection{Qty: 3, MethodOptionIDs: []int64{91}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if line.UnitPrice != 2800 || line.Amount != 8400 {
		t.Errorf("got unit=%d amount=%d, want 2800/8400", line.UnitPrice, line.Amount)
	}
}

func TestCalculate_SpecRequired(t *testing.T) {
	_, err := Calculate(milkTeaSPU(), Selection{Qty: 1})
	if !errors.Is(err, errcode.ErrSpecRequired) {
		t.Fatalf("err = %v, want ErrSpecRequired", err)
	}
}

func TestCalculate_InvalidQty(t *testing.T) {
	_, err := Calculate(dishSPU(), Selection{Qty: 0})
	var be *errcode.Error
	if !errors.As(err, &be) || be.Code != errcode.ErrInvalidParam.Code {
		t.Fatalf("err = %v, want ErrInvalidParam", err)
	}
}

func TestCalculate_UnknownSpecOption(t *testing.T) {
	_, err := Calculate(milkTeaSPU(), Selection{Qty: 1, SpecOptionID: 999})
	var be *errcode.Error
	if !errors.As(err, &be) || be.Code != errcode.ErrModifierInvalid.Code {
		t.Fatalf("err = %v, want ErrModifierInvalid", err)
	}
}

func TestCalculate_SpecOnNonSpecSPU(t *testing.T) {
	_, err := Calculate(dishSPU(), Selection{Qty: 1, SpecOptionID: 11})
	var be *errcode.Error
	if !errors.As(err, &be) || be.Code != errcode.ErrModifierInvalid.Code {
		t.Fatalf("err = %v, want ErrModifierInvalid", err)
	}
}

func TestCalculate_DuplicateAddon(t *testing.T) {
	_, err := Calculate(milkTeaSPU(), Selection{Qty: 1, SpecOptionID: 11, AddonOptionIDs: []int64{21, 21}})
	if err == nil {
		t.Fatal("want error for duplicate addon, got nil")
	}
}
