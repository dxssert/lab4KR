BEGIN;

CREATE TABLE public.users
(
    id            BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    full_name     TEXT        NOT NULL,
    login         TEXT        NOT NULL UNIQUE,
    password_hash TEXT        NOT NULL,
    role          TEXT        NOT NULL,
    is_active     BOOLEAN     NOT NULL DEFAULT TRUE,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE public.production_outputs
(
    id             BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    batch_number   TEXT        NOT NULL UNIQUE,
    product_name   TEXT        NOT NULL,
    produced_grams BIGINT      NOT NULL CHECK (produced_grams > 0),
    produced_at    TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE public.defect_reasons
(
    id   BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    code TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL
);

CREATE TABLE public.defect_records
(
    id                   BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    production_output_id BIGINT      NOT NULL,
    reason_id             BIGINT      NOT NULL,
    defect_grams          BIGINT      NOT NULL CHECK (defect_grams > 0),
    comment               TEXT        NOT NULL DEFAULT '',
    created_by            BIGINT      NOT NULL,
    created_at            TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_defect_production_output
        FOREIGN KEY (production_output_id)
            REFERENCES public.production_outputs (id),

    CONSTRAINT fk_defect_reason
        FOREIGN KEY (reason_id)
            REFERENCES public.defect_reasons (id),

    CONSTRAINT fk_defect_created_by
        FOREIGN KEY (created_by)
            REFERENCES public.users (id)
);

CREATE INDEX idx_defect_records_production_output
    ON public.defect_records (production_output_id);

CREATE INDEX idx_defect_records_reason
    ON public.defect_records (reason_id);

CREATE INDEX idx_defect_records_created_by
    ON public.defect_records (created_by);

INSERT INTO public.defect_reasons (code, name)
VALUES
    ('EQUIPMENT_FAILURE', 'Неисправность оборудования'),
    ('TEMPERATURE_DEVIATION', 'Отклонение температуры'),
    ('RAW_MATERIAL', 'Несоответствие сырья'),
    ('OPERATOR_ERROR', 'Ошибка оператора');

CREATE OR REPLACE FUNCTION public.get_defects()
RETURNS TABLE
(
    id                   BIGINT,
    production_output_id BIGINT,
    reason_id             BIGINT,
    defect_grams          BIGINT,
    comment               TEXT,
    created_by            BIGINT,
    created_at            TIMESTAMPTZ,
    updated_at            TIMESTAMPTZ
)
LANGUAGE SQL
STABLE
AS
$$
SELECT
    d.id,
    d.production_output_id,
    d.reason_id,
    d.defect_grams,
    d.comment,
    d.created_by,
    d.created_at,
    d.updated_at
FROM public.defect_records AS d
ORDER BY d.created_at DESC;
$$;

CREATE OR REPLACE FUNCTION public.get_defect(p_id BIGINT)
RETURNS TABLE
(
    id                   BIGINT,
    production_output_id BIGINT,
    reason_id             BIGINT,
    defect_grams          BIGINT,
    comment               TEXT,
    created_by            BIGINT,
    created_at            TIMESTAMPTZ,
    updated_at            TIMESTAMPTZ
)
LANGUAGE SQL
STABLE
AS
$$
SELECT
    d.id,
    d.production_output_id,
    d.reason_id,
    d.defect_grams,
    d.comment,
    d.created_by,
    d.created_at,
    d.updated_at
FROM public.defect_records AS d
WHERE d.id = p_id;
$$;

CREATE OR REPLACE FUNCTION public.create_defect(
    p_production_output_id BIGINT,
    p_reason_id BIGINT,
    p_defect_grams BIGINT,
    p_created_by BIGINT,
    p_comment TEXT
)
RETURNS TABLE
(
    id                   BIGINT,
    production_output_id BIGINT,
    reason_id             BIGINT,
    defect_grams          BIGINT,
    comment               TEXT,
    created_by            BIGINT,
    created_at            TIMESTAMPTZ,
    updated_at            TIMESTAMPTZ
)
LANGUAGE SQL
AS
$$
    INSERT INTO public.defect_records
    (
        production_output_id,
        reason_id,
        defect_grams,
        comment,
        created_by
    )
    VALUES
    (
        p_production_output_id,
        p_reason_id,
        p_defect_grams,
        COALESCE(p_comment, ''),
        p_created_by
    )
    RETURNING
        defect_records.id,
        defect_records.production_output_id,
        defect_records.reason_id,
        defect_records.defect_grams,
        defect_records.comment,
        defect_records.created_by,
        defect_records.created_at,
        defect_records.updated_at;
$$;

CREATE OR REPLACE FUNCTION public.update_defect(
    p_id BIGINT,
    p_production_output_id BIGINT,
    p_reason_id BIGINT,
    p_defect_grams BIGINT,
    p_comment TEXT
)
RETURNS TABLE
(
    id                   BIGINT,
    production_output_id BIGINT,
    reason_id             BIGINT,
    defect_grams          BIGINT,
    comment               TEXT,
    created_by            BIGINT,
    created_at            TIMESTAMPTZ,
    updated_at            TIMESTAMPTZ
)
LANGUAGE SQL
AS
$$
UPDATE public.defect_records
SET
    production_output_id = p_production_output_id,
    reason_id = p_reason_id,
    defect_grams = p_defect_grams,
    comment = COALESCE(p_comment, ''),
    updated_at = CURRENT_TIMESTAMP
WHERE defect_records.id = p_id
    RETURNING
        defect_records.id,
        defect_records.production_output_id,
        defect_records.reason_id,
        defect_records.defect_grams,
        defect_records.comment,
        defect_records.created_by,
        defect_records.created_at,
        defect_records.updated_at;
$$;

CREATE OR REPLACE FUNCTION public.delete_defect(p_id BIGINT)
RETURNS SETOF BIGINT
LANGUAGE SQL
AS
$$
DELETE FROM public.defect_records
WHERE defect_records.id = p_id
    RETURNING defect_records.id;
$$;

COMMIT;