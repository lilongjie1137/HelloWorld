# 餐饮系统 V1.0 · 图形化设计（ER 图 + 关键流程时序图）

> 使用 Mermaid 绘制，GitHub / 支持 Mermaid 的 Markdown 查看器可直接渲染。
> 配套《餐饮系统_V1.0详细设计》。已采用决策：全国统一价、规格选项绝对价、出品部门按类型、品牌统一商户号、每店边缘网关、Go。

---

## 1. 实体关系图（ER）

```mermaid
erDiagram
    TENANT ||--o{ STORE : owns
    TENANT ||--o{ EMPLOYEE : has
    TENANT ||--o{ CATEGORY : has
    TENANT ||--o{ SPU : has
    TENANT ||--o{ MODIFIER_GROUP : has

    STORE ||--o{ PRODUCTION_DEPT : has
    STORE ||--o{ TABLE_AREA : has
    STORE ||--o{ PRINTER : has
    STORE ||--o{ DEVICE : has
    STORE ||--o{ SHIFT : has

    TABLE_AREA ||--o{ DINING_TABLE : contains
    DINING_TABLE ||--o{ ORDER : "一桌多单"

    CATEGORY ||--o{ SPU : classifies
    SPU ||--o{ SPU_MODIFIER_GROUP : binds
    MODIFIER_GROUP ||--o{ SPU_MODIFIER_GROUP : binds
    MODIFIER_GROUP ||--o{ MODIFIER_OPTION : has
    SPU ||--o{ SPU_STORE : "适用门店"
    STORE ||--o{ SPU_STORE : "适用门店"

    EMPLOYEE ||--o{ EMPLOYEE_ROLE : has
    ROLE ||--o{ EMPLOYEE_ROLE : has
    ROLE ||--o{ ROLE_PERMISSION : has
    PERMISSION ||--o{ ROLE_PERMISSION : has
    EMPLOYEE ||--o{ EMPLOYEE_STORE : "数据范围"
    STORE ||--o{ EMPLOYEE_STORE : "数据范围"

    ORDER ||--o{ ORDER_ITEM : contains
    ORDER_ITEM ||--o{ ORDER_ITEM_MODIFIER : "规格/做法/加料快照"
    ORDER ||--o{ BILL : settles
    BILL ||--o{ PAYMENT : "收款"
    BILL ||--o{ DISCOUNT : "折扣"
    ORDER ||--o{ PRINT_JOB : "拆单打印"
    PRINTER ||--o{ PRINT_JOB : executes
    PRINTER ||--o{ PRINTER_DEPT : "绑定部门类型"

    TENANT ||--o{ MEMBER : "P1预留"
    MEMBER ||--o{ MEMBER_CARD : has
    MEMBER ||--|| STORED_VALUE_ACCOUNT : has
    STORED_VALUE_ACCOUNT ||--o{ STORED_VALUE_TXN : logs

    SPU {
        bigint id
        bigint category_id
        string name
        decimal base_price "未配规格时的基准单价"
        string production_dept_type "KITCHEN/BAR/COLD"
        tinyint is_combo
        string status
    }
    MODIFIER_GROUP {
        bigint id
        string group_type "SPEC/METHOD/ADDON"
        string select_type "SINGLE/MULTI"
        tinyint required
    }
    MODIFIER_OPTION {
        bigint id
        bigint group_id
        string name
        decimal price "SPEC=绝对价/ADDON=加价/METHOD=0"
    }
    ORDER {
        bigint id
        bigint table_id
        string order_no
        string source "POS/MINI_APP/TAKEOUT"
        string status
        decimal payable_amount
    }
    ORDER_ITEM {
        bigint id
        bigint spu_id
        string spu_name_snap
        string spec_name_snap
        decimal unit_price
        int qty
        string production_dept_type_snap
        string item_status
    }
    PRINT_JOB {
        bigint id
        bigint printer_id
        string production_dept_type
        string job_type "KITCHEN/BAR/LABEL/CHECKOUT/REFUND"
        string status
        string job_key "幂等键防重打"
    }
```

---

## 2. 商品属性与计价模型（类图视角）

