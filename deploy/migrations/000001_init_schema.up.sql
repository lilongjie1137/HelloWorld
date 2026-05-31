-- V1.0 基线 schema（MySQL 8, InnoDB, utf8mb4）。
-- 约定：所有业务表含 id/tenant_id/created_at/updated_at；门店级表含 store_id；金额 DECIMAL(10,2)。
-- 多租户隔离走应用层强制 tenant_id 过滤，故不设跨租户硬外键，仅建索引（便于后续分库）。

SET NAMES utf8mb4;

-- ============ 基础 / 租户 ============
CREATE TABLE tenant (
  id          BIGINT       NOT NULL PRIMARY KEY,
  name        VARCHAR(128) NOT NULL,
  code        VARCHAR(64)  NOT NULL,
  status      TINYINT      NOT NULL DEFAULT 1 COMMENT '1启用 0停用',
  created_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uk_tenant_code (code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='品牌/租户';

CREATE TABLE store (
  id          BIGINT       NOT NULL PRIMARY KEY,
  tenant_id   BIGINT       NOT NULL,
  name        VARCHAR(128) NOT NULL,
  code        VARCHAR(64)  NOT NULL,
  address     VARCHAR(255) NOT NULL DEFAULT '',
  status      TINYINT      NOT NULL DEFAULT 1,
  created_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uk_store_code (tenant_id, code),
  KEY idx_store_tenant (tenant_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='门店';

CREATE TABLE production_dept (
  id          BIGINT       NOT NULL PRIMARY KEY,
  tenant_id   BIGINT       NOT NULL,
  store_id    BIGINT       NOT NULL,
  name        VARCHAR(64)  NOT NULL,
  type        VARCHAR(16)  NOT NULL COMMENT 'KITCHEN/BAR/COLD',
  created_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY idx_dept_store (tenant_id, store_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='出品部门';

CREATE TABLE device (
  id          BIGINT       NOT NULL PRIMARY KEY,
  tenant_id   BIGINT       NOT NULL,
  store_id    BIGINT       NOT NULL,
  type        VARCHAR(16)  NOT NULL COMMENT 'POS/PRINTER/LABEL/KDS',
  name        VARCHAR(64)  NOT NULL,
  sn          VARCHAR(64)  NOT NULL,
  status      TINYINT      NOT NULL DEFAULT 1,
  created_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uk_device_sn (tenant_id, sn),
  KEY idx_device_store (tenant_id, store_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='设备登记';

-- ============ 员工 / 权限 / 交接班 ============
CREATE TABLE employee (
  id          BIGINT       NOT NULL PRIMARY KEY,
  tenant_id   BIGINT       NOT NULL,
  name        VARCHAR(64)  NOT NULL,
  phone       VARCHAR(20)  NOT NULL DEFAULT '',
  login_name  VARCHAR(64)  NOT NULL,
  pwd_hash    VARCHAR(100) NOT NULL,
  status      TINYINT      NOT NULL DEFAULT 1,
  created_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uk_emp_login (tenant_id, login_name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='员工';

CREATE TABLE role (
  id          BIGINT       NOT NULL PRIMARY KEY,
  tenant_id   BIGINT       NOT NULL,
  name        VARCHAR(64)  NOT NULL,
  code        VARCHAR(64)  NOT NULL,
  created_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uk_role_code (tenant_id, code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色';

CREATE TABLE permission (
  id          BIGINT       NOT NULL PRIMARY KEY,
  code        VARCHAR(64)  NOT NULL COMMENT '如 order:refund',
  name        VARCHAR(64)  NOT NULL,
  created_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uk_perm_code (code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='权限点（全局字典）';

CREATE TABLE role_permission (
  role_id       BIGINT NOT NULL,
  permission_id BIGINT NOT NULL,
  PRIMARY KEY (role_id, permission_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色-权限';

CREATE TABLE employee_role (
  employee_id BIGINT NOT NULL,
  role_id     BIGINT NOT NULL,
  PRIMARY KEY (employee_id, role_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='员工-角色';

CREATE TABLE employee_store (
  employee_id BIGINT NOT NULL,
  store_id    BIGINT NOT NULL,
  PRIMARY KEY (employee_id, store_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='员工数据范围:可操作门店';

CREATE TABLE shift (
  id                 BIGINT        NOT NULL PRIMARY KEY,
  tenant_id          BIGINT        NOT NULL,
  store_id           BIGINT        NOT NULL,
  employee_id        BIGINT        NOT NULL,
  open_time          DATETIME      NOT NULL,
  close_time         DATETIME      NULL,
  open_cash          DECIMAL(10,2) NOT NULL DEFAULT 0,
  close_cash         DECIMAL(10,2) NOT NULL DEFAULT 0,
  total_cash         DECIMAL(10,2) NOT NULL DEFAULT 0,
  total_order_amount DECIMAL(10,2) NOT NULL DEFAULT 0,
  status             VARCHAR(16)   NOT NULL DEFAULT 'OPEN' COMMENT 'OPEN/CLOSED',
  created_at         DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at         DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY idx_shift_store (tenant_id, store_id, status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='交接班';

-- ============ 商品中心（全国统一价）============
CREATE TABLE category (
  id          BIGINT       NOT NULL PRIMARY KEY,
  tenant_id   BIGINT       NOT NULL,
  name        VARCHAR(64)  NOT NULL,
  sort        INT          NOT NULL DEFAULT 0,
  parent_id   BIGINT       NOT NULL DEFAULT 0,
  created_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY idx_cat_tenant (tenant_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='分类';

CREATE TABLE spu (
  id                   BIGINT        NOT NULL PRIMARY KEY,
  tenant_id            BIGINT        NOT NULL,
  category_id          BIGINT        NOT NULL,
  name                 VARCHAR(128)  NOT NULL,
  base_price           DECIMAL(10,2) NOT NULL DEFAULT 0 COMMENT '全国统一价',
  production_dept_type VARCHAR(16)   NOT NULL COMMENT 'KITCHEN/BAR/COLD',
  is_combo             TINYINT       NOT NULL DEFAULT 0,
  is_weighable         TINYINT       NOT NULL DEFAULT 0,
  image_url            VARCHAR(255)  NOT NULL DEFAULT '',
  status               TINYINT       NOT NULL DEFAULT 1,
  sort                 INT           NOT NULL DEFAULT 0,
  created_at           DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at           DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY idx_spu_cat (tenant_id, category_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='SPU 商品';

CREATE TABLE modifier_group (
  id          BIGINT       NOT NULL PRIMARY KEY,
  tenant_id   BIGINT       NOT NULL,
  name        VARCHAR(64)  NOT NULL,
  group_type  VARCHAR(16)  NOT NULL COMMENT 'SPEC/METHOD/ADDON',
  select_type VARCHAR(16)  NOT NULL DEFAULT 'SINGLE' COMMENT 'SINGLE/MULTI',
  required    TINYINT      NOT NULL DEFAULT 0,
  created_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY idx_mg_tenant (tenant_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='属性组:规格/做法/加料';

CREATE TABLE modifier_option (
  id          BIGINT        NOT NULL PRIMARY KEY,
  tenant_id   BIGINT        NOT NULL,
  group_id    BIGINT        NOT NULL,
  name        VARCHAR(64)   NOT NULL,
  price       DECIMAL(10,2) NOT NULL DEFAULT 0 COMMENT 'SPEC=绝对单价/ADDON=加价/METHOD=0',
  is_default  TINYINT       NOT NULL DEFAULT 0,
  sort        INT           NOT NULL DEFAULT 0,
  status      TINYINT       NOT NULL DEFAULT 1,
  created_at  DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at  DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY idx_mo_group (tenant_id, group_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='属性选项';

CREATE TABLE spu_modifier_group (
  id          BIGINT NOT NULL PRIMARY KEY,
  tenant_id   BIGINT NOT NULL,
  spu_id      BIGINT NOT NULL,
  group_id    BIGINT NOT NULL,
  sort        INT    NOT NULL DEFAULT 0,
  UNIQUE KEY uk_spu_group (spu_id, group_id),
  KEY idx_smg_tenant (tenant_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='SPU-属性组关联';

CREATE TABLE spu_store (
  id          BIGINT  NOT NULL PRIMARY KEY,
  tenant_id   BIGINT  NOT NULL,
  spu_id      BIGINT  NOT NULL,
  store_id    BIGINT  NOT NULL,
  available   TINYINT NOT NULL DEFAULT 1,
  UNIQUE KEY uk_spu_store (spu_id, store_id),
  KEY idx_ss_store (tenant_id, store_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='菜单门店适用范围';

-- ============ 桌台 ============
CREATE TABLE table_area (
  id          BIGINT       NOT NULL PRIMARY KEY,
  tenant_id   BIGINT       NOT NULL,
  store_id    BIGINT       NOT NULL,
  name        VARCHAR(64)  NOT NULL,
  sort        INT          NOT NULL DEFAULT 0,
  created_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY idx_area_store (tenant_id, store_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='桌台区域';

CREATE TABLE dining_table (
  id          BIGINT       NOT NULL PRIMARY KEY,
  tenant_id   BIGINT       NOT NULL,
  store_id    BIGINT       NOT NULL,
  area_id     BIGINT       NOT NULL,
  name        VARCHAR(64)  NOT NULL,
  capacity    INT          NOT NULL DEFAULT 4,
  status      VARCHAR(16)  NOT NULL DEFAULT 'IDLE' COMMENT 'IDLE/OPENED/TO_CLEAN',
  qr_code     VARCHAR(128) NOT NULL DEFAULT '',
  created_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY idx_table_store (tenant_id, store_id, status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='桌台';

-- ============ 订单（下单即快照）============
CREATE TABLE `order` (
  id              BIGINT        NOT NULL PRIMARY KEY,
  tenant_id       BIGINT        NOT NULL,
  store_id        BIGINT        NOT NULL,
  table_id        BIGINT        NOT NULL DEFAULT 0,
  order_no        VARCHAR(40)   NOT NULL,
  source          VARCHAR(16)   NOT NULL DEFAULT 'POS' COMMENT 'POS/MINI_APP/TAKEOUT',
  status          VARCHAR(16)   NOT NULL DEFAULT 'OPEN',
  people_count    INT           NOT NULL DEFAULT 0,
  total_amount    DECIMAL(10,2) NOT NULL DEFAULT 0,
  discount_amount DECIMAL(10,2) NOT NULL DEFAULT 0,
  payable_amount  DECIMAL(10,2) NOT NULL DEFAULT 0,
  employee_id     BIGINT        NOT NULL DEFAULT 0,
  opened_at       DATETIME      NULL,
  settled_at      DATETIME      NULL,
  created_at      DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at      DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uk_order_no (tenant_id, order_no),
  KEY idx_order_store (tenant_id, store_id, status),
  KEY idx_order_table (tenant_id, table_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='订单';

CREATE TABLE order_item (
  id                        BIGINT        NOT NULL PRIMARY KEY,
  tenant_id                 BIGINT        NOT NULL,
  order_id                  BIGINT        NOT NULL,
  spu_id                    BIGINT        NOT NULL,
  spu_name_snap             VARCHAR(128)  NOT NULL,
  category_snap             VARCHAR(64)   NOT NULL DEFAULT '',
  spec_option_id            BIGINT        NOT NULL DEFAULT 0,
  spec_name_snap            VARCHAR(64)   NOT NULL DEFAULT '',
  unit_price                DECIMAL(10,2) NOT NULL DEFAULT 0,
  qty                       INT           NOT NULL DEFAULT 1,
  amount                    DECIMAL(10,2) NOT NULL DEFAULT 0,
  remark                    VARCHAR(255)  NOT NULL DEFAULT '',
  production_dept_type_snap VARCHAR(16)   NOT NULL COMMENT '出品部门快照→打印分组',
  item_status               VARCHAR(16)   NOT NULL DEFAULT 'PENDING' COMMENT 'PENDING/SERVED/REFUNDED',
  created_at                DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at                DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY idx_oi_order (tenant_id, order_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='订单明细(快照)';

CREATE TABLE order_item_modifier (
  id               BIGINT        NOT NULL PRIMARY KEY,
  tenant_id        BIGINT        NOT NULL,
  order_item_id    BIGINT        NOT NULL,
  group_type       VARCHAR(16)   NOT NULL COMMENT 'SPEC/METHOD/ADDON',
  option_id        BIGINT        NOT NULL,
  option_name_snap VARCHAR(64)   NOT NULL,
  price_snap       DECIMAL(10,2) NOT NULL DEFAULT 0,
  KEY idx_oim_item (tenant_id, order_item_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='订单明细属性(逐项快照)';

-- ============ 收银 / 账单 ============
CREATE TABLE bill (
  id              BIGINT        NOT NULL PRIMARY KEY,
  tenant_id       BIGINT        NOT NULL,
  store_id        BIGINT        NOT NULL,
  order_id        BIGINT        NOT NULL,
  bill_no         VARCHAR(40)   NOT NULL,
  total_amount    DECIMAL(10,2) NOT NULL DEFAULT 0,
  discount_amount DECIMAL(10,2) NOT NULL DEFAULT 0,
  payable_amount  DECIMAL(10,2) NOT NULL DEFAULT 0,
  paid_amount     DECIMAL(10,2) NOT NULL DEFAULT 0,
  change_amount   DECIMAL(10,2) NOT NULL DEFAULT 0,
  status          VARCHAR(16)   NOT NULL DEFAULT 'UNPAID' COMMENT 'UNPAID/PAID/REVERSED',
  cashier_id      BIGINT        NOT NULL DEFAULT 0,
  settled_at      DATETIME      NULL,
  created_at      DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at      DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uk_bill_no (tenant_id, bill_no),
  KEY idx_bill_order (tenant_id, order_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='账单';

CREATE TABLE payment (
  id          BIGINT        NOT NULL PRIMARY KEY,
  tenant_id   BIGINT        NOT NULL,
  bill_id     BIGINT        NOT NULL,
  pay_method  VARCHAR(16)   NOT NULL COMMENT 'CASH/SCAN_GUN/MEMBER/AGGREGATE',
  amount      DECIMAL(10,2) NOT NULL DEFAULT 0,
  trade_no    VARCHAR(64)   NOT NULL DEFAULT '',
  status      VARCHAR(16)   NOT NULL DEFAULT 'SUCCESS',
  paid_at     DATETIME      NULL,
  created_at  DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at  DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY idx_pay_bill (tenant_id, bill_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='收款';

CREATE TABLE discount (
  id          BIGINT        NOT NULL PRIMARY KEY,
  tenant_id   BIGINT        NOT NULL,
  bill_id     BIGINT        NOT NULL,
  type        VARCHAR(16)   NOT NULL COMMENT 'PERCENT/AMOUNT/ROUND',
  value       DECIMAL(10,2) NOT NULL DEFAULT 0,
  reason      VARCHAR(128)  NOT NULL DEFAULT '',
  operator_id BIGINT        NOT NULL,
  created_at  DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP,
  KEY idx_disc_bill (tenant_id, bill_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='折扣/抹零';

-- ============ 打印 / 路由 ============
CREATE TABLE printer (
  id          BIGINT       NOT NULL PRIMARY KEY,
  tenant_id   BIGINT       NOT NULL,
  store_id    BIGINT       NOT NULL,
  name        VARCHAR(64)  NOT NULL,
  ip          VARCHAR(64)  NOT NULL DEFAULT '',
  model       VARCHAR(64)  NOT NULL DEFAULT '',
  proto       VARCHAR(16)  NOT NULL DEFAULT 'ESC_POS' COMMENT 'ESC_POS/TSPL',
  type        VARCHAR(16)  NOT NULL DEFAULT 'RECEIPT' COMMENT 'RECEIPT/LABEL',
  status      TINYINT      NOT NULL DEFAULT 1,
  created_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY idx_printer_store (tenant_id, store_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='打印机';

CREATE TABLE printer_dept (
  printer_id           BIGINT      NOT NULL,
  production_dept_type VARCHAR(16) NOT NULL,
  PRIMARY KEY (printer_id, production_dept_type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='打印机-出品部门类型绑定';

CREATE TABLE print_job (
  id                   BIGINT       NOT NULL PRIMARY KEY,
  tenant_id            BIGINT       NOT NULL,
  store_id             BIGINT       NOT NULL,
  order_id             BIGINT       NOT NULL,
  printer_id           BIGINT       NOT NULL DEFAULT 0,
  production_dept_type VARCHAR(16)  NOT NULL DEFAULT '',
  job_type             VARCHAR(16)  NOT NULL COMMENT 'KITCHEN/BAR/LABEL/CHECKOUT/REFUND',
  payload              JSON         NULL,
  status               VARCHAR(16)  NOT NULL DEFAULT 'PENDING' COMMENT 'PENDING/SENT/PRINTED/FAILED',
  retry_count          INT          NOT NULL DEFAULT 0,
  job_key              VARCHAR(80)  NOT NULL COMMENT '幂等键 防重打',
  acked_at             DATETIME     NULL,
  created_at           DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at           DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uk_job_key (tenant_id, job_key),
  KEY idx_job_store_status (tenant_id, store_id, status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='打印任务';

-- ============ 预留：会员 / 储值（P1 建表, P3 实现）============
CREATE TABLE member (
  id          BIGINT       NOT NULL PRIMARY KEY,
  tenant_id   BIGINT       NOT NULL,
  phone       VARCHAR(20)  NOT NULL,
  name        VARCHAR(64)  NOT NULL DEFAULT '',
  wx_openid   VARCHAR(64)  NOT NULL DEFAULT '',
  level       INT          NOT NULL DEFAULT 0,
  status      TINYINT      NOT NULL DEFAULT 1,
  created_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uk_member_phone (tenant_id, phone)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='会员(预留)';

CREATE TABLE member_card (
  id          BIGINT       NOT NULL PRIMARY KEY,
  tenant_id   BIGINT       NOT NULL,
  member_id   BIGINT       NOT NULL,
  card_no     VARCHAR(40)  NOT NULL,
  status      TINYINT      NOT NULL DEFAULT 1,
  created_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uk_card_no (tenant_id, card_no)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='会员卡(预留)';

CREATE TABLE stored_value_account (
  id           BIGINT        NOT NULL PRIMARY KEY,
  tenant_id    BIGINT        NOT NULL,
  member_id    BIGINT        NOT NULL,
  balance      DECIMAL(10,2) NOT NULL DEFAULT 0,
  gift_balance DECIMAL(10,2) NOT NULL DEFAULT 0,
  created_at   DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at   DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uk_sva_member (tenant_id, member_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='储值账户(预留)';

CREATE TABLE stored_value_txn (
  id            BIGINT        NOT NULL PRIMARY KEY,
  tenant_id     BIGINT        NOT NULL,
  account_id    BIGINT        NOT NULL,
  type          VARCHAR(16)   NOT NULL COMMENT 'RECHARGE/CONSUME/GIFT',
  amount        DECIMAL(10,2) NOT NULL DEFAULT 0,
  balance_after DECIMAL(10,2) NOT NULL DEFAULT 0,
  ref_order_id  BIGINT        NOT NULL DEFAULT 0,
  remark        VARCHAR(128)  NOT NULL DEFAULT '',
  created_at    DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP,
  KEY idx_svt_account (tenant_id, account_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='储值流水(预留)';
