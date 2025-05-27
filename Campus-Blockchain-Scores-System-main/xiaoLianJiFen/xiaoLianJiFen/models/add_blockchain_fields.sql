-- 添加区块链交易哈希和时间戳字段到users表
ALTER TABLE users
ADD COLUMN tx_hash VARCHAR(66) NULL COMMENT '区块链交易哈希',
ADD COLUMN block_timestamp BIGINT DEFAULT 0 COMMENT '区块链交易时间戳'; 
ADD COLUMN block_number BIGINT DEFAULT 0 COMMENT '区块号';