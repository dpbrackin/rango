-- +goose Up
-- +goose StatementBegin
CREATE TABLE indexes(
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  org_id UUID REFERENCES organizations(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  engine TEXT NOT NULL
);

CREATE TABLE index_documents(
  index_id UUID REFERENCES indexes(id) on DELETE CASCADE,
  document_id UUID REFERENCES documents(id) on DELETE CASCADE,
  PRIMARY KEY(index_id, document_id)
);

CREATE UNIQUE INDEX unique_index_name_per_org
ON indexes(org_id, name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS indexes;
DROP TABLE IF EXISTS index_documents;
-- +goose StatementEnd
