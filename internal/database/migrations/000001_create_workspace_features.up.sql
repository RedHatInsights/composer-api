CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE workspace_features (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workspace_id UUID NOT NULL,
    feature_id UUID NOT NULL,
    earmark INTEGER CHECK (earmark IS NULL OR earmark >= 0),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT unique_workspace_feature UNIQUE (workspace_id, feature_id)
);

CREATE INDEX idx_workspace_features_workspace ON workspace_features(workspace_id);
