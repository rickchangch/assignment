-- 便於測試用的預設資料

INSERT INTO users (id, address) 
VALUES 
    ('csf4eh04v3c4kmbslpt0', '0x1234567890abcdef1234567890abcdef12345678'),
    ('zsf4eh04v3c4kmbslpt0', '0x2234567890abcdef1234567890abcdef12345678');

INSERT INTO campaigns (id, name, started_at, ended_at, is_distributed) 
VALUES 
    ('esf4eh04v3c4kmbslpt0', 'Uniswap Campaign', 1729468800000, 1731887999000, false); -- 10/21 - 11/17

INSERT INTO campaign_tasks (id, campaign_id, pair, points, serial, started_at, ended_at, is_distributed) 
VALUES 
    ('gsf4eh04v3c4kmbslpt0', 'esf4eh04v3c4kmbslpt0', 'USDC/ETH', 10000, 1, 1729468800000, 1730073599000, false),
    ('hsf4eh04v3c4kmbslpt0', 'esf4eh04v3c4kmbslpt0', 'USDC/ETH', 10000, 2, 1730073600000, 1730678399000, false),
    ('isf4eh04v3c4kmbslpt0', 'esf4eh04v3c4kmbslpt0', 'USDC/ETH', 10000, 3, 1730678400000, 1731283199000, false),
    ('jsf4eh04v3c4kmbslpt0', 'esf4eh04v3c4kmbslpt0', 'USDC/ETH', 10000, 4, 1731283200000, 1731887999000, false);