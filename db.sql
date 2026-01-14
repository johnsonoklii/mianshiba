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
  `deleted` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '删除状态（0=未删除, 1=已删除）',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_users_username` (`username`),
  UNIQUE KEY `uk_users_email` (`email`)
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
   created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
        ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at TIMESTAMP NULL COMMENT '软删除时间',
  `deleted` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '删除状态（0=未删除, 1=已删除）',

  PRIMARY KEY (`id`),

  KEY `idx_user_models_user_id` (`user_id`),
  KEY `idx_user_models_status` (`status`),
  KEY `idx_user_models_deleted` (`deleted`),
  KEY `idx_user_models_scope` (`scope`)

) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci;


-- 简历
CREATE TABLE resume (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    user_id BIGINT NOT NULL COMMENT '用户ID',

    file_key VARCHAR(512) NOT NULL COMMENT '对象存储中的文件唯一标识',
    filename VARCHAR(255) NOT NULL COMMENT '原始文件名',
    filetype VARCHAR(50) COMMENT '文件类型，如 pdf/docx',
    filesize BIGINT COMMENT '文件大小（字节）',

    status TINYINT NOT NULL DEFAULT 1 COMMENT '简历状态：1已上传 2解析中 3已解析 4已删除 5失败',
    parse_status TINYINT NOT NULL DEFAULT 0 COMMENT '解析状态：0未开始 1解析中 2成功 3失败',
    parse_error VARCHAR(1024) COMMENT '解析失败原因摘要',

    update_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '上传时间',
    deleted_at TIMESTAMP NULL COMMENT '软删除时间',
    deleted TINYINT(1) NOT NULL DEFAULT 0 COMMENT '删除状态（0=未删除, 1=已删除）',

    PRIMARY KEY (id),
    UNIQUE KEY uk_file_key (file_key),
    KEY idx_user_id (user_id),
    KEY idx_user_status (user_id, status),
    KEY idx_parse_status (parse_status)
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci
  COMMENT='用户简历元信息表';
