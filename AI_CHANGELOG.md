# AI Change Log — usersAppAK

Keeps the latest 15 changes made with Claude (oldest are removed). Newest first.

---

## 17. 2026-11-04 — AI agent knowledge: clock in the top bar

**Backups:** `.ai_backups/2026-11-04_before_header_clock/`.

**Description:** works with p3 #44. `brules/bAIAgent/bAISystemKnowledge.go`: the Company page line now says the top bar
shows a clock with the company time (and the computer's time when it is in another zone), and the store shows the
company time too. Version "2026-11-04 (usersAppAK #17)". No endpoint or database change. Tests pass. AI knowledge:
updated.

**Modified:** brules/bAIAgent/bAISystemKnowledge.go

**Revert:** copy the `.bak` file back ("__" = "/"), rebuild and restart usersAppAK.

---

## 16. 2026-11-04 — AI agent knowledge: company time zone

**Backups:** `.ai_backups/2026-11-04_before_company_time_zone/`.

**Description:** works with shao #32. `brules/bAIAgent/bAISystemKnowledge.go`: the Company page line now says it has
the company **"Time zone"** (default Costa Rica, GMT−6; every hour of the company uses it; changing it does not move
saved appointments: 14:00 stays 14:00). Version "2026-11-04 (usersAppAK #16)". No endpoint, database or tracking
change. Test of the knowledge (version + paths) passes. AI knowledge: updated.

**Modified:** brules/bAIAgent/bAISystemKnowledge.go

**Revert:** copy the `.bak` file back ("__" = "/"), rebuild and restart usersAppAK.

---

## 15. 2026-11-01 — Request logs: request id, JSON logs with levels, panic -> clean 500

**Backups:** `.ai_backups/2026-11-01_before_request_logs/`.

**Description:**
- New `helpers/logger.go` (+ test): `SetupLogger()` in `main.go` after `.env` (same as shao #31: slog, LOG_FORMAT
  json/text, JSON default in production with `severity` / `message`).
- New `routers/request_log.go` (+ test): `requestLogMiddleware` replaces `loggingMiddleware` and `recoveryMiddleware`
  (removed from `routers/router.go`, with `statusRecorder` and the `time` import). One line per call (request_id,
  method, route template, status, duration_ms, bytes) instead of two lines with the raw path, remote address and
  origin; header `X-Request-ID`; panic logged with stack, JSON 500 with `requestId` and CORS headers.
- `env.example.NEW.txt`: `# LOG_FORMAT=`.
- No database change, no endpoint change. AI knowledge: not needed.

**Created:** helpers/logger.go, helpers/logger_test.go, routers/request_log.go, routers/request_log_test.go
**Modified:** main.go, routers/router.go, env.example.NEW.txt

**Revert:** copy the `.bak` files back ("__" = "/"), delete the created files, rebuild and restart usersAppAK.

---

## 14. 2026-11-01 — AI agent knowledge updated to today's menu, version + path test

**Backups:** `.ai_backups/2026-11-01_before_ai_knowledge_update/`.

**Description:**
- `brules/bAIAgent/bAISystemKnowledge.go`: `aiSystemKnowledge` used the old screen names and paths (Start an order,
  /shao/approvala, /shao/oporders...). Now it has today's menu (Sell, Bookings, Customers, My store, Reports, Team)
  with the new paths, the main new features (Showcase link, Your plan, appointment products, Notify the customer per
  step, customer service email, Contact us buttons, AlephK general email) and sign-in with email, Google or Microsoft.
- New `aiKnowledgeVersion` ("2026-11-01 (usersAppAK #14)") sent to Claude; new `aiScreenPaths` (must match p3
  `assistantRoutes.js`) + `bAISystemKnowledge_test.go` (`TestAIKnowledgePaths`, `TestAIKnowledgeVersion`).
- Rule: same as shao #30 (update the AI knowledge in every change that touches screens or features).
- No endpoint, database or tracking change. AI knowledge: updated.

**Created:** brules/bAIAgent/bAISystemKnowledge_test.go
**Modified:** brules/bAIAgent/bAISystemKnowledge.go

**Revert:** copy the `.bak` file back ("__" = "/"), delete the test file, rebuild and restart usersAppAK.

---

## 13. 2026-10-31 — Domains from .env (aleph-k.com, aleph-k.online): CORS, login cookie, CAPTCHA

**Backups:** `.ai_backups/2026-10-31_before_domains/`.

**Description:** new `helpers/allowed_domains.go` (+ test), same rules as shao. `ALLOWED_DOMAINS=aleph-k.com,aleph-k.online`
= each domain and all its subdomains, http and https:
- CORS (`routers/cors_config.go`, `routers/router.go`) accepts them; domain lines allowed in the origins file.
- Login / logout cookie (`HuserLogin.go`, `securityGen/HuserLogout.go`): Domain = the allowed domain that contains the
  API host (`auth.aleph-k.com` → `aleph-k.com`; `AUTH_COOKIE_DOMAIN` may list several and wins); localhost / LAN =
  host-only as before. Secure when the browser used https (`AUTH_COOKIE_SECURE=true|false` forces it), so login works on
  http and https. Before, a set `AUTH_COOKIE_DOMAIN` always made the cookie Secure (login broken on http) and one value
  could not serve two domains.
- CAPTCHA (`brules/bCaptcha`): `*.domain` entries; empty + production = only `ALLOWED_DOMAINS`.
- Docs: `env.example.NEW.txt`, `cors_allowed_origins.example.txt`.

**Created:** helpers/allowed_domains.go, helpers/allowed_domains_test.go
**Modified:** routers/cors_config.go, routers/router.go, brules/bCaptcha/bCaptchaVerify.go,
handlers/HuserLogin/HuserLogin.go, handlers/securityGen/HuserLogout.go, env.example.NEW.txt, cors_allowed_origins.example.txt

**Revert:** copy the `.bak` files back ("__" = "/"), delete the created files, rebuild usersAppAK.

---

## 12. 2026-10-29 — Login keeps the real sign-in provider (fixes "external provider" on change email/password)

**Backups:** `.ai_backups/2026-10-29_before_provider_fix/`.

**Description:** on every login, when the user already existed, `HuserLogin` saved `ProviderID` and `PhotoURL` from
claims named "ProviderID" / "PhotoURL", which do not exist in the Firebase token, so both were saved **empty**. Change
email / change password check `ProviderID == "password"`, so an email+password user was told the account "uses an
external provider". Now the provider comes from the token (`firebase.sign_in_provider`, e.g. `password`,
`google.com`) and the photo from the `picture` claim; when the token has no value the saved one is kept. Users already
affected are fixed automatically at their next login (sign out and sign in again). No DB change.

