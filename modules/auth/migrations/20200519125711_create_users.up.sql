CREATE TABLE users (
    user_id varchar(36) not null,
    guest_id varchar(36) not null,
    device_id varchar(256) not null,

    game_id varchar(36) not null,
    app_id varchar(36) not null,

    -- primary columns on sign in via
    -- network.
    email varchar(256),
    name varchar(256),

    network_name varchar(128) not null,
    network_id varchar(128) not null,

    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),

    PRIMARY KEY(game_id, app_id, user_id, network_id, network_name)
);
COMMENT ON TABLE users IS 'Users per game, app';

CREATE INDEX idx_users_user_id ON users(user_id);
CREATE INDEX idx_users_game_id_app_id_user_id ON users(game_id, app_id, user_id);


CREATE TABLE anonymouses (
    -- GUID
    user_id varchar(36) not null,
    -- GUID
    game_id varchar(36) not null,
    -- GUID
    app_id varchar(36) not null,
    device_id varchar(256) not null,

    -- Guest-<device_id>[5]
    name varchar(64) not null,

    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),

    PRIMARY KEY(game_id, app_id, user_id)
);
COMMENT ON TABLE anonymouses IS 'Anonymouses per game, app';

CREATE INDEX idx_anonymouses_user_id ON anonymouses(user_id);


-- TODO: add users_devices as collection for the same network, network_id or user_id to
-- collect possible devices.
