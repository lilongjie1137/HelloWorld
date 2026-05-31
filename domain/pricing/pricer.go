// Package pricing 实现 V1.0 核心计价逻辑（规格绝对价模型）。
//
// 规则（与《V1.0 详细设计》3.3 一致）：
//   - 选了规格(SPEC) → 该规格选项 price 为绝对基准单价，base_price 不参与；
//   - 未配规格(如炒菜) → 用 SPU.base_price 作为基准单价；
//   - 加料(ADDON) → 每个选项 price 累加；
//   - 做法(METHOD) → 不计价；
//   - unitPrice = (有规格 ? 规格.price : base_price) + Σ(ADDON.price)
//   - amount    = unitPrice * qty
//
// 金额一律以「分」(int64) 计算，避免浮点误差；展示层再转 DECIMAL(10,2)。
package pricing

import (
	"github.com/lilongjie1137/HelloWorld/common/errcode"
	"github.com/lilongjie1137/HelloWorld/common/types"
)

// Option 一个属性选项（已展开的快照视图）。
type Option struct {
	ID    int64
	Name  string
	Price int64 // 分；SPEC=绝对单价, ADDON=加价, METHOD=0
}

// Group 一个属性组及其约束。
type Group struct {
	ID         int64
	Type       types.GroupType
	SelectType string // SINGLE / MULTI
	Required   bool
	Options    map[int64]Option // optionID -> Option
}

// SPU 计价所需的商品视图。
type SPU struct {
	ID        int64
	Name      string
	BasePrice int64 // 分
	Groups    []Group
}

// Selection 一次点单对某 SPU 的属性选择。
type Selection struct {
	Qty             int
	SpecOptionID    int64   // 0 表示未选规格
	MethodOptionIDs []int64 // 做法（不计价，仅做合法性校验与快照）
	AddonOptionIDs  []int64 // 加料
}

// LineModifier 计价结果中的逐项属性快照。
type LineModifier struct {
	GroupType types.GroupType
	OptionID  int64
	Name      string
	Price     int64
}

// Line 计价结果（写入 order_item / order_item_modifier 的快照）。
type Line struct {
	SpuID     int64
	SpuName   string
	UnitPrice int64
	Qty       int
	Amount    int64
	Modifiers []LineModifier
}

// Calculate 校验选择合法性并计算单价与金额。
func Calculate(spu SPU, sel Selection) (Line, error) {
	if sel.Qty <= 0 {
		return Line{}, errcode.ErrInvalidParam.WithMessage("数量必须大于 0")
	}

	groupByType := map[types.GroupType]Group{}
	for _, g := range spu.Groups {
		groupByType[g.Type] = g
	}

	unit := spu.BasePrice
	var mods []LineModifier

	// 规格：决定绝对基准单价。
	if specGroup, ok := groupByType[types.GroupSpec]; ok {
		if sel.SpecOptionID == 0 {
			if specGroup.Required {
				return Line{}, errcode.ErrSpecRequired
			}
		} else {
			opt, ok := specGroup.Options[sel.SpecOptionID]
			if !ok {
				return Line{}, errcode.ErrModifierInvalid.WithMessage("规格选项不存在")
			}
			unit = opt.Price // 绝对价覆盖 base_price
			mods = append(mods, LineModifier{types.GroupSpec, opt.ID, opt.Name, opt.Price})
		}
	} else if sel.SpecOptionID != 0 {
		return Line{}, errcode.ErrModifierInvalid.WithMessage("该商品无规格")
	}

	// 做法：不计价，仅校验归属并写快照。
	if err := appendNonSpec(groupByType, types.GroupMethod, sel.MethodOptionIDs, &unit, &mods, false); err != nil {
		return Line{}, err
	}
	// 加料：累加。
	if err := appendNonSpec(groupByType, types.GroupAddon, sel.AddonOptionIDs, &unit, &mods, true); err != nil {
		return Line{}, err
	}

	return Line{
		SpuID:     spu.ID,
		SpuName:   spu.Name,
		UnitPrice: unit,
		Qty:       sel.Qty,
		Amount:    unit * int64(sel.Qty),
		Modifiers: mods,
	}, nil
}

// appendNonSpec 处理 METHOD / ADDON 选项：校验存在、单选约束，按需累加价格。
func appendNonSpec(
	groups map[types.GroupType]Group,
	gt types.GroupType,
	optionIDs []int64,
	unit *int64,
	mods *[]LineModifier,
	addPrice bool,
) error {
	if len(optionIDs) == 0 {
		return nil
	}
	g, ok := groups[gt]
	if !ok {
		return errcode.ErrModifierInvalid.WithMessage("该商品无对应属性组")
	}
	if g.SelectType == "SINGLE" && len(optionIDs) > 1 {
		return errcode.ErrModifierInvalid.WithMessage("该属性为单选")
	}
	seen := map[int64]struct{}{}
	for _, id := range optionIDs {
		if _, dup := seen[id]; dup {
			return errcode.ErrModifierInvalid.WithMessage("属性选项重复")
		}
		seen[id] = struct{}{}
		opt, ok := g.Options[id]
		if !ok {
			return errcode.ErrModifierInvalid.WithMessage("属性选项不存在")
		}
		if addPrice {
			*unit += opt.Price
		}
		*mods = append(*mods, LineModifier{gt, opt.ID, opt.Name, opt.Price})
	}
	return nil
}