```mermaid
classDiagram
    class SPU {
      +name
      +base_price
      +production_dept_type
    }
    class ModifierGroup {
      +group_type : SPEC|METHOD|ADDON
      +select_type : SINGLE|MULTI
      +required
    }
    class ModifierOption {
      +name
      +price
    }
    SPU "1" --> "*" ModifierGroup : 绑定
    ModifierGroup "1" --> "*" ModifierOption

    note for ModifierOption "SPEC: price=绝对单价(大杯15/中杯12)\nADDON: price=加价(加珍珠+3)\nMETHOD: price=0(去冰/半糖)"
    note for SPU "unit_price = (有SPEC ? SPEC.price : base_price) + Σ ADDON.price"
```

---

## 3. 时序图：开台 → 混合点单 → 分单打印

```mermaid
sequenceDiagram
    participant POS as Android收银端
    participant GW as 门店边缘网关
    participant TX as 交易服务(云)
    participant MQ as Kafka
    participant PR as 打印路由服务
    participant K as 后厨打印机
    participant B as 吧台打印机

    POS->>TX: 开台(tableId, peopleCount)
    TX-->>POS: 桌台OPENED + orderId
    POS->>TX: 下单[鱼香肉丝(KITCHEN), 珍珠奶茶大杯去冰加珍珠(BAR)]
    TX->>TX: 计价(规格绝对价+加料加价) 落库order/order_item快照
    TX-->>POS: 下单成功(金额明细)
    TX->>MQ: 发事件 order.created
    MQ->>PR: 消费事件
    PR->>PR: 按 production_dept_type 拆单<br/>生成print_job(job_key幂等)
    PR->>GW: 下发打印任务(KITCHEN/BAR)
    GW->>K: 鱼香肉丝小票(ESC/POS)
    GW->>B: 珍珠奶茶小票(ESC/POS)
    K-->>GW: 打印回执
    B-->>GW: 打印回执
    GW-->>PR: 回执 PRINTED
```

---

## 4. 时序图：结账 → 收银 → 清台

```mermaid
sequenceDiagram
    participant POS as Android收银端
    participant TX as 交易服务
    participant CA as 收银服务
    participant GW as 边缘网关
    participant P as 收银小票打印机

    POS->>CA: 结账(orderId)
    CA->>CA: 算折扣/应收 生成bill
    CA-->>POS: 账单(应收金额)
    POS->>CA: 收款(CASH/SCAN_GUN, amount)
    CA->>CA: 记payment, 找零
    CA->>TX: 订单 SETTLED
    CA->>GW: 打印结账单(CHECKOUT)
    GW->>P: 结账小票
    POS->>TX: 清台(orderId)
    TX->>TX: 桌台 TO_CLEAN→IDLE
    TX-->>POS: 清台完成
```

---

## 5. 时序图：弱网/离线 — 边缘网关兜底与同步

```mermaid
sequenceDiagram
    participant POS as Android收银端
    participant GW as 门店边缘网关(本地)
    participant CLOUD as 云端
    participant PRN as 门店打印机

    Note over POS,GW: 公网中断
    POS->>GW: 下单/收款(局域网)
    GW->>GW: 本地落库(order/payment/print_job)
    GW->>PRN: 本地拆单打印(后厨/吧台)
    PRN-->>GW: 回执
    GW-->>POS: 出单/收银成功(离线)

    Note over GW,CLOUD: 公网恢复
    GW->>CLOUD: 幂等上行(order_no+操作序号)
    CLOUD->>CLOUD: 去重合并, 现金账以门店为权威
    CLOUD-->>GW: 回传权威状态
```

---

## 6. 时序图：退菜 + 补打

```mermaid
sequenceDiagram
    participant POS as Android收银端
    participant TX as 交易服务
    participant AUD as 审计
    participant GW as 边缘网关
    participant K as 后厨/吧台打印机

    POS->>TX: 退菜(orderItemId, 原因) [需refund权限]
    TX->>TX: 校验权限, item_status=REFUNDED, 重算金额
    TX->>AUD: 记录退菜审计
    TX->>GW: 退菜单(job_type=REFUND)
    GW->>K: 打印退菜通知
    POS->>GW: 补打(printJobId) [打印机异常时]
    GW->>K: 重打(job_key幂等, 不重复扣减)
```

---

## 7. 状态机：桌台 / 订单

```mermaid
stateDiagram-v2
    [*] --> IDLE
    IDLE --> OPENED : 开台
    OPENED --> OPENED : 加菜/转桌/并桌
    OPENED --> SETTLED : 结账
    SETTLED --> TO_CLEAN : 收款完成
    TO_CLEAN --> IDLE : 清台
    OPENED --> CANCELED : 取消(权限)
```

> 下一份将产出：完整接口清单 + 出入参 JSON 示例。
