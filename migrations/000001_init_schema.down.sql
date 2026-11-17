BEGIN;

DROP FUNCTION IF EXISTS public.delete_defect(BIGINT);

DROP FUNCTION IF EXISTS public.update_defect(
    BIGINT,
    BIGINT,
    BIGINT,
    BIGINT,
    TEXT
    );

DROP FUNCTION IF EXISTS public.create_defect(
    BIGINT,
    BIGINT,
    BIGINT,
    BIGINT,
    TEXT
    );

DROP FUNCTION IF EXISTS public.get_defect(BIGINT);
DROP FUNCTION IF EXISTS public.get_defects();

DROP TABLE IF EXISTS public.defect_records;
DROP TABLE IF EXISTS public.defect_reasons;
DROP TABLE IF EXISTS public.production_outputs;
DROP TABLE IF EXISTS public.users;

COMMIT;