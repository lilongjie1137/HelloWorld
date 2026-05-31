-- 种子数据：权限点字典（全局）+ 演示租户/门店/部门/角色/员工/商品/打印机/桌台。
-- 用于本地开发与联调；生产环境仅需保留 permission 字典部分。
-- 演示账号：login=admin / 密码=admin123（bcrypt）。

SET NAMES utf8mb4;

-- ===== 权限点字典（全局，生产保留）=====
INSERT INTO permission (id, code, name) VALUES
  (1,  'store:manage',      '门店管理'),
  (2,  'employee:manage',   '员工/权限管理'),
  (3,  'catalog:manage',    '商品管理'),
  (4,  'table:open',        '开台'),
  (5,  'table:transfer',    '转桌'),
  (6,  'table:merge',       '并桌'),
  (7,  'order:create',      '点单/加菜'),
  (8,  'order:refund',      '退菜'),
  (9,  'order:reopen',      '反结账'),
  (10, 'item:changeprice',  '改价'),
  (11, 'bill:checkout',     '结账收款'),
  (12, 'bill:discount',     '折扣'),
  (13, 'bill:round',        '抹零'),
  (14, 'shift:open',        '开班'),
  (15, 'shift:close',       '交接班'),
  (16, 'print:reprint',     '补打');

-- ===== 演示租户 / 门店 / 出品部门 =====
INSERT INTO tenant (id, name, code, status) VALUES (1, '示例餐饮品牌', 'DEMO', 1);
INSERT INTO store  (id, tenant_id, name, code, address, status) VALUES (1, 1, '上海旗舰店', 'SH001', '上海市黄浦区示例路1号', 1);
INSERT INTO production_dept (id, tenant_id, store_id, name, type) VALUES
  (1, 1, 1, '热菜间', 'KITCHEN'),
  (2, 1, 1, '吧台',   'BAR');

-- ===== 角色 / 权限 / 员工 =====
INSERT INTO role (id, tenant_id, name, code) VALUES
  (1, 1, '门店管理员', 'STORE_ADMIN'),
  (2, 1, '收银员',     'CASHIER'),
  (3, 1, '服务员',     'WAITER');

-- 管理员拥有全部权限
INSERT INTO role_permission (role_id, permission_id)
  SELECT 1, id FROM permission;
-- 收银员：开台/点单/结账/折扣/抹零/补打/开班交接
INSERT INTO role_permission (role_id, permission_id) VALUES
  (2,4),(2,7),(2,11),(2,12),(2,13),(2,16),(2,14),(2,15);
-- 服务员：开台/转桌/并桌/点单/退菜
INSERT INTO role_permission (role_id, permission_id) VALUES
  (3,4),(3,5),(3,6),(3,7),(3,8);

INSERT INTO employee (id, tenant_id, name, phone, login_name, pwd_hash, status) VALUES
  (1, 1, '管理员', '13800000000', 'admin', '$2a$10$AXw3nuAYh29bPcwbTgUUbOkNGkgaT/qoXbVssheMNI4wLwCbnGHfO', 1);
INSERT INTO employee_role  (employee_id, role_id)  VALUES (1, 1);
INSERT INTO employee_store (employee_id, store_id) VALUES (1, 1);

-- ===== 商品：分类 / SPU / 属性组 / 选项 =====
INSERT INTO category (id, tenant_id, name, sort, parent_id) VALUES
  (1, 1, '热菜', 1, 0),
  (2, 1, '饮品', 2, 0);

-- 鱼香肉丝（后厨，base_price=28.00，无规格）
INSERT INTO spu (id, tenant_id, category_id, name, base_price, production_dept_type, status, sort) VALUES
  (1, 1, 1, '鱼香肉丝', 28.00, 'KITCHEN', 1, 1);
-- 珍珠奶茶（吧台，靠规格绝对价，base_price=0）
INSERT INTO spu (id, tenant_id, category_id, name, base_price, production_dept_type, status, sort) VALUES
  (2, 1, 2, '珍珠奶茶', 0.00, 'BAR', 1, 1);

-- 属性组
INSERT INTO modifier_group (id, tenant_id, name, group_type, select_type, required) VALUES
  (1, 1, '杯型', 'SPEC',   'SINGLE', 1),  -- 规格(必选,绝对价)
  (2, 1, '加料', 'ADDON',  'MULTI',  0),  -- 加料
  (3, 1, '甜度冰度', 'METHOD', 'MULTI', 0), -- 做法(不计价)
  (4, 1, '菜品做法', 'METHOD', 'MULTI', 0);

-- 选项
INSERT INTO modifier_option (id, tenant_id, group_id, name, price, is_default, sort) VALUES
  (1, 1, 1, '大杯', 15.00, 0, 1),
  (2, 1, 1, '中杯', 12.00, 1, 2),
  (3, 1, 2, '加珍珠', 3.00, 0, 1),
  (4, 1, 2, '加椰果', 2.00, 0, 2),
  (5, 1, 3, '去冰', 0.00, 0, 1),
  (6, 1, 3, '半糖', 0.00, 0, 2),
  (7, 1, 4, '免葱', 0.00, 0, 1),
  (8, 1, 4, '少辣', 0.00, 0, 2);

-- SPU-属性组关联
INSERT INTO spu_modifier_group (id, tenant_id, spu_id, group_id, sort) VALUES
  (1, 1, 2, 1, 1),  -- 奶茶-杯型
  (2, 1, 2, 2, 2),  -- 奶茶-加料
  (3, 1, 2, 3, 3),  -- 奶茶-甜度冰度
  (4, 1, 1, 4, 1);  -- 鱼香肉丝-做法

-- 门店适用范围
INSERT INTO spu_store (id, tenant_id, spu_id, store_id, available) VALUES
  (1, 1, 1, 1, 1),
  (2, 1, 2, 1, 1);

-- ===== 打印机 + 出品部门绑定 =====
INSERT INTO printer (id, tenant_id, store_id, name, ip, model, proto, type, status) VALUES
  (1, 1, 1, '后厨打印机', '192.168.1.51', 'GP-58', 'ESC_POS', 'RECEIPT', 1),
  (2, 1, 1, '吧台打印机', '192.168.1.52', 'GP-58', 'ESC_POS', 'RECEIPT', 1);
INSERT INTO printer_dept (printer_id, production_dept_type) VALUES
  (1, 'KITCHEN'),
  (2, 'BAR');

-- ===== 桌台 =====
INSERT INTO table_area (id, tenant_id, store_id, name, sort) VALUES (1, 1, 1, '大厅', 1);
INSERT INTO dining_table (id, tenant_id, store_id, area_id, name, capacity, status) VALUES
  (1, 1, 1, 1, '1号桌', 4, 'IDLE'),
  (2, 1, 1, 1, '2号桌', 4, 'IDLE'),
  (3, 1, 1, 1, '3号桌', 6, 'IDLE');
