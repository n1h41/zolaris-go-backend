CREATE TABLE IF NOT EXISTS z_referrals (
    referral_id uuid PRIMARY KEY NOT NULL DEFAULT GEN_RANDOM_UUID (),
    FOREIGN KEY (referrer_id) REFERENCES z_users (user_id) ON DELETE CASCADE,
    FOREIGN KEY (referee_id) REFERENCES z_users (user_id) ON DELETE CASCADE,
    status text CHECK (status IN ('active', 'used', 'cancelled')) DEFAULT 'active',
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    used_at timestamp with time zone,
)
