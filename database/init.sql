CREATE DATABASE IF NOT EXISTS admin_system DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE admin_system;

CREATE TABLE IF NOT EXISTS `users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(50) NOT NULL COMMENT '用户名',
  `password` varchar(255) NOT NULL COMMENT '密码',
  `nickname` varchar(50) DEFAULT NULL COMMENT '昵称',
  `email` varchar(100) DEFAULT NULL COMMENT '邮箱',
  `phone` varchar(20) DEFAULT NULL COMMENT '手机号',
  `avatar` varchar(255) DEFAULT NULL COMMENT '头像',
  `status` tinyint NOT NULL DEFAULT 1 COMMENT '状态 1:启用 0:禁用',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

CREATE TABLE IF NOT EXISTS `roles` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(50) NOT NULL COMMENT '角色名称',
  `code` varchar(50) NOT NULL COMMENT '角色编码',
  `description` varchar(255) DEFAULT NULL COMMENT '描述',
  `status` tinyint NOT NULL DEFAULT 1 COMMENT '状态 1:启用 0:禁用',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色表';

CREATE TABLE IF NOT EXISTS `permissions` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(50) NOT NULL COMMENT '权限名称',
  `code` varchar(100) NOT NULL COMMENT '权限编码',
  `type` tinyint NOT NULL DEFAULT 1 COMMENT '类型 1:菜单 2:按钮',
  `parent_id` bigint unsigned NOT NULL DEFAULT 0 COMMENT '父级ID',
  `path` varchar(255) DEFAULT NULL COMMENT '路由路径',
  `method` varchar(10) DEFAULT NULL COMMENT '请求方法',
  `sort` int NOT NULL DEFAULT 0 COMMENT '排序',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='权限表';

CREATE TABLE IF NOT EXISTS `menus` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(50) NOT NULL COMMENT '菜单名称',
  `path` varchar(255) NOT NULL COMMENT '路由路径',
  `icon` varchar(50) DEFAULT NULL COMMENT '图标',
  `parent_id` bigint unsigned NOT NULL DEFAULT 0 COMMENT '父级ID',
  `component` varchar(255) DEFAULT NULL COMMENT '组件路径',
  `sort` int NOT NULL DEFAULT 0 COMMENT '排序',
  `visible` tinyint NOT NULL DEFAULT 1 COMMENT '是否可见 1:可见 0:隐藏',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='菜单表';

CREATE TABLE IF NOT EXISTS `user_roles` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL,
  `role_id` bigint unsigned NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_role` (`user_id`,`role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户角色关联表';

CREATE TABLE IF NOT EXISTS `role_permissions` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `role_id` bigint unsigned NOT NULL,
  `permission_id` bigint unsigned NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_role_permission` (`role_id`,`permission_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色权限关联表';

CREATE TABLE IF NOT EXISTS `role_menus` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `role_id` bigint unsigned NOT NULL,
  `menu_id` bigint unsigned NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_role_menu` (`role_id`,`menu_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色菜单关联表';

CREATE TABLE IF NOT EXISTS `operation_logs` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned DEFAULT NULL COMMENT '操作用户ID',
  `username` varchar(50) DEFAULT NULL COMMENT '操作用户名',
  `module` varchar(50) DEFAULT NULL COMMENT '操作模块',
  `operation` varchar(100) DEFAULT NULL COMMENT '操作类型',
  `method` varchar(10) DEFAULT NULL COMMENT '请求方法',
  `path` varchar(255) DEFAULT NULL COMMENT '请求路径',
  `ip` varchar(50) DEFAULT NULL COMMENT 'IP地址',
  `params` text COMMENT '请求参数',
  `result` text COMMENT '响应结果',
  `status` tinyint NOT NULL DEFAULT 1 COMMENT '状态 1:成功 0:失败',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='操作日志表';

INSERT INTO `roles` (`name`, `code`, `description`) VALUES
('超级管理员', 'super_admin', '系统最高权限管理员'),
('普通管理员', 'admin', '普通管理员'),
('普通用户', 'user', '普通用户');

INSERT INTO `menus` (`name`, `path`, `icon`, `parent_id`, `component`, `sort`) VALUES
('系统管理', '/system', 'setting', 0, NULL, 1),
('用户管理', '/system/users', 'user', 1, 'system/users/index', 1),
('角色管理', '/system/roles', 'user-role', 1, 'system/roles/index', 2),
('菜单管理', '/system/menus', 'menu', 1, 'system/menus/index', 3),
('日志管理', '/system/logs', 'document', 1, 'system/logs/index', 4);

INSERT INTO `permissions` (`name`, `code`, `type`, `parent_id`, `path`, `method`, `sort`) VALUES
('用户管理', 'user:manage', 1, 0, '/api/users', 'GET', 1),
('用户列表', 'user:list', 2, 1, '/api/users', 'GET', 1),
('新增用户', 'user:add', 2, 1, '/api/users', 'POST', 2),
('编辑用户', 'user:edit', 2, 1, '/api/users/:id', 'PUT', 3),
('删除用户', 'user:delete', 2, 1, '/api/users/:id', 'DELETE', 4),
('角色管理', 'role:manage', 1, 0, '/api/roles', 'GET', 2),
('角色列表', 'role:list', 2, 6, '/api/roles', 'GET', 1),
('新增角色', 'role:add', 2, 6, '/api/roles', 'POST', 2),
('编辑角色', 'role:edit', 2, 6, '/api/roles/:id', 'PUT', 3),
('删除角色', 'role:delete', 2, 6, '/api/roles/:id', 'DELETE', 4),
('分配权限', 'role:permission', 2, 6, '/api/roles/:id/permissions', 'POST', 5);
