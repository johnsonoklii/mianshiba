-- 用户表
CREATE TABLE `user` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `username` VARCHAR(50) NOT NULL,
  `email` VARCHAR(100) NOT NULL,
  `password` VARCHAR(255) NOT NULL,
  `role` VARCHAR(20) NOT NULL DEFAULT 'user',
  `avatar` VARCHAR(255) DEFAULT NULL COMMENT '头像',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at` DATETIME(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_users_username` (`username`),
  UNIQUE KEY `uk_users_email` (`email`),
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci;

-- 用户模型表
CREATE TABLE `user_model` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `user_id` BIGINT NOT NULL COMMENT '用户ID',
  `name` VARCHAR(128) NOT NULL COMMENT '模型显示名称（用户维度唯一）',
  `model_key` VARCHAR(128) NOT NULL COMMENT '模型标识（doubao-1.5-vision-lite-250315）',
  `protocol` VARCHAR(64) NOT NULL COMMENT '协议类型（openai/ark/claude/gemini/deepseek/ollama/qwen/ernie）',
  `base_url` VARCHAR(255) NOT NULL COMMENT 'API 基础地址',
  `api_key_encrypted` TEXT NOT NULL COMMENT '加密后的 API 密钥',
  `config_json` JSON DEFAULT NULL COMMENT '额外配置（如区域、访问密钥等）',
  `secret_hint` VARCHAR(32) NOT NULL DEFAULT '' COMMENT '密钥脱敏提示（如显示末尾4位）',
  `provider_name` VARCHAR(64) NOT NULL COMMENT '提供商名称（如 OpenAI、Ark、DeepSeek）',
  `meta_id` BIGINT DEFAULT NULL COMMENT '关联全局 model_meta.id（继承能力/图标）',
  `default_params` JSON DEFAULT NULL COMMENT '默认参数（如 temperature、max_tokens）',
  `scope` INT NOT NULL DEFAULT 7 COMMENT '使用范围（位掩码：1=智能体, 2=应用, 4=工作流）',
  `status` INT NOT NULL DEFAULT 1 COMMENT '状态（0=禁用, 1=启用）',
  `is_default` INT NOT NULL DEFAULT 0 COMMENT '是否为默认（0=不是, 1=是）',
  `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间（毫秒时间戳）',
  `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间（毫秒时间戳）',
  `deleted` TINYINT NOT NULL DEFAULT 0 COMMENT '删除状态（0=未删除, 1=已删除）',

  PRIMARY KEY (`id`),

  KEY `idx_user_models_user_id` (`user_id`),
  KEY `idx_user_models_status` (`status`),
  KEY `idx_user_models_deleted` (`deleted`),
  KEY `idx_user_models_scope` (`scope`)

) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci;