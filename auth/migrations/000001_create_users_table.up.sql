CREATE TABLE IF NOT EXISTS users(
    id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    username TEXT NOT NULL UNIQUE,
    is_active BOOLEAN DEFAULT FALSE,
    is_staff BOOLEAN DEFAULT FALSE,
    is_superuser BOOLEAN DEFAULT FALSE,
    date_joined TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS user_id_email_region_indx ON users (id, email, geo_location);

CREATE TABLE user_profile(
    id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL UNIQUE,
    age INT,
    geo_location TEXT,
    thumbnail TEXT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS user_profile_id_user_id_indx ON user_profile (id, user_id);

CREATE TABLE user_posts(
    id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    user_profile_id UUID NOT NULL,
    post_content TEXT,
    post_image_URL TEXT,
    post_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (user_profile_id) REFERENCES user_profile(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS user_posts_id_user_profile_id_indx ON user_posts (id, user_profile_id);
