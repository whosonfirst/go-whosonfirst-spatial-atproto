package duckdb

/*

D DESCRIBE SELECT * FROM read_parquet('https://data.geocode.earth/wof/dist/parquet/whosonfirst-data-admin-latest.parquet');
┌───────────────┬────────────────────────────────────────────────────────┬──────┬─────┬─────────┬───────┐
│  column_name  │                      column_type                       │ null │ key │ default │ extra │
├───────────────┼────────────────────────────────────────────────────────┼──────┼─────┼─────────┼───────┤
│ id            │ VARCHAR                                                │ YES  │     │         │       │
│ parent_id     │ INTEGER                                                │ YES  │     │         │       │
│ name          │ VARCHAR                                                │ YES  │     │         │       │
│ placetype     │ VARCHAR                                                │ YES  │     │         │       │
│ placelocal    │ VARCHAR                                                │ YES  │     │         │       │
│ country       │ VARCHAR                                                │ YES  │     │         │       │
│ repo          │ VARCHAR                                                │ YES  │     │         │       │
│ lat           │ DOUBLE                                                 │ YES  │     │         │       │
│ lon           │ DOUBLE                                                 │ YES  │     │         │       │
│ min_lat       │ DOUBLE                                                 │ YES  │     │         │       │
│ min_lon       │ DOUBLE                                                 │ YES  │     │         │       │
│ max_lat       │ DOUBLE                                                 │ YES  │     │         │       │
│ max_lon       │ DOUBLE                                                 │ YES  │     │         │       │
│ min_zoom      │ VARCHAR                                                │ YES  │     │         │       │
│ max_zoom      │ VARCHAR                                                │ YES  │     │         │       │
│ min_label     │ VARCHAR                                                │ YES  │     │         │       │
│ max_label     │ VARCHAR                                                │ YES  │     │         │       │
│ modified      │ DATE                                                   │ YES  │     │         │       │
│ is_funky      │ VARCHAR                                                │ YES  │     │         │       │
│ population    │ VARCHAR                                                │ YES  │     │         │       │
│ country_id    │ VARCHAR                                                │ YES  │     │         │       │
│ region_id     │ VARCHAR                                                │ YES  │     │         │       │
│ county_id     │ VARCHAR                                                │ YES  │     │         │       │
│ concord_ke    │ VARCHAR                                                │ YES  │     │         │       │
│ concord_id    │ VARCHAR                                                │ YES  │     │         │       │
│ iso_code      │ VARCHAR                                                │ YES  │     │         │       │
│ hasc_id       │ VARCHAR                                                │ YES  │     │         │       │
│ gn_id         │ VARCHAR                                                │ YES  │     │         │       │
│ wd_id         │ VARCHAR                                                │ YES  │     │         │       │
│ name_ara      │ VARCHAR                                                │ YES  │     │         │       │
│ name_ben      │ VARCHAR                                                │ YES  │     │         │       │
│ name_deu      │ VARCHAR                                                │ YES  │     │         │       │
│ name_eng      │ VARCHAR                                                │ YES  │     │         │       │
│ name_ell      │ VARCHAR                                                │ YES  │     │         │       │
│ name_fas      │ VARCHAR                                                │ YES  │     │         │       │
│ name_fra      │ VARCHAR                                                │ YES  │     │         │       │
│ name_heb      │ VARCHAR                                                │ YES  │     │         │       │
│ name_hin      │ VARCHAR                                                │ YES  │     │         │       │
│ name_hun      │ VARCHAR                                                │ YES  │     │         │       │
│ name_ind      │ VARCHAR                                                │ YES  │     │         │       │
│ name_ita      │ VARCHAR                                                │ YES  │     │         │       │
│ name_jpn      │ VARCHAR                                                │ YES  │     │         │       │
│ name_kor      │ VARCHAR                                                │ YES  │     │         │       │
│ name_nld      │ VARCHAR                                                │ YES  │     │         │       │
│ name_pol      │ VARCHAR                                                │ YES  │     │         │       │
│ name_por      │ VARCHAR                                                │ YES  │     │         │       │
│ name_rus      │ VARCHAR                                                │ YES  │     │         │       │
│ name_spa      │ VARCHAR                                                │ YES  │     │         │       │
│ name_swe      │ VARCHAR                                                │ YES  │     │         │       │
│ name_tur      │ VARCHAR                                                │ YES  │     │         │       │
│ name_ukr      │ VARCHAR                                                │ YES  │     │         │       │
│ name_urd      │ VARCHAR                                                │ YES  │     │         │       │
│ name_vie      │ VARCHAR                                                │ YES  │     │         │       │
│ name_zho      │ VARCHAR                                                │ YES  │     │         │       │
│ geom_src      │ VARCHAR                                                │ YES  │     │         │       │
│ geometry      │ BLOB                                                   │ YES  │     │         │       │
│ geometry_bbox │ STRUCT(xmin FLOAT, ymin FLOAT, xmax FLOAT, ymax FLOAT) │ YES  │     │         │       │
└───────────────┴────────────────────────────────────────────────────────┴──────┴─────┴─────────┴───────┘

*/
