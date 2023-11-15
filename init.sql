CREATE TABLE IF NOT EXISTS credentials (
                                           id SERIAL PRIMARY KEY,
                                           created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                                           updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                                           deleted_at TIMESTAMP WITH TIME ZONE,
                                           username VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL
    );
