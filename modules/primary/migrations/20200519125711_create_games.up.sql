CREATE TABLE games (
    game_id varchar(36) not null,

    name varchar(256) not null,

    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),

    PRIMARY KEY(game_id)
);


-- TODO: use enum
-- platform, market, device_type
CREATE TABLE apps (
    -- GUID
    game_id varchar(36) not null REFERENCES games(game_id),
    -- GUID
    app_id varchar(36) not null,

    version varchar(64) not null,

    -- android, ios, web
    platform varchar(64) not null,

    -- app-store, huawei, taptap, facebook, facebook-messenger
    market varchar(128) not null,

    -- tv, desktop, mobile, pad, browser
    device_type varchar(64) not null,

    -- on enable the client will request to update the app
    -- by using alert on game boot
    force_update_enabled BOOLEAN not null default false,

    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),

    PRIMARY KEY(game_id, app_id)
);

CREATE TABLE networks (
    type_name varchar(64) not null,

    id varchar(128) not null,
    secret varchar(256) not null,

    -- GUID
    game_id varchar(36) not null REFERENCES games(game_id),

    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),

    PRIMARY KEY(type_name, id, game_id)
);
