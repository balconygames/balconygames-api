CREATE TABLE leaderboards (
    -- external id, used to use to hide
    -- database real id skip possible
    -- lookup ups just in case.
    -- GUID
    id varchar(36) not null,

    name varchar(256) not null,

    -- GUID
    game_id varchar(36) not null,
    -- GUID
    app_id varchar(36) not null,

    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),

    PRIMARY KEY(id, game_id, app_id)
);
CREATE INDEX idx_leaderboards_game_id_app_id ON leaderboards(game_id, app_id);