-- =====================================================================
-- AI Agent conversations (usersAppAK)
-- Stores every conversation and message between logged users and the
-- Anthropic (Claude Haiku) agent.
-- =====================================================================

CREATE TABLE IF NOT EXISTS public.ai_conversations (
    conversation_id      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_uid             VARCHAR(128) NOT NULL,      -- Firebase UID (claims.UID)
    user_email           VARCHAR(255),
    company_id           VARCHAR(64),                -- optional context
    platform             VARCHAR(50)  NOT NULL DEFAULT 'UserManagement',
    title                VARCHAR(200),
    model                VARCHAR(100) NOT NULL,
    status               VARCHAR(20)  NOT NULL DEFAULT 'ACTIVE',
    total_input_tokens   INT NOT NULL DEFAULT 0,
    total_output_tokens  INT NOT NULL DEFAULT 0,
    created_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at           TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_ai_conv_user
    ON public.ai_conversations (user_uid, created_at DESC);

CREATE TABLE IF NOT EXISTS public.ai_conversation_messages (
    message_id       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    conversation_id  UUID NOT NULL
                     REFERENCES public.ai_conversations(conversation_id) ON DELETE CASCADE,
    user_uid         VARCHAR(128) NOT NULL,
    role             VARCHAR(20)  NOT NULL CHECK (role IN ('user', 'assistant')),
    content          TEXT NOT NULL,
    model            VARCHAR(100),
    input_tokens     INT,
    output_tokens    INT,
    stop_reason      VARCHAR(50),
    latency_ms       INT,
    error_message    TEXT,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_ai_msg_conv
    ON public.ai_conversation_messages (conversation_id, created_at);

-- =====================================================================
-- REVERT
-- DROP TABLE IF EXISTS public.ai_conversation_messages;
-- DROP TABLE IF EXISTS public.ai_conversations;
-- =====================================================================
