-- users
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(20) PRIMARY KEY,
    address VARCHAR(42) NOT NULL,
    EXCLUDE USING HASH (address WITH =) -- HASH INDEX
);
COMMENT ON TABLE users IS '用戶';
COMMENT ON COLUMN users.id IS '用戶 ID';
COMMENT ON COLUMN users.address IS '註冊的錢包地址';

-- campaigns
CREATE TABLE IF NOT EXISTS campaigns (
    id VARCHAR(20) PRIMARY KEY,
    name VARCHAR(32) NOT NULL,
    started_at BIGINT NOT NULL,
    ended_at BIGINT NOT NULL,
    is_distributed BOOLEAN NOT NULL DEFAULT false
);
COMMENT ON TABLE campaigns IS '活動';
COMMENT ON COLUMN campaigns.id IS '活動 ID';
COMMENT ON COLUMN campaigns.name IS '活動名稱';
COMMENT ON COLUMN campaigns.started_at IS '活動開始時間';
COMMENT ON COLUMN campaigns.ended_at IS '活動結束時間';
COMMENT ON COLUMN campaigns.is_distributed IS 'true: 已分發獎勵; false: 未分發獎勵';

-- user_campaigns
CREATE TABLE IF NOT EXISTS user_campaigns (
    user_id VARCHAR(20) NOT NULL,
    campaign_id VARCHAR(20) NOT NULL,
    is_completed BOOLEAN NOT NULL DEFAULT FALSE,
    amount DECIMAL(18, 2) NOT NULL DEFAULT 0.0,
    points BIGINT NOT NULL DEFAULT 0,
    PRIMARY KEY (user_id, campaign_id)
);
COMMENT ON TABLE user_campaigns IS '用戶參與的活動紀錄';
COMMENT ON COLUMN user_campaigns.user_id IS '用戶ID';
COMMENT ON COLUMN user_campaigns.campaign_id IS '活動ID';
COMMENT ON COLUMN user_campaigns.is_completed IS '用戶是否完成 onboarding task';
COMMENT ON COLUMN user_campaigns.amount IS '用戶在活動中累計的 swap 金額';
COMMENT ON COLUMN user_campaigns.points IS '用戶在活動中累計的 points';

-- campaign_tasks
CREATE TABLE IF NOT EXISTS campaign_tasks (
    id VARCHAR(20) PRIMARY KEY,
    campaign_id VARCHAR(20) NOT NULL,
    pair VARCHAR(32) NOT NULL,
    points BIGINT NOT NULL,
    serial SMALLINT NOT NULL,
    started_at BIGINT NOT NULL,
    ended_at BIGINT NOT NULL,
    is_distributed BOOLEAN NOT NULL DEFAULT false
);
COMMENT ON TABLE campaign_tasks IS '活動任務, 即 Share Pool Tasks';
COMMENT ON COLUMN campaign_tasks.id IS '活動任務 ID';
COMMENT ON COLUMN campaign_tasks.campaign_id IS '活動 ID';
COMMENT ON COLUMN campaign_tasks.pair IS '資產對，例如 USDC/ETH';
COMMENT ON COLUMN campaign_tasks.points IS '此 Share Pool Task 可分配之點數';
COMMENT ON COLUMN campaign_tasks.started_at IS '開始時間';
COMMENT ON COLUMN campaign_tasks.ended_at IS '結束時間';
COMMENT ON COLUMN campaign_tasks.is_distributed IS 'true: 已分發點數; false: 未分發點數';
CREATE INDEX idx_campaign_id_pair ON campaign_tasks (campaign_id, pair);

-- user_campaign_tasks
CREATE TABLE IF NOT EXISTS user_campaign_tasks (
    user_id VARCHAR(20) NOT NULL,
    campaign_task_id VARCHAR(20) NOT NULL,
    campaign_id VARCHAR(20) NOT NULL,
    amount DECIMAL(18, 2) NOT NULL DEFAULT 0.0,
    points BIGINT NOT NULL DEFAULT 0,
    PRIMARY KEY (user_id, campaign_task_id)
);
COMMENT ON TABLE user_campaign_tasks IS '用戶參與的 Share Pool Tasks 紀錄';
COMMENT ON COLUMN user_campaign_tasks.user_id IS '用戶ID';
COMMENT ON COLUMN user_campaign_tasks.campaign_task_id IS '活動任務 ID';
COMMENT ON COLUMN user_campaign_tasks.campaign_id IS '活動ID';
COMMENT ON COLUMN user_campaign_tasks.amount IS 'swap 交易量';
COMMENT ON COLUMN user_campaign_tasks.points IS '獲取點數量';
CREATE INDEX idx_user_id_campaign_id ON user_campaign_tasks (user_id, campaign_id);

-- point_histories
CREATE TABLE IF NOT EXISTS point_histories (
    id VARCHAR(20) PRIMARY KEY,
    user_id VARCHAR(20) NOT NULL,
    points BIGINT NOT NULL,
    campaign_id VARCHAR(20) NOT NULL,
    campaign_task_id VARCHAR(20) NOT NULL,
    created_at BIGINT NOT NULL
);
COMMENT ON TABLE point_histories IS '用戶獲取點數歷程';
COMMENT ON COLUMN point_histories.id IS '點數歷程 ID';
COMMENT ON COLUMN point_histories.user_id IS '用戶 ID';
COMMENT ON COLUMN point_histories.points IS '點數';
COMMENT ON COLUMN point_histories.campaign_id IS '活動 ID';
COMMENT ON COLUMN point_histories.campaign_task_id IS '活動任務 ID';
COMMENT ON COLUMN point_histories.created_at IS '建立時間';
