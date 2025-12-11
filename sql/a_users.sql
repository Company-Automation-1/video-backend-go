CREATE TABLE `a_users`  (
  `id` int UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `username` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '用户名',
  `email` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '邮箱',
  `password` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '密码',
  `created_at` bigint NOT NULL COMMENT '创建时间：秒级时间戳',
  `updated_at` bigint NOT NULL COMMENT '更新时间：秒级时间戳',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_email`(`email` ASC) USING BTREE COMMENT '邮箱',
  UNIQUE INDEX `idx_username`(`username` ASC) USING BTREE COMMENT '用户名'
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;
