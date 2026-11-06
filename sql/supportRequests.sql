-- =====================================================================
-- Support requests (2026-10-24) — usersAK database (usersAppAK).
-- One general table for every user, company and platform (shao, ...).
-- The "Support" form of the apps saves a row here and emails it to
-- SUPPORT_EMAIL_TO (environment variable of usersAppAK).
-- Safe to run twice.
-- =====================================================================

CREATE TABLE IF NOT EXISTS public.tbl_support_request (
    support_request_id  UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
    ticket_code         VARCHAR(20)   NOT NULL UNIQUE,      -- SUP-YYMMDD-XXXX, shown to the user
    user_uid            VARCHAR(128)  NOT NULL,
    user_email          VARCHAR(254),
    company_id          VARCHAR(100),                       -- NULL when the user had no company selected
    company_name        VARCHAR(200),
    platform            VARCHAR(40)   NOT NULL DEFAULT 'shao',
    category            VARCHAR(30)   NOT NULL,             -- technical | account | billing | howTo | other
    priority            VARCHAR(10)   NOT NULL DEFAULT 'normal', -- low | normal | high
    subject             VARCHAR(200)  NOT NULL,
    description         TEXT          NOT NULL,
    contact_method      VARCHAR(20)   NOT NULL,             -- email | phone | whatsapp | videoCall
    contact_value       VARCHAR(200)  NOT NULL,
    preferred_date      DATE,
    time_slots          VARCHAR(60),                        -- morning,afternoon,evening
    image_urls          TEXT,                               -- one URL per line (Firebase)
    current_path        VARCHAR(300),
    language            VARCHAR(5),
    status              VARCHAR(20)   NOT NULL DEFAULT 'OPEN',   -- OPEN | IN_PROGRESS | CLOSED
    email_sent          BOOLEAN       NOT NULL DEFAULT FALSE,
    handled_notes       TEXT,
    handled_at          TIMESTAMPTZ,
    created_at          TIMESTAMPTZ   NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ   NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS ix_support_request_status ON public.tbl_support_request (status, created_at);
CREATE INDEX IF NOT EXISTS ix_support_request_user   ON public.tbl_support_request (user_uid, created_at);

-- ---------------------------------------------------------------------
-- Use:
--   Open requests:  SELECT created_at, ticket_code, company_name, user_email, category, priority, subject,
--                          contact_method, contact_value
--                   FROM public.tbl_support_request WHERE status <> 'CLOSED' ORDER BY created_at;
--   Close one:      UPDATE public.tbl_support_request SET status = 'CLOSED', handled_at = now(), updated_at = now(),
--                          handled_notes = 'Called on WhatsApp' WHERE ticket_code = 'SUP-...';
-- Revert:          DROP TABLE IF EXISTS public.tbl_support_request;
-- ---------------------------------------------------------------------