**Modified:** handlers/HuserLogin/HuserLogin.go

**Revert:** copy the `.bak` file back, rebuild and restart usersAppAK.

---

## 11. 2026-10-29 — Invitation email includes a link to the app (APP_PUBLIC_URL)

**Backups:** `.ai_backups/2026-10-29_before_invitation_link/`.

**Description:** new `.env` setting `APP_PUBLIC_URL` (base address of p3). The invitation email adds "To accept or
decline the invitation, open: <APP_PUBLIC_URL>/manageinvitatinsbyuser" plus a note to create an account with the same
email first. Empty setting = email as before (no link) + a log warning. Change `.env` and restart, no rebuild.
Test `TestInvitationLink`.

**Modified:** brules/bcreateUsers/BSendInvtationEmail.go, BSendInvtationEmail_test.go, env.example.NEW.txt

**Revert:** copy the `.bak` files back, rebuild and restart usersAppAK.

---

## 10. 2026-10-28 — Invitations use the general email settings from .env (no hard-coded SMTP)

**Backups:** `.ai_backups/2026-10-28_before_invitation_smtp_env/`.

**Description:** `BSendInvtationEmail.go` had the mail server, user and **password** written in the code. Now
`LoadGeneralSMTPConfig()` reads `SMTP_HOST`, `SMTP_PORT` (587), `SMTP_USER`, `SMTP_PASSWORD`, `SMTP_FROM` (default
SMTP_USER) and the new `SMTP_FROM_NAME` (default "AlephK") from `.env` — the same settings support requests use.
Sender shows as "Name <address>", subject "<Name> invitation", signature "<Name>". Missing settings return
`ErrEmailNotConfigured` (the endpoint answers 500 "The invitation could not be sent"). Change `.env` + restart, no rebuild.
New test `BSendInvtationEmail_test.go`.

**Modified:** brules/bcreateUsers/BSendInvtationEmail.go, env.example.NEW.txt
**Created:** brules/bcreateUsers/BSendInvtationEmail_test.go

**Revert:** copy the `.bak` files back, delete the test file, rebuild and restart usersAppAK.

---

## 9. 2026-10-28 — Login cookie: Domain + Secure from AUTH_COOKIE_DOMAIN

