CREATE TABLE IF NOT EXISTS scooter_users (
                                           id SERIAL PRIMARY KEY,
                                           created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                                           updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                                           deleted_at TIMESTAMP WITH TIME ZONE,
                                           username VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL
    );
INSERT INTO scooter_users(username, password) VALUES (alex, alex1234);
INSERT INTO scooter_users(username, password) VALUES (jamon, serrano);
