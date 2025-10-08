CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    token_hash text NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    revoked_at TIMESTAMP WITH TIME ZONE NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indeks
CREATE INDEX idx_refresh_tokens_token_hash ON refresh_tokens USING btree (token_hash);
CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens USING btree (user_id);
CREATE INDEX idx_refresh_tokens_expires_at ON refresh_tokens USING btree (expires_at);