**Backups:** `.ai_backups/2026-10-28_before_cookie_domain/`.

**Description:** `HuserLogin.go` sets the `access_token` cookie with `Domain = AUTH_COOKIE_DOMAIN` and `Secure = true`
when that variable is set (not localhost / .local), the same rule logout already used. Empty variable = exactly as before
(no Domain, Secure=false). Tracking unchanged. `env.example.NEW.txt` explains the variable.

**Modified:** handlers/HuserLogin/HuserLogin.go, env.example.NEW.txt

**Revert:** copy the `.bak` files back, rebuild and restart usersAppAK.

---

## 8. 2026-10-27 — Facebook sign-in with an email counts as verified email; tracking on /createuser

**Backups:** `.ai_backups/2026-10-27_before_facebook_email_verified/`.

**Description:** Firebase marks Facebook emails as not verified, so Facebook users would see "Verify your email".
New `brules/bcreateUsers/bEmailVerified.go` → `BEffectiveEmailVerified(claim, provider, email)`: provider
`facebook.com` with a non-empty email → `"true"`; everything else keeps Firebase's value (Google, Microsoft — kept
strict on purpose for work/school accounts —, Apple, email + password, Facebook without email). Used in
`HuserLogin.go` (new user + the update comparison, so existing Facebook users are fixed at their next login) and
`HuserCreate.go`. `/createuser` now writes tracking (lt) `CreateUser` / `CreateUserFailed` with provider and
emailVerified (no token). Unit test `bEmailVerified_test.go`. No DB or frontend change.

**Created:** brules/bcreateUsers/bEmailVerified.go, bEmailVerified_test.go
**Modified:** handlers/HuserLogin/HuserLogin.go, handlers/HuserCreate/HuserCreate.go

**Revert:** copy the `.bak` files back, delete the two new files, rebuild and restart usersAppAK.

---

## 7. 2026-10-25 — Names: Aleph / AlephTelos in the AI agent instructions

**Backups:** `.ai_backups/2026-10-25_before_alephtelos_names/`.

**Description:** texts only: the agent introduces itself as "Aleph, the AlephK assistant"; the system knowledge calls
the main platform AlephTelos ("its old name was Service Hub") and AlephK a brand of Kuantum K.

**Modified:** brules/bAIAgent/bAISystemKnowledge.go

**Revert:** copy the `.bak` back and rebuild.

---

## 6. 2026-10-25 — CAPTCHA (Cloudflare Turnstile) on support requests

**Control point:** CC 2. Backups: `.ai_backups/CC2_before_captcha/`. No database change.

**Description:** new `brules/bCaptcha/bCaptchaVerify.go` (same verifier as shao) + tests. `POST /supportrequests` checks the
`X-Captcha-Token` header (400 / 503 with `data.code`, tracking `SupportRequestCaptchaFailed`). CORS allows the header.
Environment: `TURNSTILE_SECRET_KEY` (empty = CAPTCHA off; on your computer leave it empty or use the test secret `1x0000000000000000000000000000000AA`), optional `CAPTCHA_PROVIDER`, `CAPTCHA_FAIL_OPEN=1`, `TURNSTILE_ALLOWED_HOSTNAMES`. Lines added to `env.example.NEW.txt`.

**Created:** brules/bCaptcha/bCaptchaVerify.go (+_test)
**Modified:** handlers/HsupportRequest/HCreateSupportRequest.go, routers/router.go, env.example.NEW.txt

**Revert:** copy the `.bak` files back, delete the created folder, rebuild usersAppAK.

---

## 5. 2026-10-25 — Security fixes after control point CC 1 (S3, S5, S9, S10, S11, S15)

**Control point:** CC 1 = the state before this change. Backups: `.ai_backups/CC1_before_security_fixes/` (README lists every file). No database change.

**Description:**
- **S3:** `status` (and `plan`) are never taken from the company form: new = FREE / ACTIVE, update keeps the saved values
  (HcompanyCreate, BCreateOrUpdateCompany, BUpdateCompany).
- **S10:** `GET /searchUser` only for the OWNER / MANAGER of `companyId` (query param or X-Company-ID header,
  security.RequireActiveCompanyUser), exact complete email only, max 1 user, tracked (`SearchUserByEmail`, without the
  searched email). New `database/databaseUsers/dSearchUserByExactEmail.go`; the old partial search function is unused.
