BEGIN;

INSERT INTO public.users (full_name, login, password_hash, role)
VALUES ('Учебный пользователь', 'lab_demo', 'login-disabled', 'operator')
ON CONFLICT (login) DO NOTHING;

INSERT INTO public.production_outputs (batch_number, product_name, produced_grams)
VALUES ('LAB-DEMO-01', 'Хлеб пшеничный', 10000)
ON CONFLICT (batch_number) DO NOTHING;

SELECT
    p.id AS production_output_id,
    r.id AS reason_id,
    u.id AS created_by
FROM public.production_outputs AS p
CROSS JOIN public.defect_reasons AS r
CROSS JOIN public.users AS u
WHERE p.batch_number = 'LAB-DEMO-01'
  AND r.code = 'TEMPERATURE_DEVIATION'
  AND u.login = 'lab_demo';

COMMIT;