- **S11:** the Authorization header / token is no longer printed (HuserCreate, helpers/auth_helper.go); no company name /
  identification number / emails in logs (bCreateCompany, DSelectUserByEmail, HuserGetBYDesc).
- **S15:** `handlers/common/safeErrors.go`: business messages stay, technical details (DB, SQL, Firebase) go to the log
  and the browser gets a generic text — applied in 20 handlers ("Invalid JSON body" without details too).
- **S9 CORS:** same as shao (APP_ENV production / test / not set) — routers/cors_config.go.
- **S5:** .gitignore adds connection.json, databaseConnection.json, .ai_backups/, *.pem, *.key.
- Tests added: safeErrors, exact-email search. Build + tests OK (with stub libraries).

**Created:** handlers/common/safeErrors.go (+_test), brules/bcreateUsers/bsearchUsersByEmailNameOrId_test.go,
database/databaseUsers/dSearchUserByExactEmail.go, routers/cors_config.go, cors_allowed_origins.example.txt

**Modified:** (new `.env.example` content is in `env.example.NEW.txt` — rename it; remote tools cannot write .env files), .gitignore, routers/router.go, brules/bCompany/bCreateCompany.go, bUpdateCompany.go,
brules/bcreateUsers/bsearchUsersByEmailNameOrId.go, database/databaseUsers/DSelectUserByEmail.go, helpers/auth_helper.go,
handlers/users/HuserGetBYDesc.go, handlers/HuserCreate/HuserCreate.go, handlers/HuserLogin/HuserLogin.go,
HSendInvitationRequest.go, handlers/common/notifications.go, handlers/securityGen/HChangeEmail.go, HChangePassword.go,
CanChangeFirebasePasswordResponse.go, HuserMe.go, handlers/Hcompany/HcompanyCreate.go, HcompanyGetByUser.go,
hgetAllCompanies.go, hRequestUsersCompany.go, HUpdateRequestInputCancelByManager.go, HProcessCompanyUserInvitation.go,
HTransferCompanyOwner.go, HUserPendingRequestsFlagResponse.go, HGetPendingInvitationsByUser.go, HSearchCompanyUsers.go,
HCompanyUserRequestSearchResponse.go, UpdateCompanyUserResponse.go

**Revert:** copy every `.bak` of `.ai_backups/CC1_before_security_fixes/` back (README), delete the created files, rebuild usersAppAK.

---

## 4. 2026-10-24 — Support requests (general table), environment config, test routes removed

**Database (run first, usersAK):** `sql/supportRequests.sql` — new table `tbl_support_request` (one general table for all
users, companies and platforms). Safe to run twice. Usage and revert SQL inside.

**Description:**
- `POST /supportrequests` (logged-in user, `RequireActiveUser`): validates the form (topic, priority, contact method, time slots,
  up to 3 https image URLs), max 10 requests per user per day, saves the row with a ticket code `SUP-YYMMDD-XXXX` and emails it
  (in the background) to `SUPPORT_EMAIL_TO` with `SMTP_HOST / SMTP_PORT / SMTP_USER / SMTP_PASSWORD / SMTP_FROM`.
  `email_sent` is set when the email works. Tracking `CreateSupportRequest` / `CreateSupportRequestFailed`.
- `main.go` loads `.env` and reads `PORT` (default 8080); the router adds `CORS_ALLOWED_ORIGINS`.
- Removed the test routes `/profile` and `/testwithauth`; "/" answers a small JSON health check.
  Dead files to delete by hand: `handlers/tetHandler.go`, `handlers/HuserLogin/HLogout.go` (nothing uses them).
- `.gitignore` (.env, Firebase admin json, binaries) and `.env.example`.

**Files created:** `database/databaseSupportRequest/` (3 files), `brules/bSupportRequest/` (2 files),
`handlers/HsupportRequest/HCreateSupportRequest.go`, `sql/supportRequests.sql`, `.gitignore`, `.env.example`,
`.ai_backups/2026-10-24_before_report2_fixes/`.

**Files modified:** `main.go`, `routers/router.go`.

**Revert:** see the backup `README.txt`. Revert p3 #21 too (the support form would fail without the endpoint).

---

## 3. 2026-10-22 — Company plan cannot be changed by the owner

**Description:** `bCreateCompany.BCreateOrUpdateCompany`: on update the saved plan is kept; new companies always start
on FREE. `BUpdateCompany` keeps the saved plan too. Plans are changed by AlephK (SQL) after a plan request from shao.

**Revert:** restore the `.bak` files from `.ai_backups/2026-10-22_before_company_plans/`.
